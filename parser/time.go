package parser

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/adzhukov/mirea_ics/calendar"
)

const semLength = 17

func semesterStart(year int, semester int) time.Time {
	var month time.Month
	var day int

	switch semester {
	case calendar.Autumn:
		month, day = time.September, 1
	case calendar.Winter:
		month, day = time.January, 1
	case calendar.Spring:
		month, day = time.February, 13
	case calendar.Summer:
		month, day = time.June, 1
	}

	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func semesterEnd(start time.Time) time.Time {
	return start.AddDate(0, 0, semLength*7)
}

func setExamTime(event *calendar.Event, dateCell string, timeCell string) {
	timeSplit := strings.Split(timeCell, "-")
	hours, _ := strconv.Atoi(timeSplit[0])
	minutes, _ := strconv.Atoi(timeSplit[1])

	localTime := time.Minute*time.Duration(minutes) + time.Hour*time.Duration(hours)

	days := strings.Fields(dateCell)
	day, err := strconv.Atoi(days[0])
	if err != nil {
		log.Println("Unable to parse date", err)
	}

	startDate := event.Semester.Start.AddDate(0, 0, day-1)
	event.StartTime = startDate.Add(localTime)
}

func setEventTime(event *calendar.Event, timeCell string, start int) {
	daysToEvent := int(time.Monday - event.Semester.Start.Weekday() + event.Weekday - 1)
	startDate := event.Semester.Start.AddDate(0, 0, daysToEvent)

	if start != 0 {
		startDate = startDate.AddDate(0, 0, 7*(start-1))
	} else if event.WeekType == calendar.Even {
		startDate = startDate.AddDate(0, 0, 7)
	} else if startDate.Month() < event.Semester.Start.Month() {
		startDate = startDate.AddDate(0, 0, 14)
	}

	splitted := strings.Split(timeCell, "-")
	hours, _ := strconv.Atoi(splitted[0])
	minutes, _ := strconv.Atoi(splitted[1])

	localTime := time.Minute*time.Duration(minutes) + time.Hour*time.Duration(hours)
	event.StartTime = startDate.Add(localTime)
}
