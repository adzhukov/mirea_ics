package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/adzhukov/mirea_ics/parser"
)

var (
	file  string
	list  bool
	links bool
)

func init() {
	flag.StringVar(&file, "file", "", ".xlsx URI")
	flag.BoolVar(&list, "list", false, "print all groups in file")
	flag.BoolVar(&links, "links", false, "print all file links")
}

func main() {
	flag.Parse()

	if links {
		links := parser.GetLinks()
		for _, link := range links {
			fmt.Println(link)
		}

		return
	}

	if file == "" {
		log.Fatalln("file flag should be set")
	}

	if list {
		groups := parser.Groups(file)
		for _, link := range groups {
			fmt.Println(link)
		}
	}

	for _, group := range flag.Args() {
		parser.Parse(file, group)
	}
}
