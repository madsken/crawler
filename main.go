package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) != 3 {
		log.Fatal("usage: crawler <base-url> <max-concurrency> <max-pages>")
	}

	baseURL := args[0]
	maxConcurrency, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal("Max Concurrency mus be a number")
	}
	maxPages, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatal("Max Pages mus be a number")
	}

	cfg, err := NewConfig(baseURL, maxConcurrency, maxPages)
	if err != nil {
		log.Fatal(err)
	}

	cfg.wg.Add(1)
	go cfg.crawlPage(baseURL)
	cfg.wg.Wait()

	fmt.Printf("Pages scraped: %v\n", len(cfg.pages))
	for page := range cfg.pages {
		fmt.Println(page)
	}
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
		return "", fmt.Errorf("error status code: %v", resp.StatusCode)
	}
	if !strings.Contains(resp.Header.Get("content-type"), "text/html") {
		return "", fmt.Errorf("response is not text/html")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	// Check and skip other websites from base url
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return
	}
	baseURL, err := url.Parse(cfg.baseURL.String())
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

	// checks if page has been visited before (in which case we skip)
	if !cfg.pageVisit(normCurrURL) {
		return
	}

	fmt.Printf("Crawling page: %v\n", rawCurrentURL)
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Add page data to cfg
	pagedata := extractPageData(html, rawCurrentURL)
	cfg.addPageData(normCurrURL, pagedata)

	// recursively crawl found links
	for _, nextURL := range pagedata.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(nextURL)
	}
}
