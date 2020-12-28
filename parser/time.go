package parser

import (
	"strconv"
	"strings"
	"time"
)

var semesterStartDate time.Time
var semesterEndDate time.Time

func semesterStart(autumn bool) time.Time {
	month := time.February
	day := 13
	if autumn {
		month = time.September
		day = 1
	}

	return time.Date(time.Now().Year(), month, day, 0, 0, 0, 0, time.UTC)
}

func semesterEnd() time.Time {
	return semesterStartDate.AddDate(0, 0, 17*7)
}

func (event *class) setEventTime(cellValue string) {
	splitted := strings.Split(cellValue, "-")
	hours, _ := strconv.Atoi(splitted[0])
	minutes, _ := strconv.Atoi(splitted[1])

	localTime := time.Minute*time.Duration(minutes) + time.Hour*time.Duration(hours)

	daysToMonday := int(time.Monday - semesterStartDate.Weekday() + event.weekday - 1)
	startDate := semesterStartDate.AddDate(0, 0, daysToMonday)
	if !event.parity {
		startDate = startDate.AddDate(0, 0, 7)
	}
	event.startTime = startDate.Add(localTime)
}

func (event *class) endTime() time.Time {
	return event.startTime.Add(time.Minute * 90)
}

func (event *class) startAtWeek(n int) {
	event.startTime = event.startTime.AddDate(0, 0, 7*(n-1))
}
