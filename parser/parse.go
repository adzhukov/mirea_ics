package parser

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/adzhukov/mirea_ics/calendar"

	"github.com/tealeg/xlsx/v3"
)

const (
	rowTitle  = iota
	rowGroups = iota
)

func getGroupColumn(sheet *xlsx.Sheet, group string) int {
	row, _ := sheet.Row(rowGroups)
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
	if []rune(group)[2] == 'Б' {
		return 16
	}

	return 17
}

func parseSemesterInfo(title string, s *calendar.Semester) {
	switch {
	case strings.Contains(title, "осеннего"):
		s.Type = calendar.Autumn
	case strings.Contains(title, "зимней"):
		s.Type = calendar.Winter
	case strings.Contains(title, "весеннего"):
		s.Type = calendar.Spring
	case strings.Contains(title, "летней"):
		s.Type = calendar.Summer
	default:
		log.Fatalln("Unable to parse semester type:", title)
	}

	s.Year = time.Now().Year()

	splitted := strings.Split(title, "-")
	splitted = strings.Fields(splitted[len(splitted)-1])
	year, err := strconv.Atoi(splitted[0])
	if err != nil {
		log.Println("Unable to parse year", err)
	} else {
		s.Year = year
	}

	if s.Type == calendar.Autumn {
		s.Year--
	}

	s.Start = semesterStart(s.Year, s.Type)

	splitted = strings.Fields(title)
	for i := range splitted {
		if splitted[i] == "курса" {
			n, err := strconv.Atoi(splitted[i-1])
			if err != nil {
				log.Println(err)
			}
			s.Num = n
			break
		}
	}
}

func openFile(uri string) (*xlsx.File, error) {
	if !strings.HasPrefix(uri, "http") {
		return xlsx.OpenFile(uri, xlsx.RowLimit(125))
	}

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return xlsx.OpenBinary(body, xlsx.RowLimit(125))
}

func Parse(uri string, g string) {
	wb, err := openFile(uri)
	if err != nil {
		log.Fatalln(err)
	}

	cal := calendar.Calendar{
		Group: g,
	}

	sheet := wb.Sheets[0]

	row, err := sheet.Row(rowTitle)
	if err != nil {
		log.Fatalln(err)
	}

	title := ""
	row.ForEachCell(func(cell *xlsx.Cell) error {
		if cell.String() != "" {
			title = cell.String()
			return errors.New("Success")
		}
		return nil
	})

	parseSemesterInfo(title, &cal.Semester)
	cal.Semester.End = semesterEnd(semesterLength(g), cal.Semester.Start)

	switch cal.Semester.Type {
	case calendar.Autumn, calendar.Spring:
		parseNormal(sheet, g, &cal)
	case calendar.Winter, calendar.Summer:
		parseExams(sheet, g, &cal)
	}

	cal.File()
}

func Groups(file string) []string {
	wb, err := openFile(file)
	if err != nil {
		log.Fatalln(err)
	}

	row, err := wb.Sheets[0].Row(rowGroups)
	if err != nil {
		log.Fatalln(err)
	}

	groups := []string{}

	row.ForEachCell(func(cell *xlsx.Cell) error {
		if strings.Count(cell.Value, "-") == 2 {
			groups = append(groups, cell.Value)
		}

		return nil
	})

	return groups
}

func GetLinks() []string {
	resp, err := http.Get("https://mirea.ru/schedule")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	re := regexp.MustCompile(`http(s)?:\/\/.*\.xlsx`)
	return re.FindAllString(string(body), -1)
}
