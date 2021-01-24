package calendar

import (
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"
)

const delimiter = "\r\n"

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
