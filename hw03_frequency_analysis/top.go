package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type topWord struct {
	word      string
	frequency int
}

func Top10(text string) []string {
	const wordsCount = 10

	if text == "" {
		return []string{}
	}

	words := strings.Fields(text)
	uniqueWords := getUniqueWords(words)
	topWords := make([]topWord, 0, len(uniqueWords))

	for word, freq := range uniqueWords {
		topWords = append(topWords, topWord{word: word, frequency: freq})
	}

	sort.Slice(topWords, func(i, j int) bool {
		if topWords[i].frequency == topWords[j].frequency {
			return topWords[i].word < topWords[j].word
		}
		return topWords[i].frequency > topWords[j].frequency
	})

	result := make([]string, 0, wordsCount)
	for i := 0; i < wordsCount; i++ {
		result = append(result, topWords[i].word)
	}

	return result
}

func getUniqueWords(words []string) map[string]int {
	uniqueWords := make(map[string]int)
	for _, word := range words {
		uniqueWords[word]++
	}
	return uniqueWords
}
