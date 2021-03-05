package repeat

import (
	"reflect"
	"testing"
)

func TestOneWeek(t *testing.T) {
	subjects := []string{
		`1     н.....      Иностранный язык   `,
		`1н Иностранный язык`,
	}

	expected := ParsedSubject{
		Rule: Rule{
			Mode:  Once,
			Dates: []int{1},
		},
		Subject:   `Иностранный язык`,
		StartWeek: 1,
	}

	for _, subject := range subjects {
		result := Parse(subject)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v\nGot: %v\n", expected, result)
		}
	}
}

func TestAny(t *testing.T) {
	subjects := []string{
		`Иностранный язык`,
		`    Иностранный язык    `,
	}

	expected := ParsedSubject{
		Rule:    Rule{Mode: Any},
		Subject: `Иностранный язык`,
	}

	for _, subject := range subjects {
		result := Parse(subject)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v\nGot: %v\n", expected, result)
		}
	}
}

func TestRange(t *testing.T) {
	subjects := []string{
		`  3   -     17    н... Иностранный язык  `,
		`3-17н Иностранный язык`,
	}

	expected := ParsedSubject{
		Rule: Rule{
			Mode:  Range,
			Dates: []int{3, 17},
		},
		Subject:   `Иностранный язык`,
		StartWeek: 3,
	}

	for _, subject := range subjects {
		result := Parse(subject)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v\nGot: %v\n", expected, result)
		}
	}
}

func TestEnum(t *testing.T) {
	subjects := []string{
		` 3,5,  7 , 9 н. Иностранный язык `,
		`3,5,7,9н Иностранный язык`,
		`3,5,7,9   н.... Иностранный язык `,
	}

	expected := ParsedSubject{
		Rule: Rule{
			Mode:  Enum,
			Dates: []int{3, 5, 7, 9},
		},
		Subject:   `Иностранный язык`,
		StartWeek: 3,
	}

	for _, subject := range subjects {
		result := Parse(subject)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v\nGot: %v\n", expected, result)
		}
	}
}

func TestExcept(t *testing.T) {
	subjects := []string{
		`  кр...  11  н... Иностранный язык  `,
		`кр11н Иностранный язык`,
		`кр.11н. Иностранный язык`,
		`кр 11 н Иностранный язык`,
	}

	expected := ParsedSubject{
		Rule: Rule{
			Mode:   Any,
			Except: []int{11},
		},
		Subject: `Иностранный язык`,
	}

	for _, subject := range subjects {
		result := Parse(subject)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v\nGot: %v\n", expected, result)
		}
	}
}

func TestExceptEnum(t *testing.T) {
	subjects := []string{
		`кр.. 5,7 ,  11  н... Иностранный язык  `,
		`кр5,7,11н Иностранный язык`,
		`кр 5,7,11 н Иностранный язык`,
	}

	expected := ParsedSubject{
		Rule: Rule{
			Mode:   Any,
			Except: []int{5, 7, 11},
		},
		Subject: `Иностранный язык`,
	}

	for _, subject := range subjects {
		result := Parse(subject)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v\nGot: %v\n", expected, result)
		}
	}
}

func TestStartAt(t *testing.T) {
	subjects := []string{
		`с 11  н... Иностранный язык  `,
		`с11н Иностранный язык`,
	}

	expected := ParsedSubject{
		Rule:      Rule{Mode: Any},
		Subject:   `Иностранный язык`,
		StartWeek: 11,
	}

	for _, subject := range subjects {
		result := Parse(subject)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v\nGot: %v\n", expected, result)
		}
	}
}

func TestRangeWithEx(t *testing.T) {
	subjects := []string{
		`  3   -     17    н... (кр.... 11 н...) Иностранный язык  `,
		` 3 -17  н. кр  11 н. Иностранный язык `,
		`3-17н кр11н. Иностранный язык`,
		`3-17н (кр11н) Иностранный язык`,
	}

	expected := ParsedSubject{
		Rule: Rule{
			Mode:   Range,
			Dates:  []int{3, 17},
			Except: []int{11},
		},
		Subject:   `Иностранный язык`,
		StartWeek: 3,
	}

	for _, subject := range subjects {
		result := Parse(subject)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v\nGot: %v\n", expected, result)
		}
	}
}
