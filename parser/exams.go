package parser

import (
	"strings"
	"time"

	"github.com/adzhukov/mirea_ics/calendar"
	"github.com/adzhukov/mirea_ics/repeat"

	"github.com/tealeg/xlsx/v3"
)

const (
	stateDone = iota
	stateType = iota
	stateLast = iota
)

const (
	examType = iota
	examTime = iota
	examRoom = iota
)

func parseExams(sheet *xlsx.Sheet, group string, cal *calendar.Calendar) {
	groupColumn, rowNumber := getGroupColumn(sheet, cal.Group), 2

	current := calendar.Event{
		Semester: &cal.Semester,
		Weekday:  time.Sunday,
		Repeat:   repeat.Rule{Mode: repeat.Once},
		Num:      0,
	}

	state := stateDone

	for ; rowNumber < 125; rowNumber++ {
		row, _ := sheet.Row(rowNumber)
		date := row.GetCell(1).Value

		cell := row.GetCell(groupColumn).Value
		if cell == "" {
			continue
		}

		switch cell {
		case "Зачет", "Экзамен", "Консультация":
			current.ClassType = strings.ToUpper(string([]rune(cell)[:3]))
			setExamTime(&current, date, row.GetCell(groupColumn+examTime).Value)
			current.Classroom = row.GetCell(groupColumn + examRoom).Value
			state++
		default:
			if state == stateType {
				current.Subject = cell
				state++
			} else if state == stateLast {
				current.Lecturer = cell
				state = stateDone
				cal.Classes = append(cal.Classes, current)
			}
		}
	}
}
