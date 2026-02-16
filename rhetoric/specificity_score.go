package rhetoric

import (
	"strings"

	"github.com/inkcheck/shared"
)

type SpecificityResult struct {
	Mean   float64
	Range  float64
	CV     float64
	Scores []int
}

func SpecificityScore(text string) SpecificityResult {
	prose := shared.ExtractProseText(text)
	sentences := shared.SplitSentences(prose)
	if len(sentences) < 2 {
		return SpecificityResult{}
	}

	scores := make([]int, len(sentences))
	for i, s := range sentences {
		scores[i] = estimateSpecificity(s)
	}

	return buildSpecificityResult(scores)
}

func buildSpecificityResult(scores []int) SpecificityResult {
	values := make([]float64, len(scores))
	minVal, maxVal := 5, 1
	for i, s := range scores {
		values[i] = float64(s)
		if s < minVal {
			minVal = s
		}
		if s > maxVal {
			maxVal = s
		}
	}

	return SpecificityResult{
		Mean:   shared.Mean(values),
		Range:  float64(maxVal - minVal),
		CV:     shared.CoefficientOfVariation(values),
		Scores: scores,
	}
}

func estimateSpecificity(sentence string) int {
	score := 2
	lower := strings.ToLower(sentence)

	hasNumber := false
	for _, r := range sentence {
		if r >= '0' && r <= '9' {
			hasNumber = true
			break
		}
	}
	if hasNumber {
		score++
	}

	if strings.Contains(sentence, "%") || strings.Contains(sentence, "$") {
		score++
	}

	words := shared.ListWords(sentence)
	properNouns := 0
	for i, w := range words {
		if i == 0 {
			continue
		}
		if len(w) > 1 && w[0] >= 'A' && w[0] <= 'Z' {
			properNouns++
		}
	}
	if properNouns >= 2 {
		score++
	}

	if strings.ContainsAny(sentence, "\"\u201c\u201d") {
		score++
	}

	abstractSignals := []string{"important", "significant", "crucial", "essential",
		"fundamental", "effective", "various", "numerous", "overall"}
	for _, signal := range abstractSignals {
		if strings.Contains(lower, signal) {
			score--
			break
		}
	}

	concreteSignals := []string{"for example", "specifically", "in particular",
		"such as", "namely", "according to"}
	for _, signal := range concreteSignals {
		if strings.Contains(lower, signal) {
			score++
			break
		}
	}

	if score < 1 {
		score = 1
	}
	if score > 5 {
		score = 5
	}
	return score
}
