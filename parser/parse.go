package parser

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/adzhukov/mirea_ics/calendar"

	"github.com/tealeg/xlsx/v3"
)

func getGroupColumn(sheet *xlsx.Sheet, group string) int {
	row, _ := sheet.Row(1)
	column := 0
	row.ForEachCell(func(cell *xlsx.Cell) error {
		if group == cell.String() {
			column, _ = cell.GetCoordinates()
			return errors.New("Success")
		}
		return nil
	})
	return column
}

func semesterLength(group string) int {
	switch []rune(group)[2] {
	case 'Б':
		return 16
	case 'М':
		return 17
	}

	return 17
}

func parseSemesterInfo(title string, s *calendar.Semester) {
	if strings.Contains(title, "осеннего") {
		s.Type = calendar.Autumn
	} else if strings.Contains(title, "зимней") {
		s.Type = calendar.Winter
	} else if strings.Contains(title, "весеннего") {
		s.Type = calendar.Spring
	} else if strings.Contains(title, "летней") {
		s.Type = calendar.Summer
	} else {
		log.Println("Unable to parse semester type")
	}
	s.Year = time.Now().Year()

	splitted := strings.Split(title, "-")
	splitted = strings.Fields(splitted[len(splitted)-1])
	year, err := strconv.Atoi(splitted[0])
	if err != nil {
		log.Println("Unable to parse year")
		return
	}

	s.Year = year
	if s.Type == calendar.Autumn {
		s.Year--
	}

	s.Start = semesterStart(s.Year, s.Type)
}

func ParseFile(file string, g string) {
	wb, err := xlsx.OpenFile(file, xlsx.RowLimit(125))
	if err != nil {
		log.Fatal(err)
	}

	cal := calendar.Calendar{
		Group: g,
	}

	sheet := wb.Sheets[0]
	row, _ := sheet.Row(0)

	title := ""
	for i := 0; title == ""; i++ {
		title = row.GetCell(i).String()
	}

	parseSemesterInfo(title, &cal.Semester)
	cal.Semester.End = semesterEnd(semesterLength(g), cal.Semester.Start)

	switch cal.Semester.Type {
	case calendar.Autumn, calendar.Spring:
		parseNormal(sheet, g, &cal)
	case calendar.Winter, calendar.Summer:
		parseExams(sheet, g, &cal)
	}
	cal.WriteToFile()
}
