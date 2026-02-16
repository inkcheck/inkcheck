package semantic

import (
	"github.com/inkcheck/config"
	"github.com/inkcheck/shared"
)

// InformationNoveltyCurveResult holds per-paragraph novelty scores.
type InformationNoveltyCurveResult struct {
	// NoveltyScores is per-paragraph novelty (1 - max similarity to any prior paragraph).
	// First paragraph always has novelty 1.0.
	NoveltyScores []float64
	// MeanNovelty is the average novelty score.
	MeanNovelty float64
	// CV is the coefficient of variation of novelty scores.
	CV float64
}

// InformationNoveltyCurve computes per-paragraph novelty as
// 1 - max(similarity to all prior paragraphs). Config is accepted for API
// consistency with other metric functions.
func InformationNoveltyCurve(_ config.Config, m *ModelManager, text string) InformationNoveltyCurveResult {
	if m == nil {
		return InformationNoveltyCurveResult{}
	}
	paragraphs := shared.SplitParagraphs(text)
	if len(paragraphs) < 2 {
		return InformationNoveltyCurveResult{}
	}

	embeddings := make([][]float32, len(paragraphs))
	for i, p := range paragraphs {
		embeddings[i] = m.ParagraphEmbedding(p)
	}

	novelty := make([]float64, len(paragraphs))
	novelty[0] = 1.0 // First paragraph is fully novel

	for i := 1; i < len(paragraphs); i++ {
		if embeddings[i] == nil {
			novelty[i] = 1.0
			continue
		}
		maxSim := 0.0
		for j := 0; j < i; j++ {
			if embeddings[j] == nil {
				continue
			}
			sim := shared.CosineSimilarity(embeddings[j], embeddings[i])
			if sim > maxSim {
				maxSim = sim
			}
		}
		novelty[i] = 1.0 - maxSim
	}

	return InformationNoveltyCurveResult{
		NoveltyScores: novelty,
		MeanNovelty:   shared.Mean(novelty),
		CV:            shared.CoefficientOfVariation(novelty),
	}
}
