package parser

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const defaultLocation = "Moscow Technological University"

const delimiter = "\r\n"
const timeFormat = "20060102T150405"

func (event class) String() string {
	var sb strings.Builder
	writeEvent(&sb, event)
	return sb.String()
}

func (table timeTable) String() string {
	var sb strings.Builder
	writeHeader(&sb, table.group)
	writeVTimezone(&sb)
	for _, event := range table.classes {
		sb.WriteString(event.String())
	}
	writeFooter(&sb)
	return sb.String()
}

func writeToICS(table timeTable) {
	ioutil.WriteFile(table.group+".ics", []byte(table.String()), 0644)
}

func write(w io.Writer, a ...interface{}) {
	fmt.Fprint(w, append(a, delimiter)...)
}

func writeLong(w io.Writer, a ...interface{}) {
	b := bytes.NewBufferString("")
	fmt.Fprint(b, a...)
	r := b.String()
	if len(r) > 75 {
		l := limitLineLength(r, 75)
		fmt.Fprint(w, l, delimiter)
		r = r[len(l):]
		fmt.Fprint(w, " ")
	}
	for len(r) > 74 {
		l := limitLineLength(r, 74)
		fmt.Fprint(w, l, delimiter)
		r = r[len(l):]
		fmt.Fprint(w, " ")
	}
	fmt.Fprint(w, r, delimiter)
}

func limitLineLength(s string, max int) string {
	length := 0
	for _, r := range s {
		newLength := length + utf8.RuneLen(r)
		if newLength > max {
			break
		}
		length = newLength
	}
	return s[:length]
}

func writeHeader(w io.Writer, name string) {
	write(w, "BEGIN:VCALENDAR")
	write(w, "METHOD:PUBLISH")
	write(w, "VERSION:2.0")
	writeLong(w, "X-WR-CALNAME:", name)
	write(w, "PRODID:-//Apple Inc.//Mac OS X 10.15.5//EN")
	write(w, "X-APPLE-CALENDAR-COLOR:#FFCC00")
	write(w, "X-WR-TIMEZONE:Europe/Moscow")
	write(w, "CALSCALE:GREGORIAN")
}

func writeFooter(w io.Writer) {
	fmt.Fprintln(w, "END:VCALENDAR")
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

func writeAppleLocation(w io.Writer, title *string) {
	location := defaultLocation
	if title != nil {
		location = *title
	}

	writeLong(w, "X-APPLE-STRUCTURED-LOCATION;VALUE=URI;X-ADDRESS=Vernadskogo prospekt 78",
		"\\nMoscow\\nMoscow\\nRussia\\n119415;X-APPLE-MAPKIT-HANDLE=CAESxwEaEgkgE7",
		"ngwtVLQBGy+0hSe71CQCJfCgZSdXNzaWESAlJVGgZNb3Njb3cqBk1vc2NvdzIGTW9zY293Og",
		"YxMTk0MTVSFFZlcm5hZHNrb2dvIHByb3NwZWt0WgI3OGIXVmVybmFkc2tvZ28gcHJvc3Bla3",
		"QgNzgqH01vc2NvdyBUZWNobm9sb2dpY2FsIFVuaXZlcnNpdHkyF1Zlcm5hZHNrb2dvIHByb3",
		"NwZWt0IDc4MgZNb3Njb3cyBlJ1c3NpYTIGMTE5NDE1;X-APPLE-REFERENCEFRAME=0;X-TI",
		"TLE=", location, "::geo:55.670010,37.480326")
}

func writeLocation(w io.Writer, title *string) {
	location := defaultLocation
	if title != nil {
		location = *title
	}

	write(w, "LOCATION:", location, "\\nVernadskogo prospekt 78\\nMoscow\\nMoscow\\nRussia\\n119415")
}

func (event *class) byday() string {
	switch event.weekday {
	case time.Monday:
		return "MO"
	case time.Tuesday:
		return "TU"
	case time.Wednesday:
		return "WE"
	case time.Thursday:
		return "TH"
	case time.Friday:
		return "FR"
	case time.Saturday:
		return "SA"
	}
	return "SU"
}

func writeRepeatRule(w io.Writer, event *class) {
	var endDate time.Time

	switch event.repeat.mode {
	case repeatOnce:
		return
	case repeatRange:
		endDate = semesterStartDate.AddDate(0, 0, 7*event.repeat.dates[1])
	case repeatAny:
		endDate = semesterEndDate
	}

	write(w, "RRULE:FREQ=WEEKLY;",
		"INTERVAL=2;",
		"UNTIL=", endDate.UTC().Format("20060102T150405"),
		";BYDAY=", event.byday(), ";WKST=SU;")
}

func writeSummary(w io.Writer, event *class) {
	classType := strings.TrimSpace(event.classType)
	classType = strings.ToUpper(classType)
	subject := strings.TrimSpace(event.subject)
	writeLong(w, "SUMMARY:", classType, " ", subject)
}

func writeEvent(w io.Writer, event class) {
	timeNow := time.Now().UTC().Format(time.RFC3339)

	write(w, "BEGIN:VEVENT")
	write(w, "TRANSP:OPAQUE")

	write(w, "DTSTART;TZID=Europe/Moscow:", event.startTime.Format(timeFormat))
	write(w, "DTEND;TZID=Europe/Moscow:", event.endTime().Format(timeFormat))

	write(w, "UID:", uuid.New().String())

	writeRepeatRule(w, &event)
	writeAppleLocation(w, &event.classroom)
	writeLocation(w, &event.classroom)

	write(w, "X-APPLE-TRAVEL-ADVISORY-BEHAVIOR:DISABLED")
	write(w, "SEQUENCE:0")

	writeSummary(w, &event)

	writeLong(w, "DESCRIPTION:", event.lecturer)

	write(w, "DTSTAMP:", timeNow)
	write(w, "CREATED:", timeNow)
	write(w, "LAST-MODIFIED:", timeNow)

	write(w, "END:VEVENT")
}
