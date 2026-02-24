// Package wordlist embeds the Google 10,000 English word frequency list.
//
// Source: https://github.com/first20hours/google-10000-english
// See CREDITS.md for full attribution details.
package wordlist

import (
	_ "embed"
	"strings"
	"sync"
)

//go:embed google-10000-english.txt
var google10k string

var (
	once     sync.Once
	rankMap  map[string]int
)

// FrequencyRank returns a map of words to their frequency rank (1 = most common).
func FrequencyRank() map[string]int {
	once.Do(func() {
		lines := strings.Split(strings.TrimSpace(google10k), "\n")
		rankMap = make(map[string]int, len(lines))
		for i, line := range lines {
			word := strings.TrimSpace(line)
			if word != "" {
				rankMap[word] = i + 1
			}
		}
	})
	return rankMap
}
