package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Sort the map keys for deterministic output
// Build a []PageData slice in sorted order
// Marshal it with json.MarshalIndent using indent=2
// Write the result to disk with os.WriteFile
func writeJSONReport(pages map[string]PageData, filename string) error {
	if len(pages) == 0 {
		fmt.Println("no pages data, skipping json generation")
		return nil
	}

	keys := make([]string, 0, len(pages))
	for k := range pages {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sortedPages := make([]PageData, 0, len(pages))
	for _, k := range keys {
		sortedPages = append(sortedPages, pages[k])
	}

	data, err := json.MarshalIndent(sortedPages, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0o644)
	if err != nil {
		return err
	}

	fmt.Println("report.json generated")

	return nil
}
