package parser

import (
	"log"
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
	}

	state := stateDone

	colDate := 0
	for row, _ := p.Sheet.Row(1); row.GetCell(colDate).Value != "число"; {
		if colDate++; colDate == 10 {
			log.Fatalln("Unable to find date column")
		}
	}

	for rowNumber := 2; rowNumber < 125; rowNumber++ {
		row, _ := p.Sheet.Row(rowNumber)
		date := row.GetCell(colDate).Value
		cell := row.GetCell(p.Column).Value

		switch cell {
		case "":
			continue
		case "Зачет", "Экзамен", "Консультация":
			current.ClassType = strings.ToUpper(string([]rune(cell)[:3]))
			setExamTime(&current, date, row.GetCell(p.Column+examTime).Value)
			if current.StartTime.Day() > 21 && rowNumber < 10 {
				current.StartTime = current.StartTime.AddDate(0, -1, 0)
			}
			current.Classroom = row.GetCell(p.Column + examRoom).Value
			state++
		default:
			if state == stateType {
				current.Subject = cell
				state++
			} else if state == stateLast {
				current.Lecturer = strings.ReplaceAll(cell, "\n", " ")
				state = stateDone
				p.Calendar.Classes = append(p.Calendar.Classes, current)
			}
		}
	}
}
