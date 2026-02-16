package structure

import (
	"strings"

	"github.com/inkcheck/shared"
)

// ParagraphVariance computes the coefficient of variation of paragraph lengths.
// Returns 0 if the text is empty or contains fewer than 2 paragraphs.
func ParagraphVariance(text string) float64 {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}
	paragraphs := shared.SplitParagraphs(text)
	if len(paragraphs) < 2 {
		return 0
	}
	lengths := make([]float64, len(paragraphs))
	for i, p := range paragraphs {
		lengths[i] = float64(shared.CountWords(p))
	}
	return shared.CoefficientOfVariation(lengths)
}

func ParagraphLengths(text string) []int {
	paragraphs := shared.SplitParagraphs(text)
	lengths := make([]int, len(paragraphs))
	for i, p := range paragraphs {
		lengths[i] = shared.CountWords(p)
	}
	return lengths
}
