package main

import (
	"flag"
	"log"

	"github.com/adzhukov/mirea_ics/parser"
)

func main() {
	filePath := flag.String("file", "", ".xlsx file path")
	group := flag.String("group", "", "XXXX-00-00")
	flag.Parse()

	if *filePath == "" {
		log.Fatal("file flag shoul be set")
	}

	if *group == "" {
		log.Fatal("group flag should be set")
	}

	parser.ParseFile(*filePath, *group)
}
