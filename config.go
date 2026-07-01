package main

import (
	"fmt"
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	wg                 *sync.WaitGroup
	concurrencyControl chan struct{}
	maxPages           int
}

func (cfg *config) pageVisit(normURL string) (isFirst bool) {
	// increment and return if page has been crawled, start at 1 if it has not
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	if _, v := cfg.pages[normURL]; v {
		return false
	}

	if len(cfg.pages) >= cfg.maxPages {
		return
	}
	cfg.pages[normURL] = PageData{URL: normURL}
	return true
}

func (cfg *config) addPageData(normURL string, pagedata PageData) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.pages[normURL] = pagedata
}

func (cfg *config) isPagesLengthExceeded() bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	return len(cfg.pages) >= cfg.maxPages
}

func NewConfig(rawBaseURL string, maxConcurrency int, pageLimit int) (*config, error) {
	baseURLParsed, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse config: %v", err)
	}

	return &config{
		pages:              make(map[string]PageData),
		baseURL:            baseURLParsed,
		mu:                 &sync.Mutex{},
		wg:                 &sync.WaitGroup{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		maxPages:           pageLimit,
	}, nil
}
