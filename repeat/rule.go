package repeat

import (
	"log"
	"strconv"
	"strings"
)

type ParsedSubject struct {
	Rule      Rule
	Subject   string
	StartWeek int
}

type Rule struct {
	Mode  int
	Dates []int
}

const (
	Any   = iota
	Once  = iota
	Range = iota
	Enum  = iota
)

func Parse(subject string) ParsedSubject {
	var result ParsedSubject

	splitted := strings.SplitN(subject, "н.", 2)

	if len(splitted) != 2 {
		result.Rule.Mode = Any
		result.Subject = strings.TrimSpace(subject)
		return result
	}

	dates := splitted[0]
	result.Subject = strings.TrimSpace(splitted[1])

	if strings.Count(dates, "-") == 1 {
		result.parseAsRange(dates)
	} else if strings.Count(dates, ",") > 0 {
		result.parseAsEnum(dates)
	} else {
		result.parseAsSingle(dates)
	}

	return result
}

func (result *ParsedSubject) parseAsSingle(dates string) {
	trimmed := strings.TrimSpace(dates)

	if strings.HasPrefix(trimmed, "с ") {
		weekString := strings.TrimPrefix(trimmed, "с ")
		week, err := strconv.Atoi(strings.TrimSpace(weekString))
		if err != nil {
			log.Panic(err)
		}
		result.Rule.Mode = Any
		result.StartWeek = week
		return
	}

	week, err := strconv.Atoi(trimmed)
	if err != nil {
		log.Panic(err)
	}
	result.Rule.Mode = Once
	result.Rule.Dates = []int{week}
	result.StartWeek = week
}

func (result *ParsedSubject) parseAsRange(dates string) {
	d := strings.Split(dates, "-")
	start, err := strconv.Atoi(strings.TrimSpace(d[0]))
	if err != nil {

	}

	end, err := strconv.Atoi(strings.TrimSpace(d[1]))
	if err != nil {
		log.Panic(err)
	}

	result.Rule.Mode = Range
	result.Rule.Dates = []int{start, end}
	result.StartWeek = start
}

func (result *ParsedSubject) parseAsEnum(dates string) {
	d := strings.Split(dates, ",")
	for _, week := range d {
		week = strings.TrimSpace(week)
		num, err := strconv.Atoi(week)
		if err != nil {
			log.Panic(err)
		}
		result.Rule.Dates = append(result.Rule.Dates, num)
	}

	result.Rule.Mode = Enum
}
