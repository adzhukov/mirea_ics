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

type Parser struct {
	Sheet    *xlsx.Sheet
	Calendar *calendar.Calendar
	Column   int
}

const (
	rowTitle  = iota
	rowGroups = iota
)

func (p Parser) findGroup() {
	row, err := p.Sheet.Row(rowGroups)
	if err != nil {
		log.Fatalln("Unable to parse file", err)
	}

	row.ForEachCell(func(cell *xlsx.Cell) error {
		if p.Calendar.Group == cell.String() {
			p.Column, _ = cell.GetCoordinates()
			return errors.New("Success")
		}
		return nil
	})
}

func (p Parser) parseSemesterInfo() {
	row, err := p.Sheet.Row(rowTitle)
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

	s := &p.Calendar.Semester

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

	length := 17
	if []rune(p.Calendar.Group)[2] == 'Б' {
		length = 16
	}

	s.Start = semesterStart(s.Year, s.Type)
	s.End = semesterEnd(length, s.Start)

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

	p := Parser{
		Calendar: &calendar.Calendar{Group: g},
		Sheet:    wb.Sheets[0],
	}

	p.findGroup()
	if p.Column == 0 {
		log.Printf("Could not find group %s in file %s\n", g, uri)
		return
	}

	p.parseSemesterInfo()

	switch p.Calendar.Semester.Type {
	case calendar.Autumn, calendar.Spring:
		p.normal()
	case calendar.Winter, calendar.Summer:
		p.exams()
	}

	p.Calendar.File()
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

func pattern(group string) string {
	parts := []string{`http(s)?:\/\/.*`}

	switch []rune(group)[0] {
	case 'И':
		parts = append(parts, `ИИТ.*`)
	case 'К':
		parts = append(parts, `ИК.*`)
	}

	switch []rune(group)[2] {
	case 'М':
		parts = append(parts, `маг.*`)
	case 'Б':
		parts = append(parts, `[^маг]{3}.*`)
	}

	parts = append(parts, `\.xlsx`)

	return strings.Join(parts, "")
}

func GetLinks(group string) []string {
	resp, err := http.Get("https://mirea.ru/schedule")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	re := regexp.MustCompile(pattern(group))
	return re.FindAllString(string(body), -1)
}
