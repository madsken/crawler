package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	result := doc.Find("h1, h2").First().Text()

	return result
}

func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	result := doc.Find("main").Find("p").First().Text()
	if result == "" {
		result = doc.Find("p").First().Text()
	}

	return result
}
