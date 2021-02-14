package parser

import (
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
	}

	var timeCell string
	weekday := time.Sunday

	for rowNumber := 3; rowNumber < 125; rowNumber++ {
		row, _ := p.Sheet.Row(rowNumber)

		weekType := row.GetCell(columnWeek).Value
		if weekType == "" {
			break
		}

		if weekType == "I" {
			timeCell = row.GetCell(columnTime).Value
			current.WeekType = calendar.Odd
			if strings.TrimSpace(row.GetCell(1).Value) == "1" {
				weekday++
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

			setEventTime(&current, timeCell, parsed.StartWeek, weekday)

			if isValidEvent(&current) {
				p.Calendar.Classes = append(p.Calendar.Classes, current)
			}
		}
	}
}

func isValidEvent(e *calendar.Event) bool {
	if strings.HasPrefix(e.Subject, "……") {
		return false
	}

	if strings.ToUpper(e.ClassType) == "С/Р" {
		return false
	}

	switch strings.TrimSpace(e.Subject) {
	case "Военная подготовка":
		return false
	case "Военная", "подготовка":
		return false
	case "День", "самостоятельных", "занятий":
		return false
	}

	return true
}
