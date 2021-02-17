package calendar

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Calendar struct {
	Semester Semester
	Group    string
	Classes  []Event
}

type SemesterType int

const (
	_      SemesterType = iota
	Autumn SemesterType = iota
	Winter SemesterType = iota
	Spring SemesterType = iota
	Summer SemesterType = iota
)

type Semester struct {
	Year  int
	Num   int
	Type  SemesterType
	Start time.Time
	End   time.Time
}

func (cal Calendar) String() string {
	var sb strings.Builder
	cal.WriteTo(&sb)
	return sb.String()
}

func (cal Calendar) File() {
	file, err := os.Create(cal.name() + ".ics")
	if err != nil {
		panic(err)
	}

	cal.WriteTo(file)
	file.Close()
}

func (cal Calendar) WriteTo(w io.Writer) {
	cal.writeHeader(w)
	writeVTimezone(w)
	for _, event := range cal.Classes {
		event.WriteTo(w)
	}
	writeFooter(w)
}

func (cal Calendar) name() string {
	if cal.Semester.Type == 0 {
		return cal.Group
	}

	t := "Семестр"
	if cal.Semester.Type == Winter || cal.Semester.Type == Summer {
		t = "Сессия"
	}
	return fmt.Sprintf("%s %d %s", cal.Group, cal.Semester.Num, t)
}

func (cal Calendar) color() string {
	switch cal.Semester.Type {
	case Autumn, Spring:
		return "86F79B"
	case Winter, Summer:
		return "E23E24"
	}

	return "FFFFFF"
}

func (cal Calendar) writeHeader(w io.Writer) {
	write(w, "BEGIN:VCALENDAR")
	write(w, "METHOD:PUBLISH")
	write(w, "VERSION:2.0")
	writeLong(w, "X-WR-CALNAME:", cal.name())
	write(w, "PRODID:-//Apple Inc.//Mac OS X 10.15.5//EN")
	write(w, "X-APPLE-CALENDAR-COLOR:#", cal.color())
	write(w, "X-WR-TIMEZONE:Europe/Moscow")
	write(w, "CALSCALE:GREGORIAN")
}

func writeVTimezone(w io.Writer) {
	write(w, "BEGIN:VTIMEZONE")
	write(w, "TZID:Europe/Moscow")
	write(w, "BEGIN:STANDARD")
	write(w, "TZOFFSETFROM:+023017")
	write(w, "DTSTART:20010101T000000")
	write(w, "TZNAME:GMT+3")
	write(w, "TZOFFSETTO:+023017")
	write(w, "END:STANDARD")
	write(w, "END:VTIMEZONE")
}

func writeFooter(w io.Writer) {
	write(w, "END:VCALENDAR")
}
