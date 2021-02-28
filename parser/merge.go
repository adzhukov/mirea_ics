package parser

import (
	"github.com/adzhukov/mirea_ics/calendar"
	"bytes"
	"os"
	"strings"
	"log"
	"bufio"
	"io"
)

const calBegin = "END:VTIMEZONE"
const calEnd = "END:VCALENDAR"

func merge(group string) {
	group = normalizeGroup(group)

	fileInfo, err := os.ReadDir("./")
	if err != nil {
		log.Fatalln(err)
	}

	empty := calendar.Calendar{Group: group}.String()
	empty = empty[:len(empty) - len(calEnd) - 2]
	buf := bytes.NewBufferString(empty)
	
	for _, file := range fileInfo {
		if name := file.Name(); strings.HasPrefix(name, group + " ") {
			writeEvents(buf, name)
		}
	}

	buf.WriteString(calEnd)
	file, err := os.Create(group + ".ics")
	if err != nil {
		panic(err)
	}

	buf.WriteTo(file)
	file.Close()
}

func writeEvents(w io.Writer, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if bytes.Compare(scanner.Bytes(), []byte(calBegin)) == 0 {
			break
		}
	}

	for scanner.Scan() {
		read := scanner.Bytes()
		if bytes.Compare(read, []byte(calEnd)) == 0 {
			break
		}

		w.Write(read)
		w.Write([]byte("\r\n"))
	}
}
