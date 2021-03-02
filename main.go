package main

import (
	"flag"
	"fmt"

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
		links := parser.GetLinks(args)
		for _, link := range links {
			fmt.Println(link)
		}

		return
	}

	if list {
		groups := parser.Groups(file)
		for _, link := range groups {
			fmt.Println(link)
		}
	}

	if all {
		parser.ParseAllGroups(file, merge)
		return
	}

	if merge {
		for _, group := range args {
			parser.Merge(group)
		}

		return
	}

	for _, group := range args {
		parser.Parse(file, group)
	}
}
