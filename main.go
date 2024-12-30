package main

// coger top level postings
// filtrar por la primera lista
// 3 lineas de la descripcion

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type JobEntry struct {
	CompanyTitle string
	Location     string
	URL          string
	JobTitle     string
}

func main() {
	file, err := os.Open(os.ExpandEnv("$HOME/Desktop/hn.html"))
	if err != nil {
		log.Fatal("fetch: failed to open file", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("fetch: failed to read file", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(data)))
	if err != nil {
		log.Fatal(err)
	}

	entries := make([]JobEntry, 0)

	doc.Find("tr.athing.comtr").Each(func(i int, tr *goquery.Selection) {
		// if id, exists := tr.Attr("id"); exists {
		// 	fmt.Printf("ID %d: %s\n", i, id)
		// }

		tr.Find("td.default div.commtext").Each(func(i int, div *goquery.Selection) {
			commentBody := div.Text()

			if !strings.Contains(commentBody, "|") {
				return
			}

			parsed, err := parseEntry(commentBody)
			if err != nil {
				//				log.Println("parse: bad body", err)
				return
			}

			entries = append(entries, parsed)
		})
	})

	for e := range entries {
		fmt.Println(e)
	}
}

func parseEntry(body string) (JobEntry, error) {
	entry := JobEntry{}
	firstLine := strings.Split(body, "\n")[0]
	entry.CompanyTitle = strings.Split(firstLine, "|")[0]

	if strings.Contains("Amsterdam", firstLine) ||
		strings.Contains("EU", firstLine) ||
		strings.Contains("Netherlands ", firstLine) ||
		(strings.Contains("REMOTE", firstLine) && !strings.Contains("REMOTE (US", firstLine)) {
		entry.Location = firstLine
	} else {
		return JobEntry{}, fmt.Errorf("parse: location not interesting")
	}

	entry.URL = "donthave.com"
	entry.JobTitle = "engineer"

	return entry, nil
}
