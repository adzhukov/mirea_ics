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
	merge bool
	all   bool
)

func init() {
	flag.StringVar(&file, "file", "", ".xlsx URI")
	flag.BoolVar(&list, "list", false, "print all groups in file")
	flag.BoolVar(&links, "links", false, "print all file links")
	flag.BoolVar(&merge, "merge", false, "merge calendars")
	flag.BoolVar(&all, "all", false, "make calendars for all groups in file")
}

func main() {
	flag.Parse()
	args := flag.Args()

	if links {
		group := "nil"
		if len(args) != 0 {
			group = args[0]
		}
		links := parser.GetLinks(group)
		for _, link := range links {
			fmt.Println(link)
		}

		return
	}

	if merge {
		for _, group := range args {
			parser.Merge(group)
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

	if all {
		parser.ParseAllGroups(file)
		return
	}

	for _, group := range args {
		parser.Parse(file, group)
	}
}
