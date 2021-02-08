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
	Mode   RepeatMode
	Dates  []int
	Except []int
}

type RepeatMode int

const (
	_     RepeatMode = iota
	Any   RepeatMode = iota
	Once  RepeatMode = iota
	Range RepeatMode = iota
	Enum  RepeatMode = iota
)

func Parse(subject string) ParsedSubject {
	result := ParsedSubject{
		Rule:    Rule{Mode: Any},
		Subject: strings.TrimSpace(subject),
	}

	splitted := strings.SplitN(subject, "н", 2)
	if len(splitted) != 2 {
		return result
	}

	hasNumbers := strings.ContainsAny(splitted[0], "0123456789")
	if !hasNumbers {
		return result
	}

	dates := splitted[0]
	splitted[1] = strings.Trim(splitted[1], ".")
	result.Subject = strings.TrimSpace(splitted[1])

	switch {
	case strings.Count(dates, "-") == 1:
		result.parseAsRange(dates)
	case strings.Count(dates, ",") > 0:
		result.parseAsEnum(dates)
	default:
		result.parseAsSingle(dates)
	}

	return result
}

func (result *ParsedSubject) parseAsSingle(dates string) {
	trimmed := strings.TrimSpace(dates)

	if strings.HasPrefix(trimmed, "с") {
		weekString := strings.TrimPrefix(trimmed, "с")
		week, err := strconv.Atoi(strings.TrimSpace(weekString))
		if err != nil {
			log.Panic(err)
		}
		result.Rule.Mode = Any
		result.StartWeek = week
		return
	}

	if strings.HasPrefix(trimmed, "кр") {
		weekString := strings.TrimLeft(trimmed, "кр. ")
		week, err := strconv.Atoi(strings.TrimSpace(weekString))
		if err != nil {
			log.Panic(err)
		}
		result.Rule.Mode = Any
		result.Rule.Except = []int{week}
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
		log.Panic(err)
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
	trimmed := strings.TrimSpace(dates)
	if strings.HasPrefix(trimmed, "кр") {
		weekString := strings.TrimLeft(trimmed, "кр. ")
		result.Rule.Mode = Any
		weeks := strings.Split(weekString, ",")
		for _, week := range weeks {
			w, err := strconv.Atoi(strings.TrimSpace(week))
			if err != nil {
				log.Panic(err)
			}
			result.Rule.Except = append(result.Rule.Except, w)
		}
		return
	}

	d := strings.Split(dates, ",")
	for _, week := range d {
		week = strings.TrimSpace(week)
		num, err := strconv.Atoi(week)
		if err != nil {
			log.Panic(err)
		}
		result.Rule.Dates = append(result.Rule.Dates, num)
	}

	result.StartWeek = result.Rule.Dates[0]
	result.Rule.Mode = Enum
}
