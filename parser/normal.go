package parser

import (
	"strconv"
	"strings"
	"time"

	"github.com/adzhukov/mirea_ics/calendar"
	"github.com/adzhukov/mirea_ics/repeat"
)

const (
	offsetSubject   = iota
	offsetType      = iota
	offsetLecturer  = iota
	offsetClassroom = iota
)

const (
	columnTime = 2
	columnWeek = 4
)

func (p *Parser) normal() {
	current := calendar.Event{
		Semester: &p.Calendar.Semester,
		Weekday:  time.Sunday,
	}

	var timeCell string

	for rowNumber := 3; rowNumber < 125; rowNumber++ {
		row, _ := p.Sheet.Row(rowNumber)

		weekType := row.GetCell(columnWeek).Value
		if weekType == "" {
			break
		}

		if weekType == "I" {
			timeCell = row.GetCell(columnTime).Value
			current.WeekType = calendar.Odd
			current.Num, _ = strconv.Atoi(row.GetCell(1).Value)
			if current.Num == 1 {
				current.Weekday++
			}
		} else {
			current.WeekType = calendar.Even
		}

		subjects := strings.Split(row.GetCell(p.Column+offsetSubject).Value, "\n")
		classType := strings.Split(row.GetCell(p.Column+offsetType).Value, "\n")
		lecturer := strings.Split(row.GetCell(p.Column+offsetLecturer).Value, "\n")
		classroom := strings.Split(row.GetCell(p.Column+offsetClassroom).Value, "\n")

		for i, subject := range subjects {
			if subject == "" {
				continue
			}

			current.ClassType = classType[0]
			if len(subjects) == 1 {
				current.ClassType = strings.Join(classType, "")
			} else if i < len(classType) {
				current.ClassType = classType[i]
			}

			current.Lecturer = lecturer[0]
			if len(subjects) == 1 {
				current.Lecturer = strings.Join(lecturer, ", ")
			} else if i < len(lecturer) {
				current.Lecturer = lecturer[i]
			}

			current.Classroom = classroom[0]
			if len(subjects) == 1 {
				current.Classroom = strings.Join(classroom, " ")
			} else if i < len(classroom) {
				current.Classroom = classroom[i]
			}

			parsed := repeat.Parse(subject)
			current.Repeat = parsed.Rule
			current.Subject = parsed.Subject

			setEventTime(&current, timeCell, parsed.StartWeek)
			p.Calendar.Classes = append(p.Calendar.Classes, current)
		}
	}
}
