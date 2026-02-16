package rhetoric

import (
	"strings"

	"github.com/inkcheck/shared"
)

// CounterargumentResult holds counterargument engagement analysis.
type CounterargumentResult struct {
	Instances     int
	DensityPer100 float64
	Phrases       []string
}

var counterargumentPhrases = []string{
	"on the other hand", "however", "nevertheless", "nonetheless",
	"critics argue", "some might say", "one could argue",
	"it could be argued", "opponents suggest", "skeptics point out",
	"while it is true", "while this may", "although",
	"despite this", "in contrast", "conversely",
	"admittedly", "granted", "to be fair",
	"some may disagree", "a common objection",
}

// CounterargumentEngagement detects counterargument language in the text.
// Higher density suggests more engagement with opposing viewpoints.
// TODO: LLM judge for nuanced counterargument detection
func CounterargumentEngagement(text string) CounterargumentResult {
	prose := shared.ExtractProseText(text)
	wordCount := shared.CountWords(prose)
	if wordCount == 0 {
		return CounterargumentResult{}
	}

	lower := strings.ToLower(prose)
	var found []string
	total := 0

	for _, phrase := range counterargumentPhrases {
		n := countOccurrences(lower, phrase)
		if n > 0 {
			found = append(found, phrase)
			total += n
		}
	}

	density := float64(total) / float64(wordCount) * 100

	return CounterargumentResult{
		Instances:     total,
		DensityPer100: density,
		Phrases:       found,
	}
}
