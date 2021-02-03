package parser

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/adzhukov/mirea_ics/calendar"
)

func semesterStart(year int, semester calendar.SemesterType) time.Time {
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

func semesterEnd(length int, start time.Time) time.Time {
	return start.AddDate(0, 0, length*7)
}

func cellToTime(cell string) (int, int) {
	splitted := strings.Split(cell, "-")
	hours, err := strconv.Atoi(splitted[0])
	if err != nil {
		log.Fatalln("Unable to parse time", err)
	}

	minutes, err := strconv.Atoi(splitted[1])
	if err != nil {
		log.Fatalln("Unable to parse time", err)
	}

	return hours, minutes
}

func setExamTime(event *calendar.Event, dateCell string, timeCell string) {
	hours, minutes := cellToTime(timeCell)
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

	hours, minutes := cellToTime(timeCell)

	localTime := time.Minute*time.Duration(minutes) + time.Hour*time.Duration(hours)
	event.StartTime = startDate.Add(localTime)
}
