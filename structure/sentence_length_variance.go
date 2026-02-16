package structure

import "github.com/inkcheck/shared"

func SentenceLengthVariance(text string) float64 {
	sentences := shared.SplitSentences(shared.ExtractProseText(text))
	if len(sentences) < 2 {
		return 0
	}
	lengths := make([]float64, len(sentences))
	for i, s := range sentences {
		lengths[i] = float64(shared.CountWords(s))
	}
	return shared.CoefficientOfVariation(lengths)
}
