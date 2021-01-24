package parser

import (
	"strconv"
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

func parseNormal(sheet *xlsx.Sheet, group string, cal *calendar.Calendar) {
	groupColumn, rowNumber := getGroupColumn(sheet, cal.Group), 3

	current := calendar.Event{
		Semester: &cal.Semester,
		Weekday:  time.Sunday,
	}

	for {
		row, _ := sheet.Row(rowNumber)

		weekType := row.GetCell(4).Value
		if weekType == "" {
			break
		}

		isOddWeek := weekType == "I"

		if isOddWeek {
			current.Num, _ = strconv.Atoi(row.GetCell(1).Value)
			if current.Num == 1 {
				current.Weekday++
			}
			setEventTime(&current, row.GetCell(2).Value)
		}

		current.Subject = row.GetCell(groupColumn + offsetSubject).Value

		if current.Subject != "" {
			current.ClassType = row.GetCell(groupColumn + offsetType).Value
			current.Lecturer = row.GetCell(groupColumn + offsetLecturer).Value
			current.Classroom = row.GetCell(groupColumn + offsetClassroom).Value
			current.Parity = isOddWeek

			parsed := repeat.Parse(current.Subject)
			current.Repeat = parsed.Rule
			if parsed.Subject != "" {
				current.Subject = parsed.Subject
			}

			if parsed.StartWeek != 0 {
				startAtWeek(&current, parsed.StartWeek)
			}

			cal.Classes = append(cal.Classes, current)
		}

		rowNumber++
	}
}
