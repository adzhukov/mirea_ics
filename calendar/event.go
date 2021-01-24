package calendar

import (
	"io"
	"strings"
	"time"

	"github.com/adzhukov/mirea_ics/repeat"
	"github.com/google/uuid"
)

const defaultLocation = "Moscow Technological University"
const timeFormat = "20060102T150405"

type Event struct {
	Subject   string
	ClassType string
	Classroom string
	StartTime time.Time
	Lecturer  string
	Parity    bool
	Num       int
	Repeat    repeat.Rule
	Weekday   time.Weekday
	Semester  *Semester
}

func (event Event) String() string {
	var sb strings.Builder
	writeEvent(&sb, event)
	return sb.String()
}

func (event *Event) byday() string {
	switch event.Weekday {
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

func writeRepeatRule(w io.Writer, event *Event) {
	var endDate time.Time

	switch event.Repeat.Mode {
	case repeat.Once:
		return
	case repeat.Range:
		endDate = event.Semester.Start.AddDate(0, 0, 7*event.Repeat.Dates[1])
	case repeat.Any:
		endDate = event.Semester.End
	}

	write(w, "RRULE:FREQ=WEEKLY;",
		"INTERVAL=2;",
		"UNTIL=", endDate.UTC().Format("20060102T150405"),
		";BYDAY=", event.byday(), ";WKST=SU;")
}

func writeSummary(w io.Writer, event *Event) {
	classType := strings.TrimSpace(event.ClassType)
	classType = strings.ToUpper(classType)
	subject := strings.TrimSpace(event.Subject)
	writeLong(w, "SUMMARY:", classType, " ", subject)
}

func (event *Event) endTime() time.Time {
	return event.StartTime.Add(time.Minute * 90)
}

func writeEvent(w io.Writer, event Event) {
	timeNow := time.Now().UTC().Format(time.RFC3339)

	write(w, "BEGIN:VEVENT")
	write(w, "TRANSP:OPAQUE")

	write(w, "DTSTART;TZID=Europe/Moscow:", event.StartTime.Format(timeFormat))
	write(w, "DTEND;TZID=Europe/Moscow:", event.endTime().Format(timeFormat))

	write(w, "UID:", uuid.New().String())

	writeRepeatRule(w, &event)
	writeAppleLocation(w, &event.Classroom)
	writeLocation(w, &event.Classroom)

	write(w, "X-APPLE-TRAVEL-ADVISORY-BEHAVIOR:DISABLED")
	write(w, "SEQUENCE:0")

	writeSummary(w, &event)

	writeLong(w, "DESCRIPTION:", event.Lecturer)

	write(w, "DTSTAMP:", timeNow)
	write(w, "CREATED:", timeNow)
	write(w, "LAST-MODIFIED:", timeNow)

	write(w, "END:VEVENT")
}
