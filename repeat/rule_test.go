package repeat

import (
	"reflect"
	"testing"
)

func TestOneWeek(t *testing.T) {
	subject := `1 н. Иностанный язык`
	expected := ParsedSubject{
		Rule: Rule{
			Mode:  Once,
			Dates: []int{1},
		},
		Subject:   `Иностанный язык`,
		StartWeek: 1,
	}
	result := Parse(subject)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v\nGot: %v\n", expected, result)
	}
}

func TestOneWeekExtraSpaces(t *testing.T) {
	subject := `1     н.      Иностанный язык   `
	expected := ParsedSubject{
		Rule: Rule{
			Mode:  Once,
			Dates: []int{1},
		},
		Subject:   `Иностанный язык`,
		StartWeek: 1,
	}
	result := Parse(subject)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v\nGot: %v\n", expected, result)
	}
}

func TestAny(t *testing.T) {
	subject := `Иностанный язык`
	expected := ParsedSubject{
		Rule:    Rule{Mode: Any},
		Subject: `Иностанный язык`,
	}
	result := Parse(subject)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v\nGot: %v\n", expected, result)
	}
}

func TestAnyExtraSpaces(t *testing.T) {
	subject := `    Иностанный язык    `
	expected := ParsedSubject{
		Rule:    Rule{Mode: Any},
		Subject: `Иностанный язык`,
	}
	result := Parse(subject)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v\nGot: %v\n", expected, result)
	}
}

func TestRange(t *testing.T) {
	subject := `  3   -     17    н. Иностанный язык  `
	expected := ParsedSubject{
		Rule: Rule{
			Mode:  Range,
			Dates: []int{3, 17},
		},
		Subject:   `Иностанный язык`,
		StartWeek: 3,
	}
	result := Parse(subject)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v\nGot: %v\n", expected, result)
	}
}

func TestEnum(t *testing.T) {
	subject := ` 3,5,  7 , 9 н. Иностанный язык `
	expected := ParsedSubject{
		Rule: Rule{
			Mode:  Enum,
			Dates: []int{3, 5, 7, 9},
		},
		Subject:   `Иностанный язык`,
		StartWeek: 3,
	}
	result := Parse(subject)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v\nGot: %v\n", expected, result)
	}
}

func TestExcept(t *testing.T) {
	subject := `кр.  11  н. Иностранный язык  `
	expected := ParsedSubject{
		Rule: Rule{
			Mode:   Any,
			Except: []int{11},
		},
		Subject: `Иностанный язык`,
	}
	result := Parse(subject)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v\nGot: %v\n", expected, result)
	}
}

func TestExceptEnum(t *testing.T) {
	subject := `кр 5,7 ,  11  н. Иностанный язык  `
	expected := ParsedSubject{
		Rule: Rule{
			Mode:   Any,
			Except: []int{5, 7, 11},
		},
		Subject: `Иностанный язык`,
	}
	result := Parse(subject)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v\nGot: %v\n", expected, result)
	}
}

func TestStartAt(t *testing.T) {
	subject := `с 11  н. Иностанный язык  `
	expected := ParsedSubject{
		Rule:      Rule{Mode: Any},
		Subject:   `Иностанный язык`,
		StartWeek: 11,
	}
	result := Parse(subject)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v\nGot: %v\n", expected, result)
	}
}