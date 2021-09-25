package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(input string) []string {
	counter := make(map[string]int)
	words := make([]string, 0)

	for _, w := range strings.Fields(input) {
		if _, ok := counter[w]; !ok {
			words = append(words, w)
		}
		counter[w]++
	}

	sort.Slice(words, func(i, j int) bool {
		if counter[words[i]] == counter[words[j]] {
			return words[i] < words[j]
		}
		return counter[words[i]] > counter[words[j]]
	})

	top := 10
	if top > len(words) {
		top = len(words)
	}
	return words[:top]
}
