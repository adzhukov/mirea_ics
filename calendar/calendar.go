package calendar

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

type Calendar struct {
	Semester Semester
	Group    string
	Classes  []Event
}

type Semester struct {
	Type  bool
	Start time.Time
	End   time.Time
}

func (cal Calendar) String() string {
	var sb strings.Builder
	cal.writeHeader(&sb)
	writeVTimezone(&sb)
	for _, event := range cal.Classes {
		sb.WriteString(event.String())
	}
	writeFooter(&sb)
	return sb.String()
}

func (cal Calendar) WriteToFile() {
	ioutil.WriteFile(cal.filename(), []byte(cal.String()), 0644)
}

func (cal Calendar) filename() string {
	return cal.Group + ".ics"
}

func (cal Calendar) writeHeader(w io.Writer) {
	write(w, "BEGIN:VCALENDAR")
	write(w, "METHOD:PUBLISH")
	write(w, "VERSION:2.0")
	writeLong(w, "X-WR-CALNAME:", cal.Group)
	write(w, "PRODID:-//Apple Inc.//Mac OS X 10.15.5//EN")
	write(w, "X-APPLE-CALENDAR-COLOR:#FFCC00")
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
	fmt.Fprintln(w, "END:VCALENDAR")
}
