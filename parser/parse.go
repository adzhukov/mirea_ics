package parser

import (
	"errors"
	"fmt"
	"io"
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

func (p *Parser) findGroup() {
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

func (p *Parser) parseSemesterInfo() {
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

	s.Year = time.Now().Year()

	splitted := strings.Split(title, "-")
	splitted = strings.Fields(splitted[len(splitted)-1])
	year, err := strconv.Atoi(splitted[0])
	if err != nil {
		log.Println("Unable to parse year", err)
	} else {
		s.Year = year
	}

	splitted = strings.Fields(title)
	for i := range splitted {
		if splitted[i] == "курса" {
			n, err := strconv.Atoi(splitted[i-1])
			if err != nil {
				log.Println(err)
			}
			s.Num = n * 2
			break
		}
	}

	switch {
	case strings.Contains(title, "осеннего"):
		s.Type = calendar.Autumn
		s.Year--
		s.Num--
	case strings.Contains(title, "зимней"):
		s.Type = calendar.Winter
		s.Num--
	case strings.Contains(title, "весеннего"):
		s.Type = calendar.Spring
	case strings.Contains(title, "летней"):
		s.Type = calendar.Summer
	default:
		log.Fatalln("Unable to parse semester type:", title)
	}

	length := 17
	if []rune(p.Calendar.Group)[2] == 'Б' {
		if s.Num == 8 {
			length = 8
		} else {
			length = 16
		}
	}

	s.Start = semesterStart(s.Year, s.Type)
	s.End = semesterEnd(length, s.Start)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return xlsx.OpenBinary(body, xlsx.RowLimit(125))
}

func parse(file *xlsx.File, g string) {
	p := Parser{
		Calendar: &calendar.Calendar{Group: normalizeGroup(g)},
		Sheet:    file.Sheets[0],
	}

	p.findGroup()
	if p.Column == 0 {
		log.Fatalf("Could not find group %s\n", g)
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

func groups(file *xlsx.File) []string {
	row, err := file.Sheets[0].Row(rowGroups)
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

func filterGroups(body []byte, str string) []string {
	re := regexp.MustCompile(`https?:\/\/.*\.xlsx`)
	links := re.FindAllString(string(body), -1)

	filters := []func(string) bool{}

	group := []rune(str)
	switch group[0] {
	case 'И':
		filters = append(filters, func(s string) bool {
			return strings.Contains(s, `ИИТ`)
		})
	case 'К':
		filters = append(filters, func(s string) bool {
			return strings.Contains(s, `ИК`)
		})
	}

	switch group[2] {
	case 'М':
		filters = append(filters, func(s string) bool {
			return strings.Contains(s, `маг`)
		})
	case 'Б':
		filters = append(filters, func(s string) bool {
			return !strings.Contains(s, `маг`)
		})
	}

	if year := groupYear(group); year != 0 {
		filters = append(filters, func(year int) func(string) bool {
			re := regexp.MustCompile(fmt.Sprintf(`%d(?:_-\s)?к`, year))
			return func(s string) bool {
				return re.MatchString(s)
			}
		}(year))
	}

	results := []string{}
	for _, link := range links {
		valid := true
		for _, filter := range filters {
			if !filter(link) {
				valid = false
				break
			}
		}

		if valid {
			results = append(results, link)
		}
	}

	return results
}

func normalizeGroup(group string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case 'i', 'I':
			return 'И'
		case 'k', 'K':
			return 'К'
		case 'm', 'M':
			return 'М'
		case 'b', 'B':
			return 'Б'
		case 'o', 'O':
			return 'О'
		case 'r', 'R':
			return 'Р'
		case 'v', 'V':
			return 'В'
		case 'n', 'N':
			return 'Н'
		case 'a', 'A':
			return 'А'
		default:
			return r
		}
	}, group)
}
