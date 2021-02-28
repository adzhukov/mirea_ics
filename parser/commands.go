package parser

import (
	"io/ioutil"
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

func GetLinks(group string) []string {
	resp, err := http.Get("https://mirea.ru/schedule")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return filterGroups(body, normalizeGroup(group))
}

func Merge(group string) {
	merge(group)
}
