// Package signature produces a 10-dimensional style signature for a text
// document. Each axis is a weighted composite of normalized sub-metrics,
// encoding tone of voice and writing style as a radar (decagon) chart.
//
// See requirements.md for the full specification.
package signature

import (
	"math"

	"github.com/inkcheck/rhetoric"
	"github.com/inkcheck/semantic"
"github.com/inkcheck/structure"
)

// AxisCount is the number of radar axes.
const AxisCount = 10

// Axis indices into the signature vector.
const (
	Formality = iota
	Confidence
	Rhythm
	Economy
	Precision
	Coherence
	Vocabulary
	Stance
	EmotionalTone
	TemporalOrientation
)

// AxisNames maps axis index to display name.
var AxisNames = [AxisCount]string{
	"formality", "confidence", "rhythm", "economy", "precision",
	"coherence", "vocabulary", "stance", "emotional_tone", "temporal_orientation",
}

// SubMetricValue holds raw and normalized values for a single sub-metric.
type SubMetricValue struct {
	Raw        float64
	Normalized float64
	Weight     float64
}

// AxisResult holds the composite score and sub-metric breakdown for one axis.
type AxisResult struct {
	Score      float64
	SubMetrics map[string]SubMetricValue
}

// Signature is a 10-dimensional style vector with per-axis breakdowns.
type Signature struct {
	Axes [AxisCount]AxisResult
}

// Vector returns the 10-dimensional score vector.
func (s *Signature) Vector() [AxisCount]float64 {
	var v [AxisCount]float64
	for i := range s.Axes {
		v[i] = s.Axes[i].Score
	}
	return v
}

// DocumentInfo holds contextual annotations (not radar axes).
type DocumentInfo struct {
	WordCount          int
	SentenceCount      int
	ParagraphCount     int
	ReadingTimeSeconds int
	ReadabilityFormula string
	ReadabilityScore   float64
	ReadabilityGrade   string
}

// RawMetrics collects all raw sub-metric sources needed to compute a signature.
// This decouples signature computation from the inkcheck.Result type so that
// callers can populate it from any source.
type RawMetrics struct {
	// Structure
	SentenceLengthCV float64
	ParagraphLengthCV float64
	OpenerDiversity  structure.SentenceOpenerDiversityResult
	SentenceType     structure.SentenceTypeResult

	// Rhetoric
	VoiceConsistency rhetoric.VoiceConsistencyResult
	Hedging          rhetoric.HedgingResult
	Specificity      rhetoric.SpecificityResult
	ClaimSupport     rhetoric.ClaimSupportResult
	ArgumentStructure rhetoric.ArgumentStructureResult
	Stance           rhetoric.StanceResult
	Contraction      rhetoric.ContractionResult
	Temporal         rhetoric.TemporalResult
	Economy          rhetoric.EconomyResult
	VocabSophistication rhetoric.VocabSophisticationResult
	TransitionWordDensity rhetoric.TransitionResult

	// Semantic (optional — requires word embeddings)
	TopicCoherence semantic.TopicCoherenceResult
	Redundancy     semantic.RedundancyResult
	EmotionalTone  semantic.EmotionalToneResult

	// Document-level
	WordCount      int
	SentenceCount  int
	ParagraphCount int
}

// Compute produces a Signature from raw metrics.
func Compute(m RawMetrics) Signature {
	defs := axesDefs()
	var sig Signature

	rawValues := extractRawValues(m)

	for i, def := range defs {
		subs := make(map[string]SubMetricValue, len(def.SubMetrics))
		norms := make([]float64, len(def.SubMetrics))
		weights := make([]float64, len(def.SubMetrics))

		for j, sm := range def.SubMetrics {
			raw := rawValues[sm.Name]
			n := normalize(raw, sm.Lo, sm.Hi, sm.Curve, sm.Invert)
			subs[sm.Name] = SubMetricValue{
				Raw:        raw,
				Normalized: n,
				Weight:     sm.Weight,
			}
			norms[j] = n
			weights[j] = sm.Weight
		}

		sig.Axes[i] = AxisResult{
			Score:      composite(norms, weights),
			SubMetrics: subs,
		}
	}

	return sig
}

// extractRawValues maps sub-metric names to their raw values from RawMetrics.
func extractRawValues(m RawMetrics) map[string]float64 {
	wordCount := float64(m.WordCount)

	// Transition density: total transitions per 100 words
	transitionDensity := 0.0
	if wordCount > 0 {
		transitionDensity = float64(m.TransitionWordDensity.Total) / wordCount * 100
	}

	// Specificity score: normalize mean from 1-5 scale to 0-1
	specificityNorm := 0.0
	if m.Specificity.Mean > 0 {
		specificityNorm = (m.Specificity.Mean - 1) / 4.0
	}

	// Redundancy score: approximate from pair count relative to paragraph count
	redundancyScore := 0.0
	if m.ParagraphCount > 2 {
		maxPairs := float64(m.ParagraphCount*(m.ParagraphCount-1)) / 2.0
		if maxPairs > 0 {
			redundancyScore = math.Min(float64(m.Redundancy.PairCount)/maxPairs, 1.0)
		}
	}

	// Claim-support ratio: normalize to 0-1 (ratio of 1.0 = balanced = 0.5 on scale)
	claimSupportNorm := 0.0
	if m.ClaimSupport.Ratio > 0 {
		// Ratio is support/claim; cap at 5.0. Map to 0-1 as ratio/5.
		claimSupportNorm = math.Min(m.ClaimSupport.Ratio/5.0, 1.0)
	}

	// Positive affect: approximate from valence (semantic model)
	// Valence is [-1, 1]; map positive half to [0, 1] for affect ratio
	positiveAffect := 0.0
	if m.EmotionalTone.Valence > 0 {
		positiveAffect = math.Min(m.EmotionalTone.Valence, 1.0) * 0.15
	}

	// Arousal level: map from [-1, 1] to [0, 1]
	arousalLevel := (m.EmotionalTone.Arousal + 1) / 2.0

	// Emotional intensity: approximate from arousal magnitude
	emotionalIntensity := math.Abs(m.EmotionalTone.Arousal) * 3.0

	// Empathy marker density: use audience awareness second-person as proxy
	// (This is a rough approximation; a dedicated empathy lexicon would be better)
	empathyDensity := 0.0
	if wordCount > 0 {
		empathyDensity = float64(m.Stance.SecondPerson) / wordCount * 100
	}

	// Past tense ratio: past tense verbs / total verb-like words (approximate)
	pastTenseRatio := 0.0
	totalVerbIndicators := m.Temporal.PastTenseCount + m.Temporal.FutureModalCount
	if totalVerbIndicators > 0 {
		pastTenseRatio = float64(m.Temporal.PastTenseCount) / float64(totalVerbIndicators)
	}

	// Active voice ratio = 1 - passive voice ratio
	activeVoiceRatio := 1.0 - m.VoiceConsistency.PassiveRatio

	return map[string]float64{
		// Formality
		"formal_word_ratio":  m.VocabSophistication.FormalWordRatio,
		"passive_voice_ratio": m.VoiceConsistency.PassiveRatio,
		"contraction_rate":    m.Contraction.Rate,

		// Confidence
		"hedging_density":   m.Hedging.Density,
		"modal_verb_density": m.Hedging.AssertiveModalDensity,
		"active_voice_ratio": activeVoiceRatio,

		// Rhythm
		"sentence_length_cv":   m.SentenceLengthCV,
		"paragraph_length_cv":  m.ParagraphLengthCV,
		"opener_diversity":     m.OpenerDiversity.Entropy,
		"sentence_type_entropy": m.SentenceType.Entropy,

		// Economy
		"avg_sentence_length":  m.Economy.AvgSentenceLength,
		"wordy_phrase_density": m.Economy.WordyPhraseDensity,
		"words_per_clause":     m.Economy.WordsPerClause,
		"syntactic_complexity": m.Economy.SubordinationIndex,

		// Precision
		"specificity_score":  specificityNorm,
		"vague_word_density": m.Specificity.VagueDensity,
		"redundancy_score":   redundancyScore,

		// Coherence
		"transition_density":  transitionDensity,
		"topic_coherence":     m.TopicCoherence.MeanSimilarity,
		"argument_structure":  m.ArgumentStructure.CoherenceScore,
		"claim_support_ratio": claimSupportNorm,

		// Vocabulary
		"mattr":               m.VocabSophistication.MATTR,
		"lexical_density":     m.VocabSophistication.LexicalDensity,
		"low_freq_word_ratio": m.VocabSophistication.LowFreqWordRatio,

		// Stance
		"reader_centricity": m.Stance.ReaderCentricity,

		// Emotional Tone
		"positive_affect_ratio":   positiveAffect,
		"emotional_intensity":     emotionalIntensity,
		"empathy_marker_density":  empathyDensity,
		"arousal_level":           arousalLevel,

		// Temporal Orientation
		"future_modal_density":      m.Temporal.FutureModalDensity,
		"past_tense_ratio":          pastTenseRatio,
		"evidential_marker_density": m.Temporal.EvidentialDensity,
		"aspiration_marker_density": m.Temporal.AspirationDensity,
	}
}

// FromAnalysis is a convenience constructor that builds RawMetrics from
// the individual analyzer results and document-level counts.
func FromAnalysis(
	sentLenCV, paraLenCV float64,
	opener structure.SentenceOpenerDiversityResult,
	sentType structure.SentenceTypeResult,
	voice rhetoric.VoiceConsistencyResult,
	hedging rhetoric.HedgingResult,
	specificity rhetoric.SpecificityResult,
	claimSupport rhetoric.ClaimSupportResult,
	argStructure rhetoric.ArgumentStructureResult,
	stance rhetoric.StanceResult,
	contraction rhetoric.ContractionResult,
	temporal rhetoric.TemporalResult,
	economy rhetoric.EconomyResult,
	vocab rhetoric.VocabSophisticationResult,
	transitions rhetoric.TransitionResult,
	topicCoherence semantic.TopicCoherenceResult,
	redundancy semantic.RedundancyResult,
	emotionalTone semantic.EmotionalToneResult,
	wordCount, sentenceCount, paragraphCount int,
) Signature {
	raw := RawMetrics{
		SentenceLengthCV:      sentLenCV,
		ParagraphLengthCV:     paraLenCV,
		OpenerDiversity:       opener,
		SentenceType:          sentType,
		VoiceConsistency:      voice,
		Hedging:               hedging,
		Specificity:           specificity,
		ClaimSupport:          claimSupport,
		ArgumentStructure:     argStructure,
		Stance:                stance,
		Contraction:           contraction,
		Temporal:              temporal,
		Economy:               economy,
		VocabSophistication:   vocab,
		TransitionWordDensity: transitions,
		TopicCoherence:        topicCoherence,
		Redundancy:            redundancy,
		EmotionalTone:         emotionalTone,
		WordCount:             wordCount,
		SentenceCount:         sentenceCount,
		ParagraphCount:        paragraphCount,
	}
	return Compute(raw)
}

// CosineSimilarity computes the cosine similarity between two signature vectors.
func CosineSimilarity(a, b [AxisCount]float64) float64 {
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
