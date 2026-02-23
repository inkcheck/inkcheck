package rhetoric

import (
	"strings"
	"unicode"

	"github.com/inkcheck/shared"
)

// ContractionResult holds the count and rate of contractions in a text.
type ContractionResult struct {
	Count int     // number of contraction tokens found
	Rate  float64 // Count / total_words (0.0–1.0)
}

// contractionSet is the canonical set of English contractions (all lowercase).
var contractionSet = map[string]bool{
	"can't": true, "won't": true, "don't": true,
	"doesn't": true, "didn't": true, "isn't": true, "aren't": true,
	"wasn't": true, "weren't": true, "hasn't": true, "haven't": true,
	"hadn't": true, "wouldn't": true, "couldn't": true, "shouldn't": true,
	"mustn't": true, "needn't": true, "daren't": true,
	"it's": true, "that's": true, "what's": true, "who's": true,
	"there's": true, "here's": true, "where's": true, "when's": true,
	"why's": true, "how's": true, "let's": true,
	"i'm": true, "you're": true, "he's": true, "she's": true,
	"we're": true, "they're": true,
	"i've": true, "you've": true, "we've": true, "they've": true,
	"i'll": true, "you'll": true, "he'll": true, "she'll": true,
	"we'll": true, "they'll": true,
	"i'd": true, "you'd": true, "he'd": true, "she'd": true,
	"we'd": true, "they'd": true, "it'd": true,
	"ain't": true, "gonna": true, "wanna": true, "gotta": true,
}

// stripKeepApostrophe removes leading and trailing punctuation from a word
// while preserving internal and trailing apostrophes used in contractions.
// NOTE: This intentionally differs from shared.StripPunctuation, which strips
// all non-letter characters including apostrophes needed for contraction matching.
func stripKeepApostrophe(word string) string {
	runes := []rune(word)
	start, end := 0, len(runes)
	for start < end && !unicode.IsLetter(runes[start]) && runes[start] != '\'' {
		start++
	}
	for end > start && !unicode.IsLetter(runes[end-1]) && runes[end-1] != '\'' {
		end--
	}
	if start >= end {
		return ""
	}
	return string(runes[start:end])
}

// ContractionRate counts contractions in text and returns the rate per total words.
// Uses shared.ExtractProseText to strip markdown before analysis.
func ContractionRate(text string) ContractionResult {
	prose := shared.ExtractProseText(text)
	tokens := strings.Fields(prose)
	totalWords := len(tokens)

	if totalWords == 0 {
		return ContractionResult{}
	}

	count := 0
	for _, tok := range tokens {
		cleaned := stripKeepApostrophe(strings.ToLower(tok))
		if contractionSet[cleaned] {
			count++
		}
	}

	return ContractionResult{
		Count: count,
		Rate:  float64(count) / float64(totalWords),
	}
}
