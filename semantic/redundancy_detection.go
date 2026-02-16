package semantic

import (
	"github.com/inkcheck/config"
	"github.com/inkcheck/shared"
)

// RedundancyPair represents two non-adjacent paragraphs with high similarity.
type RedundancyPair struct {
	IndexA     int
	IndexB     int
	Similarity float64
}

// RedundancyResult holds redundancy detection results.
type RedundancyResult struct {
	Pairs     []RedundancyPair
	PairCount int
}

// RedundancyDetection finds all non-adjacent paragraph pairs with
// cosine similarity above the configured threshold.
func RedundancyDetection(cfg config.Config, m *ModelManager, text string) RedundancyResult {
	if m == nil {
		return RedundancyResult{}
	}
	paragraphs := shared.SplitParagraphs(text)
	if len(paragraphs) < 3 {
		return RedundancyResult{}
	}

	embeddings := make([][]float32, len(paragraphs))
	for i, p := range paragraphs {
		embeddings[i] = m.ParagraphEmbedding(p)
	}

	var pairs []RedundancyPair
	for i := 0; i < len(embeddings); i++ {
		if embeddings[i] == nil {
			continue
		}
		for j := i + 2; j < len(embeddings); j++ {
			if embeddings[j] == nil {
				continue
			}
			sim := shared.CosineSimilarity(embeddings[i], embeddings[j])
			if sim >= cfg.RedundancyThreshold {
				pairs = append(pairs, RedundancyPair{
					IndexA:     i,
					IndexB:     j,
					Similarity: sim,
				})
			}
		}
	}

	return RedundancyResult{
		Pairs:     pairs,
		PairCount: len(pairs),
	}
}
