package rhetoric

import (
	"strings"

	"github.com/inkcheck/shared"
)

// TemporalResult holds the epistemic-temporal orientation analysis of a text.
// High scores indicate prospective/aspirational orientation;
// low scores indicate retrospective/evidential orientation.
type TemporalResult struct {
	FutureModalCount   int
	PastTenseCount     int
	EvidentialCount    int
	AspirationCount    int
	FutureModalDensity float64 // per 100 words
	PastTenseDensity   float64 // per 100 words
	EvidentialDensity  float64 // per 100 words
	AspirationDensity  float64 // per 100 words
}

// futureModalTokens are single-word future-oriented modals matched as whole tokens.
var futureModalTokens = map[string]bool{
	"will": true, "shall": true, "gonna": true,
}

// futureModalPhrases are genuinely multi-word future constructions.
// Single-token modals (will, shall) are already counted in futureModalTokens,
// so phrases like "we will" or "will be" are excluded to avoid double-counting.
var futureModalPhrases = []string{
	"going to", "are going to", "is going to", "am going to",
}

// irregularPastForms are common irregular past-tense verb forms.
var irregularPastForms = map[string]bool{
	"was": true, "were": true, "had": true, "went": true, "came": true,
	"saw": true, "made": true, "took": true, "got": true, "gave": true,
	"found": true, "knew": true, "thought": true, "told": true, "became": true,
	"left": true, "felt": true, "kept": true, "put": true, "said": true,
	"sent": true, "set": true, "cut": true, "hit": true, "led": true,
	"met": true, "ran": true, "sat": true, "stood": true, "heard": true,
	"held": true, "lost": true, "meant": true, "paid": true, "sold": true,
	"spent": true, "won": true, "wrote": true, "brought": true, "bought": true,
	"caught": true, "fought": true, "taught": true, "sought": true, "built": true,
	"dealt": true, "drove": true, "flew": true, "grew": true, "hung": true,
	"laid": true, "lent": true, "rode": true, "rose": true, "shone": true,
	"slept": true, "slid": true, "stole": true, "swam": true, "swore": true,
	"threw": true, "tore": true, "woke": true, "wore": true, "read": true,
	"began": true, "broke": true, "chose": true, "fell": true, "forgot": true,
	"froze": true, "hid": true, "spoke": true, "understood": true, "withdrew": true,
}

// notPastTense is an exclusion list for common false positives of the -ed heuristic.
var notPastTense = map[string]bool{
	"need": true, "seed": true, "feed": true, "breed": true, "speed": true,
	"bleed": true, "creed": true, "greed": true, "steed": true, "reed": true,
	"weed": true, "deed": true, "heed": true,
	"naked": true, "sacred": true, "wicked": true, "beloved": true,
	"rugged": true, "ragged": true, "crooked": true, "learned": true,
	"alleged": true, "aged": true, "blessed": true, "cursed": true,
	"dogged": true, "fixed": true, "mixed": true, "used": true,
	"supposed": true, "biased": true, "based": true,
	"hundred": true, "kindred": true,
}

var evidentialMarkers = []string{
	"according to", "studies show", "research indicates", "data suggests",
	"evidence shows", "it has been found", "findings suggest", "the data shows",
	"research shows", "studies indicate", "evidence suggests", "as shown by",
	"as demonstrated by", "statistics show", "the results show", "analysis shows",
	"it was found", "it has been shown", "the study found", "results indicate",
}

var aspirationMarkers = []string{
	"we aim to", "our goal is", "we plan to", "we're working toward",
	"our vision", "we hope to", "we're committed to", "we believe",
	"we aspire", "our mission is", "we strive to", "our objective is",
	"we intend to", "looking forward to", "we are committed", "we endeavour",
	"we endeavor", "we seek to", "our purpose is", "we want to",
	"our ambition", "we're dedicated to", "our aim is",
}

// TemporalOrientation analyses the epistemic-temporal mode of a text.
// It counts future modals, past-tense indicators, evidential phrases,
// and aspiration/intention phrases, returning densities per 100 words.
func TemporalOrientation(text string) TemporalResult {
	prose := shared.ExtractProseText(text)
	lower := strings.ToLower(prose)
	words := strings.Fields(lower)
	totalWords := len(words)

	if totalWords == 0 {
		return TemporalResult{}
	}

	// Count future modals: single tokens
	futureCount := 0
	for _, w := range words {
		token := shared.StripPunctuation(w)
		if futureModalTokens[token] {
			futureCount++
		}
	}
	// Count future modal phrases
	for _, phrase := range futureModalPhrases {
		futureCount += shared.CountOccurrences(lower, phrase)
	}

	// Count past tense: irregular forms + regular -ed suffix heuristic
	pastCount := 0
	for _, w := range words {
		token := shared.StripPunctuation(w)
		if token == "" {
			continue
		}
		if irregularPastForms[token] {
			pastCount++
		} else if len(token) > 4 && strings.HasSuffix(token, "ed") && !notPastTense[token] {
			pastCount++
		}
	}

	// Count evidential markers
	evidentialCount := 0
	for _, phrase := range evidentialMarkers {
		evidentialCount += shared.CountOccurrences(lower, phrase)
	}

	// Count aspiration markers
	aspirationCount := 0
	for _, phrase := range aspirationMarkers {
		aspirationCount += shared.CountOccurrences(lower, phrase)
	}

	perHundred := func(n int) float64 {
		return float64(n) / float64(totalWords) * 100
	}

	return TemporalResult{
		FutureModalCount:   futureCount,
		PastTenseCount:     pastCount,
		EvidentialCount:    evidentialCount,
		AspirationCount:    aspirationCount,
		FutureModalDensity: perHundred(futureCount),
		PastTenseDensity:   perHundred(pastCount),
		EvidentialDensity:  perHundred(evidentialCount),
		AspirationDensity:  perHundred(aspirationCount),
	}
}
