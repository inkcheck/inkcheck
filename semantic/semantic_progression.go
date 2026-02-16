package semantic

import (
	"github.com/inkcheck/config"
	"github.com/inkcheck/shared"
)

// SemanticProgressionResult holds semantic drift analysis.
type SemanticProgressionResult struct {
	// DriftRates is 1-similarity between consecutive paragraphs.
	DriftRates []float64
	// MeanDrift is the average drift rate.
	MeanDrift float64
	// CV is the coefficient of variation of drift rates.
	CV float64
}

// SemanticProgression measures how much the topic shifts between consecutive
// paragraphs. Drift rate = 1 - similarity. Config is accepted for API
// consistency with other metric functions.
func SemanticProgression(_ config.Config, m *ModelManager, text string) SemanticProgressionResult {
	if m == nil {
		return SemanticProgressionResult{}
	}
	paragraphs := shared.SplitParagraphs(text)
	if len(paragraphs) < 2 {
		return SemanticProgressionResult{}
	}

	embeddings := make([][]float32, len(paragraphs))
	for i, p := range paragraphs {
		embeddings[i] = m.ParagraphEmbedding(p)
	}

	drifts := make([]float64, 0, len(paragraphs)-1)
	for i := 0; i < len(embeddings)-1; i++ {
		if embeddings[i] == nil || embeddings[i+1] == nil {
			continue
		}
		sim := shared.CosineSimilarity(embeddings[i], embeddings[i+1])
		drifts = append(drifts, 1.0-sim)
	}

	if len(drifts) == 0 {
		return SemanticProgressionResult{}
	}

	return SemanticProgressionResult{
		DriftRates: drifts,
		MeanDrift:  shared.Mean(drifts),
		CV:         shared.CoefficientOfVariation(drifts),
	}
}
