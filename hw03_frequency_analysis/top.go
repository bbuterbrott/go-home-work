package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"regexp"
	"sort"
	"strings"
)

// Top10 returns top 10 words from the input string.
func Top10(s string) []string {
	words := strings.Split(s, " ")
	re := regexp.MustCompile(`[\p{L}-]*`)
	reEmpty := regexp.MustCompile(`\s+`)
	wordCount := make(map[string]int)
	for _, w := range words {
		filteredWords := re.FindAllString(w, -1)

		for _, fw := range filteredWords {
			filteredWord := strings.ToLower(fw)
			if filteredWord == "" || filteredWord == "-" || reEmpty.MatchString(filteredWord) {
				continue
			}
			wordCount[filteredWord]++
		}
	}

	type kv struct {
		Key   string
		Value int
	}

	kvs := make([]kv, len(wordCount))
	i := 0
	for k, v := range wordCount {
		kvs[i] = kv{k, v}
		i++
	}

	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].Value > kvs[j].Value
	})

	var top10 []string
	var top10kvs []kv
	if len(kvs) > 10 {
		top10 = make([]string, 10)
		top10kvs = make([]kv, 10)
		copy(top10kvs, kvs[:10])
	} else {
		top10 = make([]string, len(kvs))
		top10kvs = kvs
	}

	for i, kv := range top10kvs {
		top10[i] = kv.Key
	}

	return top10
}
