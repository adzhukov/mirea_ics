package parser

import (
	"errors"
	"log"
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

func getGroupColumn(sheet *xlsx.Sheet, cal calendar.Calendar) int {
	row, _ := sheet.Row(1)
	column := 0
	row.ForEachCell(func(cell *xlsx.Cell) error {
		if cal.Group == cell.String() {
			column, _ = cell.GetCoordinates()
			return errors.New("Success")
		}
		return nil
	})
	return column
}

func parse(sheet *xlsx.Sheet, group string) calendar.Calendar {
	var cal calendar.Calendar

	cell, _ := sheet.Cell(0, 0)
	cal.Semester.Type = strings.Contains(cell.String(), "осеннего")
	cal.Semester.Start = semesterStart(cal.Semester.Type)
	cal.Semester.End = semesterEnd(cal.Semester.Start)
	cal.Group = group

	groupColumn, rowNumber := getGroupColumn(sheet, cal), 3
	var current calendar.Event
	current.Semester = &cal.Semester
	current.Weekday = time.Sunday

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

	return cal
}

func ParseFile(file string, g string) {
	wb, err := xlsx.OpenFile(file, xlsx.RowLimit(125))
	if err != nil {
		log.Fatal(err)
	}

	calendar := parse(wb.Sheets[0], g)
	calendar.WriteToFile()
}
