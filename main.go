package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		log.Fatal("no website provided")
	}

	if len(args) > 1 {
		log.Fatal("too many arguments provided")
	}
	baseURL := args[0]
	pages := make(map[string]int)

	crawlPage(baseURL, baseURL, pages)
	fmt.Printf("Pages scraped: %v\n", len(pages))
}

func getHTML(rawURL string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "BootCrawler/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("error status code: %v\n", resp.StatusCode)
	}
	if !strings.Contains(resp.Header.Get("content-type"), "text/html") {
		return "", fmt.Errorf("response is not text/html\n")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return
	}

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return
	}

	if currentURL.Hostname() != baseURL.Hostname() {
		return // Skips other websites
	}

	normCurrURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	// increment and return if page has been crawled, start at 1 if it has not
	if _, v := pages[normCurrURL]; v {
		pages[normCurrURL]++
		return
	}
	pages[normCurrURL] = 1

	fmt.Printf("Crawling page: %v\n", rawCurrentURL)
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}

	nextURLs, err := getURLsFromHTML(html, baseURL)
	if err != nil {
		return
	}

	for _, nextURL := range nextURLs {
		crawlPage(rawBaseURL, nextURL, pages)
	}

}
