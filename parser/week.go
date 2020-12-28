package parser

import (
	"log"
	"strconv"
	"strings"
)

type repeatRule struct {
	mode  int
	dates []int
}

const (
	repeatOnce  = iota
	repeatRange = iota
	repeatAny   = iota
)

func (event *class) parseDates() {
	event.repeat.mode = repeatAny
	splitted := strings.SplitN(event.subject, "н.", 2)

	if len(splitted) < 2 {
		return
	}

	event.subject = splitted[1]
	dates := strings.Split(splitted[0], "-")
	switch len(dates) {
	case 1:
		parseDatesAsSinge(event, dates)
	case 2:
		parseDatesAsRange(event, dates)
	default:
		log.Println("Unable to parse", dates)
	}
}

func parseDatesAsSinge(event *class, dates []string) {
	week, err := strconv.Atoi(strings.TrimSpace(dates[0]))
	if err != nil {
		log.Panic(err)
	}
	event.repeat.mode = repeatOnce
	event.repeat.dates = []int{week}
	event.startAtWeek(week)
}

func parseDatesAsRange(event *class, dates []string) {
	start, err := strconv.Atoi(strings.TrimSpace(dates[0]))
	if err != nil {
		log.Panic(err)
	}

	end, err := strconv.Atoi(strings.TrimSpace(dates[1]))
	if err != nil {
		log.Panic(err)
	}

	event.repeat.mode = repeatRange
	event.repeat.dates = []int{start, end}
	event.startAtWeek(start)
}
