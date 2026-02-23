package semantic

import (
	"github.com/inkcheck/config"
	"github.com/inkcheck/shared"
)

// EmotionalToneResult holds valence and arousal scores derived from
// word-vector seed projection (Russell circumplex model).
//
// Both scores are in [-1, +1]:
//   - Valence: −1 = very negative, 0 = neutral, +1 = very positive
//   - Arousal: −1 = very calm/passive, 0 = neutral, +1 = very excited/energetic
//
// CoveredWords is the count of content words that had embeddings in the model.
// If CoveredWords is 0 the model was nil or no words matched; scores will be 0.
type EmotionalToneResult struct {
	Valence      float64
	Arousal      float64
	CoveredWords int
}

// valencePositiveSeeds are high-valence (positive emotion) anchor words.
var valencePositiveSeeds = []string{
	"joy", "happy", "love", "excellent", "wonderful", "beautiful",
	"delight", "pleasure", "success", "triumph", "grateful", "hope",
	"cheerful", "positive", "bright", "amazing", "fantastic", "brilliant",
}

// valenceNegativeSeeds are low-valence (negative emotion) anchor words.
var valenceNegativeSeeds = []string{
	"sad", "terrible", "hate", "awful", "horrible", "disgusting",
	"grief", "pain", "failure", "misery", "angry", "fear",
	"despair", "negative", "dark", "dreadful", "pathetic",
}

// arousalHighSeeds are high-arousal (excited/active) anchor words.
var arousalHighSeeds = []string{
	"excited", "energetic", "thrilled", "passionate", "intense", "urgent",
	"wild", "frantic", "explosive", "electrifying", "fierce", "dynamic",
	"vibrant", "powerful", "bold", "dramatic", "stirring", "feverish",
}

// arousalLowSeeds are low-arousal (calm/passive) anchor words.
var arousalLowSeeds = []string{
	"calm", "peaceful", "quiet", "serene", "relaxed", "gentle",
	"slow", "still", "tranquil", "sleepy", "dull", "passive",
	"mellow", "subdued", "soft", "hushed", "languid", "placid",
}

// EmotionalTone estimates valence and arousal of the text by projecting
// content-word embeddings onto seed-word axes.
//
// The valence axis is the unit vector from the centroid of negative seed
// embeddings to the centroid of positive seed embeddings. Arousal uses the
// same construction with its own seed sets. Each content word is projected
// onto both axes via cosine similarity with the axis centroid pair, and the
// mean projection across all covered words is returned.
//
// Config is accepted for API consistency. Requires a non-nil ModelManager;
// returns zero result if model is nil or no words have embeddings.
func EmotionalTone(_ config.Config, m *ModelManager, text string) EmotionalToneResult {
	if m == nil {
		return EmotionalToneResult{}
	}

	prose := shared.ExtractProseText(text)
	words := shared.ListWords(prose)
	if len(words) == 0 {
		return EmotionalToneResult{}
	}

	// Build seed centroids
	posValCentroid := seedCentroid(m, valencePositiveSeeds)
	negValCentroid := seedCentroid(m, valenceNegativeSeeds)
	posArousalCentroid := seedCentroid(m, arousalHighSeeds)
	negArousalCentroid := seedCentroid(m, arousalLowSeeds)

	if posValCentroid == nil || negValCentroid == nil ||
		posArousalCentroid == nil || negArousalCentroid == nil {
		return EmotionalToneResult{}
	}

	valenceSum := 0.0
	arousalSum := 0.0
	covered := 0

	for _, w := range words {
		vec := m.Embedding(w)
		if vec == nil {
			continue
		}
		// Project onto valence axis: similarity to positive pole minus negative pole
		valPos := shared.CosineSimilarity(vec, posValCentroid)
		valNeg := shared.CosineSimilarity(vec, negValCentroid)
		valenceSum += valPos - valNeg

		// Project onto arousal axis
		aroPos := shared.CosineSimilarity(vec, posArousalCentroid)
		aroNeg := shared.CosineSimilarity(vec, negArousalCentroid)
		arousalSum += aroPos - aroNeg

		covered++
	}

	if covered == 0 {
		return EmotionalToneResult{}
	}

	valence := clamp(valenceSum/float64(covered), -1, 1)
	arousal := clamp(arousalSum/float64(covered), -1, 1)

	return EmotionalToneResult{
		Valence:      valence,
		Arousal:      arousal,
		CoveredWords: covered,
	}
}

// seedCentroid computes the average embedding of a list of seed words.
// Returns nil if no seed has an embedding in the model.
func seedCentroid(m *ModelManager, seeds []string) []float32 {
	var sum []float32
	count := 0
	for _, seed := range seeds {
		vec := m.Embedding(seed)
		if vec == nil {
			continue
		}
		if sum == nil {
			sum = make([]float32, len(vec))
		}
		for i, v := range vec {
			sum[i] += v
		}
		count++
	}
	if count == 0 {
		return nil
	}
	for i := range sum {
		sum[i] /= float32(count)
	}
	return sum
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
