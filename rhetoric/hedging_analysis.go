package rhetoric

import (
	"strings"

	"github.com/inkcheck/shared"
)

const (
	// minWordsForHedgingAnalysis is the minimum number of words required
	// to perform meaningful hedging analysis.
	minWordsForHedgingAnalysis = 10
)

type HedgingResult struct {
	Total      int
	Density    float64
	Distinct   int
	Variety    float64
	Categories HedgingCategories
	Hedges     []HedgeInstance
}

type HedgingCategories struct {
	Modal        int
	Approximator int
	Plausibility int
	Attribution  int
	Frequency    int
}

type HedgeInstance struct {
	Text     string `json:"text"`
	Category string `json:"category"`
}

func HedgingAnalysis(text string) HedgingResult {
	prose := shared.ExtractProseText(text)
	wordCount := shared.CountWords(prose)
	if wordCount < minWordsForHedgingAnalysis {
		return HedgingResult{}
	}

	lower := strings.ToLower(prose)
	var hedges []HedgeInstance

	for _, h := range hedgingPatterns {
		n := countOccurrences(lower, h.phrase)
		for range n {
			hedges = append(hedges, HedgeInstance{
				Text:     h.phrase,
				Category: h.category,
			})
		}
	}

	return buildHedgingResult(hedges, wordCount)
}

func buildHedgingResult(hedges []HedgeInstance, wordCount int) HedgingResult {
	total := len(hedges)
	density := 0.0
	if wordCount > 0 {
		density = float64(total) / float64(wordCount) * 100
	}

	distinct := make(map[string]bool)
	var cats HedgingCategories
	for _, h := range hedges {
		distinct[strings.ToLower(h.Text)] = true
		switch strings.ToLower(h.Category) {
		case "modal":
			cats.Modal++
		case "approximator":
			cats.Approximator++
		case "plausibility":
			cats.Plausibility++
		case "attribution":
			cats.Attribution++
		case "frequency":
			cats.Frequency++
		}
	}

	variety := 0.0
	if total > 0 {
		variety = float64(len(distinct)) / float64(total)
	}

	return HedgingResult{
		Total:      total,
		Density:    density,
		Distinct:   len(distinct),
		Variety:    variety,
		Categories: cats,
		Hedges:     hedges,
	}
}

type hedgePattern struct {
	phrase   string
	category string
}

var hedgingPatterns = []hedgePattern{
	{"may", "modal"}, {"might", "modal"}, {"could", "modal"},
	{"approximately", "approximator"}, {"roughly", "approximator"},
	{"somewhat", "approximator"}, {"nearly", "approximator"},
	{"perhaps", "plausibility"}, {"probably", "plausibility"},
	{"possibly", "plausibility"}, {"likely", "plausibility"},
	{"it is possible", "plausibility"}, {"it's possible", "plausibility"},
	{"it seems", "attribution"}, {"appears to", "attribution"},
	{"tends to", "attribution"}, {"arguably", "attribution"},
	{"is thought to", "attribution"},
	{"sometimes", "frequency"}, {"often", "frequency"},
	{"generally", "frequency"}, {"usually", "frequency"},
	{"typically", "frequency"}, {"in some cases", "frequency"},
	{"in many cases", "frequency"}, {"frequently", "frequency"},
}
