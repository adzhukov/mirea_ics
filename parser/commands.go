package parser

import (
	"io"
	"log"
	"net/http"
)

func Parse(uri string, g string) {
	wb, err := openFile(uri)
	if err != nil {
		log.Fatalln(err)
	}

	parse(wb, g)
}

func Groups(file string) []string {
	wb, err := openFile(file)
	if err != nil {
		log.Fatalln(err)
	}

	return groups(wb)
}

func ParseAllGroups(file string, mergeFlag bool) {
	wb, err := openFile(file)
	if err != nil {
		log.Fatalln(err)
	}

	for _, group := range groups(wb) {
		parse(wb, group)
		if mergeFlag {
			merge(group)
		}
	}
}

func Merge(group string) {
	merge(group)
}

func GetLinks(groups []string) []string {
	resp, err := http.Get("https://mirea.ru/schedule")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	m := make(map[string]struct{})
	for _, group := range groups {
		links := filterGroups(body, normalizeGroup(group))
		for _, link := range links {
			m[link] = struct{}{}
		}
	}

	links := make([]string, 0, len(m))
	for link := range m {
		links = append(links, link)
	}

	return links
}
