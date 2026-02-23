package semantic

import (
	"math"
	"testing"

	"github.com/inkcheck/config"
)

// newTestModel builds a minimal ModelManager from a word→vector map.
// Vectors are pre-normalized to unit length for predictable cosine sims.
func newTestModel(vecs map[string][]float32) *ModelManager {
	dim := 0
	for _, v := range vecs {
		dim = len(v)
		break
	}
	return &ModelManager{embeddings: vecs, dim: dim}
}

// unit normalises a []float32 slice in-place and returns it.
func unit(v []float32) []float32 {
	var norm float64
	for _, x := range v {
		norm += float64(x) * float64(x)
	}
	norm = math.Sqrt(norm)
	if norm == 0 {
		return v
	}
	for i := range v {
		v[i] = float32(float64(v[i]) / norm)
	}
	return v
}

// buildTestModel creates a 4-dimensional model with just enough words to
// exercise the seed-projection logic without a real embedding file.
//
// Dimensions (conceptual):
//
//	[0] positive valence pole
//	[1] negative valence pole
//	[2] high arousal pole
//	[3] low arousal pole
func buildTestModel() *ModelManager {
	vecs := map[string][]float32{
		// ---- valence positive seeds ----
		"joy":       unit([]float32{1, 0, 0, 0}),
		"happy":     unit([]float32{1, 0, 0, 0}),
		"love":      unit([]float32{1, 0, 0, 0}),
		"excellent": unit([]float32{1, 0, 0, 0}),
		"wonderful": unit([]float32{1, 0, 0, 0}),
		"beautiful": unit([]float32{1, 0, 0, 0}),
		"delight":   unit([]float32{1, 0, 0, 0}),
		"pleasure":  unit([]float32{1, 0, 0, 0}),
		"success":   unit([]float32{1, 0, 0, 0}),
		"triumph":   unit([]float32{1, 0, 0, 0}),
		"grateful":  unit([]float32{1, 0, 0, 0}),
		"hope":      unit([]float32{1, 0, 0, 0}),
		"cheerful":  unit([]float32{1, 0, 0, 0}),
		"positive":  unit([]float32{1, 0, 0, 0}),
		"bright":    unit([]float32{1, 0, 0, 0}),
		"amazing":   unit([]float32{1, 0, 0, 0}),
		"fantastic": unit([]float32{1, 0, 0, 0}),
		"brilliant": unit([]float32{1, 0, 0, 0}),
		// ---- valence negative seeds ----
		"sad":        unit([]float32{0, 1, 0, 0}),
		"terrible":   unit([]float32{0, 1, 0, 0}),
		"hate":       unit([]float32{0, 1, 0, 0}),
		"awful":      unit([]float32{0, 1, 0, 0}),
		"horrible":   unit([]float32{0, 1, 0, 0}),
		"disgusting": unit([]float32{0, 1, 0, 0}),
		"grief":      unit([]float32{0, 1, 0, 0}),
		"pain":       unit([]float32{0, 1, 0, 0}),
		"failure":    unit([]float32{0, 1, 0, 0}),
		"misery":     unit([]float32{0, 1, 0, 0}),
		"angry":      unit([]float32{0, 1, 0, 0}),
		"fear":       unit([]float32{0, 1, 0, 0}),
		"despair":    unit([]float32{0, 1, 0, 0}),
		"negative":   unit([]float32{0, 1, 0, 0}),
		"dark":       unit([]float32{0, 1, 0, 0}),
		"dreadful":   unit([]float32{0, 1, 0, 0}),
		"pathetic":   unit([]float32{0, 1, 0, 0}),
		// ---- arousal high seeds ----
		"excited":       unit([]float32{0, 0, 1, 0}),
		"energetic":     unit([]float32{0, 0, 1, 0}),
		"thrilled":      unit([]float32{0, 0, 1, 0}),
		"passionate":    unit([]float32{0, 0, 1, 0}),
		"intense":       unit([]float32{0, 0, 1, 0}),
		"urgent":        unit([]float32{0, 0, 1, 0}),
		"wild":          unit([]float32{0, 0, 1, 0}),
		"frantic":       unit([]float32{0, 0, 1, 0}),
		"explosive":     unit([]float32{0, 0, 1, 0}),
		"electrifying":  unit([]float32{0, 0, 1, 0}),
		"fierce":        unit([]float32{0, 0, 1, 0}),
		"dynamic":       unit([]float32{0, 0, 1, 0}),
		"vibrant":       unit([]float32{0, 0, 1, 0}),
		"powerful":      unit([]float32{0, 0, 1, 0}),
		"bold":          unit([]float32{0, 0, 1, 0}),
		"dramatic":      unit([]float32{0, 0, 1, 0}),
		"stirring":      unit([]float32{0, 0, 1, 0}),
		"feverish":      unit([]float32{0, 0, 1, 0}),
		// ---- arousal low seeds ----
		"calm":      unit([]float32{0, 0, 0, 1}),
		"peaceful":  unit([]float32{0, 0, 0, 1}),
		"quiet":     unit([]float32{0, 0, 0, 1}),
		"serene":    unit([]float32{0, 0, 0, 1}),
		"relaxed":   unit([]float32{0, 0, 0, 1}),
		"gentle":    unit([]float32{0, 0, 0, 1}),
		"slow":      unit([]float32{0, 0, 0, 1}),
		"still":     unit([]float32{0, 0, 0, 1}),
		"tranquil":  unit([]float32{0, 0, 0, 1}),
		"sleepy":    unit([]float32{0, 0, 0, 1}),
		"dull":      unit([]float32{0, 0, 0, 1}),
		"passive":   unit([]float32{0, 0, 0, 1}),
		"mellow":    unit([]float32{0, 0, 0, 1}),
		"subdued":   unit([]float32{0, 0, 0, 1}),
		"soft":      unit([]float32{0, 0, 0, 1}),
		"hushed":    unit([]float32{0, 0, 0, 1}),
		"languid":   unit([]float32{0, 0, 0, 1}),
		"placid":    unit([]float32{0, 0, 0, 1}),
		// ---- test content words ----
		// Pure positive valence: aligned with dim 0
		"wonderful_word": unit([]float32{1, 0, 0, 0}),
		// Pure negative valence: aligned with dim 1
		"dreadful_word": unit([]float32{0, 1, 0, 0}),
		// Pure high arousal: aligned with dim 2
		"energetic_word": unit([]float32{0, 0, 1, 0}),
		// Pure low arousal: aligned with dim 3
		"tranquil_word": unit([]float32{0, 0, 0, 1}),
	}
	return newTestModel(vecs)
}

func TestEmotionalTone_NilModel(t *testing.T) {
	cfg := config.DefaultConfig()
	result := EmotionalTone(cfg, nil, "This is a happy, wonderful day!")

	if result.Valence != 0 || result.Arousal != 0 || result.CoveredWords != 0 {
		t.Errorf("expected zero result for nil model, got %+v", result)
	}
}

func TestEmotionalTone_EmptyText(t *testing.T) {
	m := buildTestModel()
	cfg := config.DefaultConfig()
	result := EmotionalTone(cfg, m, "")

	if result.Valence != 0 || result.Arousal != 0 || result.CoveredWords != 0 {
		t.Errorf("expected zero result for empty text, got %+v", result)
	}
}

func TestEmotionalTone_PositiveValence(t *testing.T) {
	m := buildTestModel()
	cfg := config.DefaultConfig()
	// Text whose only model-covered word is aligned with the positive-valence pole.
	result := EmotionalTone(cfg, m, "wonderful_word")

	if result.CoveredWords == 0 {
		t.Fatal("expected CoveredWords > 0")
	}
	if result.Valence <= 0 {
		t.Errorf("expected positive Valence, got %v", result.Valence)
	}
}

func TestEmotionalTone_NegativeValence(t *testing.T) {
	m := buildTestModel()
	cfg := config.DefaultConfig()
	result := EmotionalTone(cfg, m, "dreadful_word")

	if result.CoveredWords == 0 {
		t.Fatal("expected CoveredWords > 0")
	}
	if result.Valence >= 0 {
		t.Errorf("expected negative Valence, got %v", result.Valence)
	}
}

func TestEmotionalTone_HighArousal(t *testing.T) {
	m := buildTestModel()
	cfg := config.DefaultConfig()
	result := EmotionalTone(cfg, m, "energetic_word")

	if result.CoveredWords == 0 {
		t.Fatal("expected CoveredWords > 0")
	}
	if result.Arousal <= 0 {
		t.Errorf("expected positive Arousal, got %v", result.Arousal)
	}
}

func TestEmotionalTone_LowArousal(t *testing.T) {
	m := buildTestModel()
	cfg := config.DefaultConfig()
	result := EmotionalTone(cfg, m, "tranquil_word")

	if result.CoveredWords == 0 {
		t.Fatal("expected CoveredWords > 0")
	}
	if result.Arousal >= 0 {
		t.Errorf("expected negative Arousal, got %v", result.Arousal)
	}
}

func TestEmotionalTone_ScoreRange(t *testing.T) {
	m := buildTestModel()
	cfg := config.DefaultConfig()
	texts := []string{
		"wonderful_word dreadful_word",
		"energetic_word tranquil_word",
		"wonderful_word",
		"dreadful_word",
	}
	for _, text := range texts {
		result := EmotionalTone(cfg, m, text)
		if result.Valence < -1 || result.Valence > 1 {
			t.Errorf("Valence out of range [-1,1]: %v for %q", result.Valence, text)
		}
		if result.Arousal < -1 || result.Arousal > 1 {
			t.Errorf("Arousal out of range [-1,1]: %v for %q", result.Arousal, text)
		}
	}
}

func TestEmotionalTone_NonOrthogonalSeeds(t *testing.T) {
	// Model where seed words have overlapping dimensions, simulating
	// real embeddings where emotional categories aren't perfectly separable.
	vecs := map[string][]float32{
		// Positive valence seeds with some arousal bleed
		"joy":       unit([]float32{0.9, -0.1, 0.3, 0}),
		"happy":     unit([]float32{0.8, -0.2, 0.2, 0.1}),
		"love":      unit([]float32{0.85, -0.1, 0.1, 0.2}),
		"excellent": unit([]float32{0.9, -0.15, 0.2, 0}),
		// Negative valence seeds with some arousal bleed
		"sad":       unit([]float32{-0.1, 0.9, 0, 0.3}),
		"terrible":  unit([]float32{-0.2, 0.85, 0.2, 0.1}),
		"hate":      unit([]float32{-0.15, 0.8, 0.3, 0}),
		"awful":     unit([]float32{-0.1, 0.9, 0.1, 0.1}),
		// High arousal seeds
		"excited":   unit([]float32{0.2, 0, 0.9, -0.1}),
		"energetic": unit([]float32{0.1, 0, 0.85, -0.15}),
		"thrilled":  unit([]float32{0.3, 0, 0.8, -0.1}),
		"intense":   unit([]float32{0.1, 0.1, 0.9, -0.1}),
		// Low arousal seeds
		"calm":     unit([]float32{0.1, 0, -0.1, 0.9}),
		"peaceful": unit([]float32{0.15, 0, -0.15, 0.85}),
		"quiet":    unit([]float32{0, 0.1, -0.1, 0.9}),
		"serene":   unit([]float32{0.1, 0, -0.05, 0.9}),
		// Content word (not a seed): clearly positive and high-arousal
		"celebration": unit([]float32{0.7, -0.1, 0.6, -0.1}),
	}
	m := newTestModel(vecs)
	cfg := config.DefaultConfig()

	result := EmotionalTone(cfg, m, "celebration")

	if result.CoveredWords == 0 {
		t.Fatal("expected CoveredWords > 0")
	}
	if result.Valence <= 0 {
		t.Errorf("expected positive Valence for 'celebration', got %v", result.Valence)
	}
	if result.Arousal <= 0 {
		t.Errorf("expected positive Arousal for 'celebration', got %v", result.Arousal)
	}
	if result.Valence < -1 || result.Valence > 1 {
		t.Errorf("Valence out of range: %v", result.Valence)
	}
	if result.Arousal < -1 || result.Arousal > 1 {
		t.Errorf("Arousal out of range: %v", result.Arousal)
	}
}

func TestEmotionalTone_MissingSeedModel(t *testing.T) {
	// A model with no seed words at all — centroids will be nil, so result should be zero.
	vecs := map[string][]float32{
		"randomword": unit([]float32{0.5, 0.5, 0.5, 0.5}),
	}
	m := newTestModel(vecs)
	cfg := config.DefaultConfig()
	result := EmotionalTone(cfg, m, "randomword")

	if result.Valence != 0 || result.Arousal != 0 {
		t.Errorf("expected zero scores for model with no seed words, got v=%v a=%v", result.Valence, result.Arousal)
	}
}

func TestEmotionalTone_NoModelCoverage(t *testing.T) {
	// A model with only seed words — content text has no coverage.
	m := buildTestModel()
	cfg := config.DefaultConfig()
	result := EmotionalTone(cfg, m, "xyzzy quux frobnicator blarg")

	if result.CoveredWords != 0 {
		t.Errorf("expected CoveredWords = 0 for unknown words, got %d", result.CoveredWords)
	}
	if result.Valence != 0 || result.Arousal != 0 {
		t.Errorf("expected zero scores for no coverage, got v=%v a=%v", result.Valence, result.Arousal)
	}
}
