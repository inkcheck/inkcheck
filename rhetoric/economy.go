package rhetoric

import (
	"strings"

	"github.com/inkcheck/shared"
)

// EconomyResult holds metrics about the efficiency and conciseness of writing.
type EconomyResult struct {
	AvgSentenceLength  float64 // mean words per sentence
	WordyPhraseCount   int
	WordyPhraseDensity float64 // wordy phrases per 100 words
	WordsPerClause     float64 // mean words per clause (approximated)
	ClauseCount        int
	SubordinationIndex float64 // subordinating conjunctions per sentence
}

// wordyPhrases is a curated list of verbose or redundant phrase patterns.
var wordyPhrases = []string{
	"in order to", "due to the fact that", "it is important to note",
	"at this point in time", "for the purpose of", "in the event that",
	"it should be noted that", "in light of the fact that", "the fact that",
	"it goes without saying", "as a matter of fact", "in terms of",
	"with regard to", "with respect to", "on the part of", "in the case of",
	"by means of", "for the reason that", "in the process of",
	"at the present time", "each and every", "first and foremost",
	"needless to say", "whether or not", "in my personal opinion",
	"completely eliminate", "future plans", "past history", "end result",
	"free gift", "unexpected surprise", "added bonus", "close proximity",
	"advance planning", "at a later date", "during the course of",
	"in the near future", "on a daily basis", "owing to the fact that",
	"prior to the", "subsequent to the", "in connection with",
	"with reference to", "in relation to", "in excess of",
	"a large number of", "a majority of", "are of the opinion",
	"at this moment in time", "has the ability to", "in the majority of cases",
	"in the not too distant future", "it is of great importance",
	"it may be argued that", "on the occasion of",
	"take into consideration", "the question as to whether",
}

// subordinatingConjunctions used to approximate clause boundaries.
var subordinatingConjunctions = []string{
	"although", "because", "since", "while", "whereas", "even though",
	"unless", "until", "whenever", "wherever", "whoever", "whatever",
	"provided that", "so that", "in order that", "rather than",
	"as soon as", "as long as", "now that", "in case",
}

// clauseBoundaryMarkers are single-character clause separators.
const clauseBoundaryChars = ";:"

// EconomyAnalysis measures the efficiency and conciseness of the text.
// It uses shared.SplitSentences for sentence splitting and approximates
// clause count via punctuation and subordinating conjunction counting.
func EconomyAnalysis(text string) EconomyResult {
	prose := shared.ExtractProseText(text)
	if prose == "" {
		return EconomyResult{}
	}

	lower := strings.ToLower(prose)
	words := strings.Fields(prose)
	totalWords := len(words)
	if totalWords == 0 {
		return EconomyResult{}
	}

	sentences := shared.SplitSentences(prose)
	sentenceCount := len(sentences)

	// Average sentence length
	avgSentLen := 0.0
	if sentenceCount > 0 {
		avgSentLen = float64(totalWords) / float64(sentenceCount)
	}

	// Wordy phrase count
	wordyCount := 0
	for _, phrase := range wordyPhrases {
		wordyCount += shared.CountOccurrences(lower, phrase)
	}
	wordyDensity := 0.0
	if totalWords > 0 {
		wordyDensity = float64(wordyCount) / float64(totalWords) * 100
	}

	// Clause count approximation:
	// Start with one clause per sentence, add boundaries from punctuation and conjunctions.
	clauseCount := sentenceCount
	// Count clause-boundary punctuation (semicolons, colons)
	for _, ch := range prose {
		if strings.ContainsRune(clauseBoundaryChars, ch) {
			clauseCount++
		}
	}
	// Count commas that introduce a new clause (followed by a conjunction or relative pronoun)
	clauseCommaPatterns := []string{
		", and ", ", but ", ", or ", ", nor ", ", so ", ", yet ", ", for ",
		", which ", ", who ", ", whom ", ", whose ", ", that ",
	}
	for _, pattern := range clauseCommaPatterns {
		clauseCount += strings.Count(lower, pattern)
	}
	// Add subordinating conjunction occurrences
	subordCount := 0
	for _, conj := range subordinatingConjunctions {
		n := shared.CountOccurrences(lower, conj)
		subordCount += n
		clauseCount += n
	}

	wordsPerClause := 0.0
	if clauseCount > 0 {
		wordsPerClause = float64(totalWords) / float64(clauseCount)
	}

	subordIndex := 0.0
	if sentenceCount > 0 {
		subordIndex = float64(subordCount) / float64(sentenceCount)
	}

	return EconomyResult{
		AvgSentenceLength:  avgSentLen,
		WordyPhraseCount:   wordyCount,
		WordyPhraseDensity: wordyDensity,
		WordsPerClause:     wordsPerClause,
		ClauseCount:        clauseCount,
		SubordinationIndex: subordIndex,
	}
}
