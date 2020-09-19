package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"regexp"
	"sort"
	"strings"
)

const (
	topWordCount = 10
)

var re = regexp.MustCompile(`[\p{L}-]*`)

// Top10 returns top 10 words from the input string.
func Top10(s string) []string {
	words := re.FindAllString(s, -1)
	wordCount := make(map[string]int)
	for _, fw := range words {
		filteredWord := strings.ToLower(fw)
		if filteredWord == "" || filteredWord == "-" {
			continue
		}
		wordCount[filteredWord]++
	}

	type kv struct {
		Key   string
		Value int
	}

	kvs := make([]kv, 0, len(wordCount))
	for k, v := range wordCount {
		kvs = append(kvs, kv{k, v})
	}

	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].Value > kvs[j].Value
	})

	cutTop := make([]string, 0, topWordCount)
	for i, kv := range kvs {
		cutTop = append(cutTop, kv.Key)
		if i >= topWordCount {
			break
		}
	}

	return cutTop
}
