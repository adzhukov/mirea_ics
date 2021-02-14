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
		month, day = time.February, 11
	case calendar.Summer:
		month, day = time.June, 1
	}

	date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	switch date.Weekday() {
	case time.Saturday:
		date = date.AddDate(0, 0, 2)
	case time.Sunday:
		date = date.AddDate(0, 0, 1)
	}

	return date
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

func setEventTime(event *calendar.Event, timeCell string, start int, weekday time.Weekday) {
	daysToEvent := int(time.Monday - event.Semester.Start.Weekday() + weekday - 1)
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

func groupYear(group []rune) int {
	y, err := strconv.Atoi(string(group[len(group)-2:]))
	if err != nil {
		return 0
	}

	now := time.Now()
	year := now.Year() - y - 1999
	if now.Month() < time.August {
		year--
	}

	return year
}
