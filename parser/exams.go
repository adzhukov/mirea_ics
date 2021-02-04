package parser

import (
	"strings"

	"github.com/adzhukov/mirea_ics/calendar"
	"github.com/adzhukov/mirea_ics/repeat"
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

func (p *Parser) exams() {
	current := calendar.Event{
		Semester: &p.Calendar.Semester,
		Repeat:   repeat.Rule{Mode: repeat.Once},
		Num:      0,
	}

	state := stateDone

	for rowNumber := 2; rowNumber < 125; rowNumber++ {
		row, _ := p.Sheet.Row(rowNumber)
		date := row.GetCell(1).Value
		cell := row.GetCell(p.Column).Value

		switch cell {
		case "":
			continue
		case "Зачет", "Экзамен", "Консультация":
			current.ClassType = strings.ToUpper(string([]rune(cell)[:3]))
			setExamTime(&current, date, row.GetCell(p.Column+examTime).Value)
			current.Classroom = row.GetCell(p.Column + examRoom).Value
			state++
		default:
			if state == stateType {
				current.Subject = cell
				state++
			} else if state == stateLast {
				current.Lecturer = cell
				state = stateDone
				p.Calendar.Classes = append(p.Calendar.Classes, current)
			}
		}
	}
}
