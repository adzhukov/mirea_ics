package parser

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/google/uuid"
)

func writeHeader(sb *strings.Builder, name string) {
	sb.WriteString("BEGIN:VCALENDAR\n")
	sb.WriteString("METHOD:PUBLISH\n")
	sb.WriteString("VERSION:2.0\n")
	sb.WriteString("X-WR-CALNAME:")
	sb.WriteString(name)
	sb.WriteString("\n")
	sb.WriteString("PRODID:-//Apple Inc.//Mac OS X 10.15.5//EN\n")
	sb.WriteString("X-APPLE-CALENDAR-COLOR:#FFCC00\n")
	sb.WriteString("X-WR-TIMEZONE:Europe/Moscow\n")
	sb.WriteString("CALSCALE:GREGORIAN\n")
}

func writeFooter(sb *strings.Builder) {
	sb.WriteString("END:VCALENDAR")
}

func writeVTimezone(sb *strings.Builder) {
	sb.WriteString("BEGIN:VTIMEZONE\n")
	sb.WriteString("TZID:Europe/Moscow\n")
	sb.WriteString("BEGIN:STANDARD\n")
	sb.WriteString("TZOFFSETFROM:+023017\n")
	sb.WriteString("DTSTART:20010101T000000\n")
	sb.WriteString("TZNAME:GMT+3\n")
	sb.WriteString("TZOFFSETTO:+023017\n")
	sb.WriteString("END:STANDARD\n")
	sb.WriteString("END:VTIMEZONE\n")
}

func writeAppleLocation(sb *strings.Builder, title *string) {
	sb.WriteString("X-APPLE-STRUCTURED-LOCATION;VALUE=URI;X-ADDRESS=Vernadskogo prospekt 78\n")
	sb.WriteString(" \\nMoscow\\nMoscow\\nRussia\\n119415;X-APPLE-MAPKIT-HANDLE=CAESxwEaEgkgE7\n")
	sb.WriteString(" ngwtVLQBGy+0hSe71CQCJfCgZSdXNzaWESAlJVGgZNb3Njb3cqBk1vc2NvdzIGTW9zY293Og\n")
	sb.WriteString(" YxMTk0MTVSFFZlcm5hZHNrb2dvIHByb3NwZWt0WgI3OGIXVmVybmFkc2tvZ28gcHJvc3Bla3\n")
	sb.WriteString(" QgNzgqH01vc2NvdyBUZWNobm9sb2dpY2FsIFVuaXZlcnNpdHkyF1Zlcm5hZHNrb2dvIHByb3\n")
	sb.WriteString(" NwZWt0IDc4MgZNb3Njb3cyBlJ1c3NpYTIGMTE5NDE1;X-APPLE-REFERENCEFRAME=0;X-TI\n")
	sb.WriteString(" TLE=")
	if title != nil {
		sb.WriteString(*title)
	} else {
		sb.WriteString("Moscow Technological University")
	}
	sb.WriteString("::geo:55.670010,37.480326\n")
}

func writeLocation(sb *strings.Builder, title *string) {
	sb.WriteString("LOCATION:")
	if title != nil {
		sb.WriteString(*title)
	} else {
		sb.WriteString("RTU MIREA")
	}
	sb.WriteString("\\nVernadskogo prospekt 78\\nMoscow\\nMoscow\\nRussia\\n119415\n")
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

func writeRepeatRule(sb *strings.Builder, event *class) {
	var endDate time.Time

	switch event.repeat.mode {
	case repeatOnce:
		return
	case repeatRange:
		endDate = semesterStartDate.AddDate(0, 0, 7*event.repeat.dates[1])
	case repeatAny:
		endDate = semesterEndDate
	}

	sb.WriteString("RRULE:FREQ=WEEKLY;")
	sb.WriteString("INTERVAL=2;")
	sb.WriteString("UNTIL=")
	sb.WriteString(endDate.UTC().Format("20060102T150405"))
	sb.WriteString(";BYDAY=")
	sb.WriteString(event.byday())
	sb.WriteString(";WKST=SU;\n")
}

func limitLineLength(s string, limit int) string {
	if limit >= len(s) {
		return s
	}
	var chunks []string
	chunk := make([]rune, limit)
	len := 0
	for _, r := range s {
		chunk[len] = r
		len++
		if len == limit {
			chunks = append(chunks, string(chunk))
			chunk[0] = ' '
			len = 1
		}
	}
	if len > 0 {
		chunks = append(chunks, string(chunk[:len]))
	}
	return strings.Join(chunks, "\n")
}

func writeSummary(sb *strings.Builder, event *class) {
	sb.WriteString("SUMMARY:")
	sb.WriteString(limitLineLength(event.classType+" "+event.subject, 74))
	sb.WriteString("\n")
}

func writeEvent(sb *strings.Builder, event class) {
	timeNow := time.Now().UTC().Format(time.RFC3339)

	sb.WriteString("BEGIN:VEVENT\n")
	sb.WriteString("TRANSP:OPAQUE\n")

	sb.WriteString("DTSTART;TZID=Europe/Moscow:")
	sb.WriteString(event.startTime.Format("20060102T150405"))
	sb.WriteString("\n")
	sb.WriteString("DTEND;TZID=Europe/Moscow:")
	sb.WriteString(event.endTime().Format("20060102T150405"))
	sb.WriteString("\n")

	sb.WriteString("UID:")
	sb.WriteString(uuid.New().String())
	sb.WriteString("\n")

	writeRepeatRule(sb, &event)
	writeAppleLocation(sb, &event.classroom)
	writeLocation(sb, &event.classroom)
	sb.WriteString("X-APPLE-TRAVEL-ADVISORY-BEHAVIOR:DISABLED\n")
	sb.WriteString("SEQUENCE:0\n")

	writeSummary(sb, &event)
	sb.WriteString("DESCRIPTION:")
	sb.WriteString(event.lecturer)
	sb.WriteString("\n")

	sb.WriteString("DTSTAMP:")
	sb.WriteString(timeNow)
	sb.WriteString("\n")
	sb.WriteString("CREATED:")
	sb.WriteString(timeNow)
	sb.WriteString("\n")
	sb.WriteString("LAST-MODIFIED:")
	sb.WriteString(timeNow)
	sb.WriteString("\n")

	sb.WriteString("END:VEVENT\n")
}

func writeToICS(table timeTable) {
	var sb strings.Builder
	writeHeader(&sb, table.group)
	writeVTimezone(&sb)
	for _, event := range table.classes {
		writeEvent(&sb, event)
	}
	writeFooter(&sb)
	ioutil.WriteFile(table.group+".ics", []byte(sb.String()), 0644)
}
