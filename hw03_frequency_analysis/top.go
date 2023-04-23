package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	const topWordsCount = 10

	if text == "" {
		return nil
	}

	words := strings.Fields(text)

	uniqueWords := make(map[string]int)
	for _, word := range words {
		uniqueWords[word]++
	}

	topWords := make([]string, 0, len(uniqueWords))
	for word := range uniqueWords {
		topWords = append(topWords, word)
	}

	sort.Slice(topWords, func(i, j int) bool {
		if uniqueWords[topWords[i]] == uniqueWords[topWords[j]] {
			return topWords[i] < topWords[j]
		}
		return uniqueWords[topWords[i]] > uniqueWords[topWords[j]]
	})

	if len(topWords) < topWordsCount {
		return topWords
	}
	return topWords[:topWordsCount]
}
