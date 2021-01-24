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
)

func Parse(subject string) ParsedSubject {
	var result ParsedSubject

	splitted := strings.SplitN(subject, "н.", 2)

	if len(splitted) < 2 {
		return result
	}

	result.Subject = splitted[1]
	datesStr := strings.Split(splitted[0], "-")
	switch len(datesStr) {
	case 1:
		result.parseAsSingle(datesStr)
	case 2:
		result.parseAsRange(datesStr)
	default:
		log.Println("Unable to parse", datesStr)
	}

	return result
}

func (result *ParsedSubject) parseAsSingle(dates []string) {
	week, err := strconv.Atoi(strings.TrimSpace(dates[0]))
	if err != nil {
		log.Panic(err)
	}
	result.Rule.Mode = Once
	result.Rule.Dates = []int{week}
	result.StartWeek = week
}

func (result *ParsedSubject) parseAsRange(dates []string) {
	start, err := strconv.Atoi(strings.TrimSpace(dates[0]))
	if err != nil {
		log.Panic(err)
	}

	end, err := strconv.Atoi(strings.TrimSpace(dates[1]))
	if err != nil {
		log.Panic(err)
	}

	result.Rule.Mode = Range
	result.Rule.Dates = []int{start, end}
	result.StartWeek = start
}
