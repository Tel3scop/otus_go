package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordRank struct {
	Value string
	Count int
}

func Top10(input string) []string {
	words := strings.Fields(input)
	wordCounts := countWords(words)
	wordsRank := rankWords(wordCounts)
	sortWords(wordsRank)
	return getTopWords(wordsRank, 10)
}

func countWords(words []string) map[string]int {
	wordCounts := make(map[string]int)
	for _, word := range words {
		wordCounts[word]++
	}
	return wordCounts
}

func rankWords(wordCounts map[string]int) []wordRank {
	wordsRank := make([]wordRank, 0, len(wordCounts))
	for word, count := range wordCounts {
		wordsRank = append(wordsRank, wordRank{Value: word, Count: count})
	}
	return wordsRank
}

func sortWords(wordsRank []wordRank) {
	sort.Slice(wordsRank, func(i, j int) bool {
		if wordsRank[i].Count > wordsRank[j].Count {
			return true
		}
		if wordsRank[i].Count < wordsRank[j].Count {
			return false
		}
		return wordsRank[i].Value < wordsRank[j].Value
	})
}

func getTopWords(wordsRank []wordRank, top int) []string {
	result := make([]string, 0, top)
	for i := 0; i < len(wordsRank) && i < top; i++ {
		result = append(result, wordsRank[i].Value)
	}
	return result
}
