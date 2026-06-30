package main

import (
	"net/url"
	"strings"
)

func normalizeURL(rawUrl string) (string, error) {
	parsed, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	result := parsed.Host + parsed.Path
	result, _ = strings.CutSuffix(result, "/")
	result = strings.ToLower(result)
	return result, nil
}
