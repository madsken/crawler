package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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

	html, err := getHTML(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(html)

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
