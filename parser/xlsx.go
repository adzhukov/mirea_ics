package parser

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx/v3"
)

type timeTable struct {
	classes      []class
	group        string
	semesterType bool
}

type class struct {
	subject   string
	classType string
	classroom string
	startTime time.Time
	lecturer  string
	parity    bool
	num       int
	repeat    repeatRule
	weekday   time.Weekday
}

const (
	offsetSubject   = iota
	offsetType      = iota
	offsetLecturer  = iota
	offsetClassroom = iota
)

var table timeTable

func setSemesterType(sheet *xlsx.Sheet) {
	cell, _ := sheet.Cell(0, 0)
	table.semesterType = strings.Contains(cell.String(), "осеннего")
	semesterStartDate = semesterStart(table.semesterType)
	semesterEndDate = semesterEnd()
}

func getGroupColumn(sheet *xlsx.Sheet) int {
	row, _ := sheet.Row(1)
	column := 0
	row.ForEachCell(func(cell *xlsx.Cell) error {
		if table.group == cell.String() {
			column, _ = cell.GetCoordinates()
			return errors.New("Success")
		}
		return nil
	})
	return column
}

func parse(sheet *xlsx.Sheet) {
	setSemesterType(sheet)
	groupColumn, rowNumber := getGroupColumn(sheet), 3
	var currentClass class
	currentClass.weekday = time.Sunday
	for {
		row, _ := sheet.Row(rowNumber)
		weekType := row.GetCell(4).Value

		if weekType == "" {
			break
		}

		isOddWeek := weekType == "I"

		if isOddWeek {
			currentClass.num, _ = strconv.Atoi(row.GetCell(1).Value)
			if currentClass.num == 1 {
				currentClass.weekday++
			}
			currentClass.setEventTime(row.GetCell(2).Value)
		}

		currentClass.subject = row.GetCell(groupColumn + offsetSubject).Value

		if currentClass.subject != "" {
			currentClass.classType = row.GetCell(groupColumn + offsetType).Value
			currentClass.lecturer = row.GetCell(groupColumn + offsetLecturer).Value
			currentClass.classroom = row.GetCell(groupColumn + offsetClassroom).Value
			currentClass.parity = isOddWeek
			currentClass.parseDates()

			table.classes = append(table.classes, currentClass)
		}

		rowNumber++
	}
}

func ParseFile(file string, group string) {
	wb, err := xlsx.OpenFile(file, xlsx.RowLimit(125))

	if err != nil {
		log.Fatal(err)
	}

	table.group = group

	parse(wb.Sheets[0])
	writeToICS(table)
}
