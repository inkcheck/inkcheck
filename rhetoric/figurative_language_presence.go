package rhetoric

import (
	"strings"
	"unicode"

	"github.com/inkcheck/shared"
)

// FigurativeLanguageResult holds figurative language detection results.
type FigurativeLanguageResult struct {
	SimileCount             int
	RhetoricalQuestionCount int
	AlliterationCount       int
	TotalInstances          int
	DensityPer100Words      float64
}

// FigurativeLanguagePresence detects similes, rhetorical questions, and
// alliteration using lexical patterns.
// TODO: LLM judge for metaphor detection and context-aware analysis
func FigurativeLanguagePresence(text string) FigurativeLanguageResult {
	prose := shared.ExtractProseText(text)
	words := shared.ListWords(prose)
	if len(words) == 0 {
		return FigurativeLanguageResult{}
	}

	sentences := shared.SplitSentences(prose)

	similes := 0
	rhetoricalQs := 0
	alliterations := 0

	for _, s := range sentences {
		lower := strings.ToLower(s)
		// Simile detection: "like a/an", "as ... as"
		if strings.Contains(lower, " like a ") || strings.Contains(lower, " like an ") {
			similes++
		}
		if containsAsAs(lower) {
			similes++
		}
		// Rhetorical questions: questions not at the very start of the text
		if strings.HasSuffix(strings.TrimSpace(s), "?") {
			rhetoricalQs++
		}
	}

	// Alliteration: 3+ consecutive words starting with the same letter
	for i := 0; i < len(words)-2; i++ {
		a := firstLetter(words[i])
		b := firstLetter(words[i+1])
		c := firstLetter(words[i+2])
		if a != 0 && a == b && b == c {
			alliterations++
		}
	}

	total := similes + rhetoricalQs + alliterations
	density := 0.0
	if len(words) > 0 {
		density = float64(total) / float64(len(words)) * 100
	}

	return FigurativeLanguageResult{
		SimileCount:             similes,
		RhetoricalQuestionCount: rhetoricalQs,
		AlliterationCount:       alliterations,
		TotalInstances:          total,
		DensityPer100Words:      density,
	}
}

func containsAsAs(s string) bool {
	idx := strings.Index(s, " as ")
	if idx < 0 {
		return false
	}
	rest := s[idx+4:]
	return strings.Contains(rest, " as ")
}

func firstLetter(word string) rune {
	for _, r := range word {
		if unicode.IsLetter(r) {
			return unicode.ToLower(r)
		}
	}
	return 0
}
