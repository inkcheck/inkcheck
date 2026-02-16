package rhetoric

import (
	"strings"

	"github.com/inkcheck/shared"
)

// RhetoricalDiversityResult holds sentence type distribution results.
type RhetoricalDiversityResult struct {
	Questions    int
	Exclamations int
	Imperatives  int
	Conditionals int
	Declaratives int
	Total        int
	Entropy      float64
}

// RhetoricalDiversity classifies sentences by type and computes Shannon entropy
// of the distribution. Higher entropy means more diverse sentence types.
// TODO: LLM judge for nuanced sentence type classification
func RhetoricalDiversity(text string) RhetoricalDiversityResult {
	prose := shared.ExtractProseText(text)
	sentences := shared.SplitSentences(prose)
	if len(sentences) == 0 {
		return RhetoricalDiversityResult{}
	}

	var questions, exclamations, imperatives, conditionals, declaratives int

	for _, s := range sentences {
		trimmed := strings.TrimSpace(s)
		lower := strings.ToLower(trimmed)

		switch {
		case strings.HasSuffix(trimmed, "?"):
			questions++
		case strings.HasSuffix(trimmed, "!"):
			exclamations++
		case isImperative(lower):
			imperatives++
		case isConditional(lower):
			conditionals++
		default:
			declaratives++
		}
	}

	total := len(sentences)
	counts := []int{questions, exclamations, imperatives, conditionals, declaratives}
	dist := make([]float64, len(counts))
	for i, c := range counts {
		dist[i] = float64(c) / float64(total)
	}

	return RhetoricalDiversityResult{
		Questions:    questions,
		Exclamations: exclamations,
		Imperatives:  imperatives,
		Conditionals: conditionals,
		Declaratives: declaratives,
		Total:        total,
		Entropy:      shared.Entropy(dist),
	}
}

var imperativeVerbs = []string{
	"consider", "note", "remember", "imagine", "think", "look",
	"see", "let", "take", "make", "do", "try", "ensure", "avoid",
}

func isImperative(lower string) bool {
	for _, v := range imperativeVerbs {
		if strings.HasPrefix(lower, v+" ") {
			return true
		}
	}
	return false
}

func isConditional(lower string) bool {
	return strings.HasPrefix(lower, "if ") ||
		strings.HasPrefix(lower, "when ") ||
		strings.HasPrefix(lower, "unless ") ||
		strings.HasPrefix(lower, "provided ") ||
		strings.Contains(lower, " would ") ||
		strings.Contains(lower, " could have ")
}
