package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type keyVal struct {
	word  string
	count int
}

func Top10(inputText string) []string {
	words := strings.Fields(inputText)

	wordCnt := make(map[string]int)
	for _, w := range words {
		wordCnt[w]++
	}

	counts := make([]keyVal, 0, len(wordCnt))
	for w, c := range wordCnt {
		counts = append(counts, keyVal{w, c})
	}

	sort.Slice(counts, func(i, j int) bool {
		if counts[i].count == counts[j].count {
			return counts[i].word < counts[j].word
		}
		return counts[i].count > counts[j].count
	})

	outputLen := min(10, len(counts))

	output := make([]string, outputLen)
	for i := 0; i < outputLen; i++ {
		output[i] = counts[i].word
	}

	return output
}
