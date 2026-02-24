package signature_test

import (
	"math"
	"testing"

	"github.com/inkcheck/rhetoric"
	"github.com/inkcheck/semantic"
	"github.com/inkcheck/signature"
	"github.com/inkcheck/structure"
)

func TestCompute_ZeroMetrics(t *testing.T) {
	sig := signature.Compute(signature.RawMetrics{})

	for i, axis := range sig.Axes {
		if axis.Score < 0 || axis.Score > 1 {
			t.Errorf("axis %d (%s) score out of [0,1]: %v",
				i, signature.AxisNames[i], axis.Score)
		}
		for name, sm := range axis.SubMetrics {
			if sm.Normalized < 0 || sm.Normalized > 1 {
				t.Errorf("axis %d sub-metric %s normalized out of [0,1]: %v",
					i, name, sm.Normalized)
			}
		}
	}
}

func TestCompute_TypicalMetrics(t *testing.T) {
	raw := signature.RawMetrics{
		SentenceLengthCV:  0.41,
		ParagraphLengthCV: 0.62,
		OpenerDiversity:   structure.SentenceOpenerDiversityResult{Entropy: 1.83},
		SentenceType:      structure.SentenceTypeResult{Entropy: 0.72},
		VoiceConsistency:  rhetoric.VoiceConsistencyResult{PassiveRatio: 0.19},
		Hedging:           rhetoric.HedgingResult{Density: 1.8, AssertiveModalDensity: 2.1},
		Specificity:       rhetoric.SpecificityResult{Mean: 3.44, VagueDensity: 1.4},
		ClaimSupport:      rhetoric.ClaimSupportResult{Ratio: 0.74},
		ArgumentStructure: rhetoric.ArgumentStructureResult{CoherenceScore: 0.71},
		Stance:            rhetoric.StanceResult{ReaderCentricity: 0.031, SecondPerson: 5},
		Contraction:       rhetoric.ContractionResult{Rate: 0.01},
		Temporal: rhetoric.TemporalResult{
			FutureModalDensity: 1.4,
			PastTenseCount:     8,
			FutureModalCount:   4,
			EvidentialDensity:  0.7,
			AspirationDensity: 0.9,
		},
		Economy: rhetoric.EconomyResult{
			AvgSentenceLength:  16.3,
			WordyPhraseDensity: 1.2,
			WordsPerClause:     9.1,
			SubordinationIndex: 0.22,
		},
		VocabSophistication: rhetoric.VocabSophisticationResult{
			MATTR:            0.72,
			FormalWordRatio:  0.08,
			LexicalDensity:   0.51,
			LowFreqWordRatio: 0.14,
		},
		TransitionWordDensity: rhetoric.TransitionResult{Total: 27},
		TopicCoherence:        semantic.TopicCoherenceResult{MeanSimilarity: 0.68},
		EmotionalTone:         semantic.EmotionalToneResult{Valence: 0.2, Arousal: 0.1},
		WordCount:             847,
		SentenceCount:         52,
		ParagraphCount:        9,
	}

	sig := signature.Compute(raw)

	for i, axis := range sig.Axes {
		if axis.Score < 0 || axis.Score > 1 {
			t.Errorf("axis %d (%s) score out of [0,1]: %v",
				i, signature.AxisNames[i], axis.Score)
		}
		if len(axis.SubMetrics) == 0 {
			t.Errorf("axis %d (%s) has no sub-metrics",
				i, signature.AxisNames[i])
		}
	}
}

func TestVector(t *testing.T) {
	sig := signature.Compute(signature.RawMetrics{})
	v := sig.Vector()
	if len(v) != signature.AxisCount {
		t.Errorf("Vector length = %d, want %d", len(v), signature.AxisCount)
	}
}

func TestCosineSimilarity_Identical(t *testing.T) {
	a := [signature.AxisCount]float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5}
	got := signature.CosineSimilarity(a, a)
	if math.Abs(got-1.0) > 1e-9 {
		t.Errorf("identical vectors: cosine similarity = %v, want 1.0", got)
	}
}

func TestCosineSimilarity_Zero(t *testing.T) {
	a := [signature.AxisCount]float64{}
	b := [signature.AxisCount]float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	got := signature.CosineSimilarity(a, b)
	if got != 0 {
		t.Errorf("zero vector: cosine similarity = %v, want 0", got)
	}
}

func TestAxisNames(t *testing.T) {
	expected := []string{
		"formality", "confidence", "rhythm", "economy", "precision",
		"coherence", "vocabulary", "stance", "emotional_tone", "temporal_orientation",
	}
	for i, name := range expected {
		if signature.AxisNames[i] != name {
			t.Errorf("AxisNames[%d] = %q, want %q", i, signature.AxisNames[i], name)
		}
	}
}
