package semantic

import (
	"github.com/inkcheck/config"
	"github.com/inkcheck/shared"
)

// TopicCoherenceResult holds topic coherence analysis.
type TopicCoherenceResult struct {
	// PairSimilarities is the cosine similarity between consecutive paragraphs.
	PairSimilarities []float64
	// MeanSimilarity is the average consecutive-paragraph similarity.
	MeanSimilarity float64
	// CV is the coefficient of variation of similarities.
	CV float64
}

// TopicCoherence measures how consistently topics flow between consecutive
// paragraphs using word embeddings. Config is accepted for API consistency
// with other metric functions.
func TopicCoherence(_ config.Config, m *ModelManager, text string) TopicCoherenceResult {
	if m == nil {
		return TopicCoherenceResult{}
	}
	paragraphs := shared.SplitParagraphs(text)
	if len(paragraphs) < 2 {
		return TopicCoherenceResult{}
	}

	embeddings := make([][]float32, len(paragraphs))
	for i, p := range paragraphs {
		embeddings[i] = m.ParagraphEmbedding(p)
	}

	similarities := make([]float64, 0, len(paragraphs)-1)
	for i := 0; i < len(embeddings)-1; i++ {
		if embeddings[i] == nil || embeddings[i+1] == nil {
			continue
		}
		sim := shared.CosineSimilarity(embeddings[i], embeddings[i+1])
		similarities = append(similarities, sim)
	}

	if len(similarities) == 0 {
		return TopicCoherenceResult{}
	}

	return TopicCoherenceResult{
		PairSimilarities: similarities,
		MeanSimilarity:   shared.Mean(similarities),
		CV:               shared.CoefficientOfVariation(similarities),
	}
}
