package parser

import (
	"strconv"
	"strings"
	"time"

	"github.com/adzhukov/mirea_ics/calendar"
)

const semLength = 17

func semesterStart(autumn bool) time.Time {
	month, day := time.February, 13
	if autumn {
		month, day = time.September, 1
	}

	return time.Date(time.Now().Year(), month, day, 0, 0, 0, 0, time.UTC)
}

func semesterEnd(start time.Time) time.Time {
	return start.AddDate(0, 0, semLength*7)
}

func setEventTime(event *calendar.Event, cellValue string) {
	splitted := strings.Split(cellValue, "-")
	hours, _ := strconv.Atoi(splitted[0])
	minutes, _ := strconv.Atoi(splitted[1])

	localTime := time.Minute*time.Duration(minutes) + time.Hour*time.Duration(hours)

	daysToMonday := int(time.Monday - event.Semester.Start.Weekday() + event.Weekday - 1)
	startDate := event.Semester.Start.AddDate(0, 0, daysToMonday)
	if !event.Parity {
		startDate = startDate.AddDate(0, 0, 7)
	}
	event.StartTime = startDate.Add(localTime)
}

func startAtWeek(event *calendar.Event, n int) {
	event.StartTime = event.StartTime.AddDate(0, 0, 7*(n-1))
}
