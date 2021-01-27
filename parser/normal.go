package parser

import (
	"strconv"
	"strings"
	"time"

	"github.com/adzhukov/mirea_ics/calendar"
	"github.com/adzhukov/mirea_ics/repeat"

	"github.com/tealeg/xlsx/v3"
)

const (
	offsetSubject   = iota
	offsetType      = iota
	offsetLecturer  = iota
	offsetClassroom = iota
)

const timeColumn = 2

func parseNormal(sheet *xlsx.Sheet, group string, cal *calendar.Calendar) {
	base, rowNumber := getGroupColumn(sheet, cal.Group), 3

	current := calendar.Event{
		Semester: &cal.Semester,
		Weekday:  time.Sunday,
	}

	var timeCell string

	for {
		row, _ := sheet.Row(rowNumber)

		weekType := row.GetCell(4).Value
		if weekType == "" {
			break
		}

		if weekType == "I" {
			timeCell = row.GetCell(timeColumn).Value
			current.WeekType = calendar.Odd
			current.Num, _ = strconv.Atoi(row.GetCell(1).Value)
			if current.Num == 1 {
				current.Weekday++
			}
		} else {
			current.WeekType = calendar.Even
		}

		subjectValue := row.GetCell(base + offsetSubject).Value
		classType := strings.Split(row.GetCell(base+offsetType).Value, "\n")
		lecturer := strings.Split(row.GetCell(base+offsetLecturer).Value, "\n")
		classroom := strings.Split(row.GetCell(base+offsetClassroom).Value, "\n")

		for i, subject := range strings.Split(subjectValue, "\n") {
			if subject == "" {
				continue
			}

			current.Subject = subject
			current.ClassType = classType[0]
			if i < len(classType) {
				current.ClassType = classType[i]
			}

			current.Lecturer = lecturer[0]
			if i < len(lecturer) {
				current.Lecturer = lecturer[i]
			}

			current.Classroom = classroom[0]
			if i < len(classroom) {
				current.Classroom = classroom[i]
			}

			parsed := repeat.Parse(current.Subject)
			current.Repeat = parsed.Rule
			if parsed.Subject != "" {
				current.Subject = parsed.Subject
			}

			setEventTime(&current, timeCell, parsed.StartWeek)
			cal.Classes = append(cal.Classes, current)
		}

		rowNumber++
	}
}
