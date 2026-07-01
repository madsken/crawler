package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	URL            string   `json:"url"`
	Heading        string   `json:"heading"`
	FirstParagraph string   `json:"first_paragraph"`
	OutgoingLinks  []string `json:"outgoing_links"`
	ImageURLs      []string `json:"image_urls"`
}

func NewPageData(u, h, p string) PageData {
	return PageData{
		URL:            u,
		Heading:        h,
		FirstParagraph: p,
		ImageURLs:      nil,
		OutgoingLinks:  nil,
	}
}

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

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var result []string

	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok {
			return
		}
		href = strings.TrimSpace(href)
		if href == "" {
			return
		}

		// Parse href as url
		u, err := url.Parse(href)
		if err != nil {
			fmt.Printf("Could not parse href %q: %v\n", href, err)
			return
		}

		// ResolveReference: creates a new url, based on baseurl and reference url.
		// if reference url is abs url, just use abs url.
		// if ref url is rel url, construct new url based on baseurl and rel url
		resolved := baseURL.ResolveReference(u)
		result = append(result, resolved.String())
	})

	return result, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var result []string

	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if !ok {
			return
		}
		src = strings.TrimSpace(src)
		if src == "" {
			return
		}

		u, err := url.Parse(src)
		if err != nil {
			fmt.Printf("Could not parse href %q: %v\n", src, err)
		}

		resolved := baseURL.ResolveReference(u)
		result = append(result, resolved.String())
	})

	return result, nil
}

func extractPageData(html, pageURL string) PageData {
	pageData := NewPageData(pageURL, getHeadingFromHTML(html), getFirstParagraphFromHTML(html))

	baseURL, err := url.Parse(pageURL)
	if err != nil {
		return pageData
	}

	links, err := getURLsFromHTML(html, baseURL)
	if err != nil {
		return pageData
	}

	imgs, err := getImagesFromHTML(html, baseURL)
	if err != nil {
		return pageData
	}
	pageData.OutgoingLinks = links
	pageData.ImageURLs = imgs

	return pageData
}
