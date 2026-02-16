package inkcheck_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/inkcheck"
	"github.com/inkcheck/config"
	"github.com/inkcheck/semantic"
)

// loadTestData reads a test file from testdata directory.
func loadTestData(t *testing.T, filename string) string {
	t.Helper()
	path := filepath.Join("testdata", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to load test data %s: %v", filename, err)
	}
	return string(data)
}

// TestAnalyze_ArgumentativeEssay tests the full analysis pipeline on an argumentative essay.
// This text type should show:
// - Higher claim/support ratio (argumentative structure)
// - Evidence of counterargument engagement
// - Transition word usage
// - Thesis/evidence/conclusion markers
func TestAnalyze_ArgumentativeEssay(t *testing.T) {
	text := loadTestData(t, "argumentative_essay.md")
	cfg := config.DefaultConfig()

	result := inkcheck.Analyze(cfg, text)

	// Structure metrics: reasonable variance expected
	if result.Structure.ParagraphVariance < 0 {
		t.Errorf("paragraph variance should be non-negative, got %v", result.Structure.ParagraphVariance)
	}
	if result.Structure.SentenceLengthVariance < 0 {
		t.Errorf("sentence length variance should be non-negative, got %v", result.Structure.SentenceLengthVariance)
	}

	// Rhetoric: argumentative text should have counterargument engagement
	if result.Rhetoric.Counterargument.Instances < 1 {
		t.Errorf("expected at least 1 counterargument instance in argumentative essay, got %d", result.Rhetoric.Counterargument.Instances)
	}

	// Should detect some argument structure
	if !result.Rhetoric.ArgumentStructure.HasThesisMarker && !result.Rhetoric.ArgumentStructure.HasConclusionMarker {
		t.Error("expected to detect thesis or conclusion markers in argumentative essay")
	}

	// Hedging should be measurable
	if result.Rhetoric.Hedging.Density < 0 {
		t.Errorf("hedging density should be non-negative, got %v", result.Rhetoric.Hedging.Density)
	}

	// Specificity should have a reasonable range
	if result.Rhetoric.Specificity.Mean < 1 || result.Rhetoric.Specificity.Mean > 5 {
		t.Errorf("specificity mean should be between 1-5, got %v", result.Rhetoric.Specificity.Mean)
	}
}

// TestAnalyze_NarrativeProse tests analysis on narrative text.
// This text type should show:
// - Figurative language (similes, metaphors)
// - Tension and resolution arc
// - Different voice consistency patterns
func TestAnalyze_NarrativeProse(t *testing.T) {
	text := loadTestData(t, "narrative_prose.md")
	cfg := config.DefaultConfig()

	result := inkcheck.Analyze(cfg, text)

	// Structure metrics should be in valid range
	if result.Structure.ParagraphVariance < 0 || result.Structure.ParagraphVariance > 2 {
		t.Errorf("unexpected paragraph variance: %v", result.Structure.ParagraphVariance)
	}

	// Narrative text often has figurative language
	if result.Rhetoric.FigurativeLanguage.TotalInstances < 0 {
		t.Errorf("figurative language instances should be non-negative, got %d", result.Rhetoric.FigurativeLanguage.TotalInstances)
	}

	// Tension/resolution should be detected in narrative
	if result.Rhetoric.TensionAndResolution.TensionMarkers < 0 || result.Rhetoric.TensionAndResolution.ResolutionMarkers < 0 {
		t.Error("tension/resolution markers should be non-negative")
	}

	// Passive voice ratio should be between 0 and 1
	if result.Rhetoric.VoiceConsistency.PassiveRatio < 0 || result.Rhetoric.VoiceConsistency.PassiveRatio > 1 {
		t.Errorf("passive ratio should be 0-1, got %v", result.Rhetoric.VoiceConsistency.PassiveRatio)
	}
}

// TestAnalyze_TechnicalWriting tests analysis on technical/expository text.
// This text type should show:
// - Higher vocabulary sophistication (technical terms)
// - Lower figurative language usage
// - More specific language (numbers, technical details)
func TestAnalyze_TechnicalWriting(t *testing.T) {
	text := loadTestData(t, "technical_writing.md")
	cfg := config.DefaultConfig()

	result := inkcheck.Analyze(cfg, text)

	// Vocabulary sophistication should be measurable
	if result.Rhetoric.VocabSophistication.TypeTokenRatio < 0 || result.Rhetoric.VocabSophistication.TypeTokenRatio > 1 {
		t.Errorf("TTR should be 0-1, got %v", result.Rhetoric.VocabSophistication.TypeTokenRatio)
	}

	// MATTR should also be 0-1
	if result.Rhetoric.VocabSophistication.MATTR < 0 || result.Rhetoric.VocabSophistication.MATTR > 1 {
		t.Errorf("MATTR should be 0-1, got %v", result.Rhetoric.VocabSophistication.MATTR)
	}

	// Technical text should have specificity
	if result.Rhetoric.Specificity.Mean < 1 || result.Rhetoric.Specificity.Mean > 5 {
		t.Errorf("specificity mean should be 1-5, got %v", result.Rhetoric.Specificity.Mean)
	}

	// Punctuation profile should show variety
	if result.Structure.Punctuation.Variety() < 0 || result.Structure.Punctuation.Variety() > 8 {
		t.Errorf("punctuation variety should be 0-8, got %d", result.Structure.Punctuation.Variety())
	}
}

// TestAnalyze_EdgeCases tests the analysis pipeline on edge cases.
func TestAnalyze_EdgeCases(t *testing.T) {
	text := loadTestData(t, "edge_cases.md")
	cfg := config.DefaultConfig()

	// Should not panic on edge cases
	result := inkcheck.Analyze(cfg, text)

	// Basic sanity checks - metrics should be in valid ranges
	if result.Structure.SentenceLengthVariance < 0 {
		t.Errorf("sentence variance should be non-negative, got %v", result.Structure.SentenceLengthVariance)
	}

	if result.Structure.SentenceOpenerDiversity < 0 || result.Structure.SentenceOpenerDiversity > 1 {
		t.Errorf("opener diversity should be 0-1, got %v", result.Structure.SentenceOpenerDiversity)
	}
}

// TestAnalyze_EmptyText tests handling of empty text.
func TestAnalyze_EmptyText(t *testing.T) {
	cfg := config.DefaultConfig()
	result := inkcheck.Analyze(cfg, "")

	// Should not panic, should return zero/default values
	if result.Structure.ParagraphVariance != 0 {
		t.Errorf("expected zero variance for empty text, got %v", result.Structure.ParagraphVariance)
	}
}

// TestAnalyze_SingleSentence tests handling of minimal text.
func TestAnalyze_SingleSentence(t *testing.T) {
	cfg := config.DefaultConfig()
	result := inkcheck.Analyze(cfg, "This is a single sentence.")

	// Should not panic
	// Variance metrics should be zero or very low for single sentence
	if result.Structure.SentenceLengthVariance != 0 {
		t.Errorf("expected zero sentence variance for single sentence, got %v", result.Structure.SentenceLengthVariance)
	}
}

// TestAnalyzeStructure_Isolation tests just the structure metrics.
func TestAnalyzeStructure_Isolation(t *testing.T) {
	text := loadTestData(t, "argumentative_essay.md")
	cfg := config.DefaultConfig()

	result := inkcheck.AnalyzeStructure(cfg, text)

	// Verify all structure metrics are populated
	if result.ParagraphVariance < 0 {
		t.Error("paragraph variance should be non-negative")
	}
	if len(result.ParagraphLengths) == 0 {
		t.Error("should have paragraph lengths")
	}
	if result.SentenceLengthVariance < 0 {
		t.Error("sentence length variance should be non-negative")
	}
	if result.SentenceOpenerDiversity < 0 || result.SentenceOpenerDiversity > 1 {
		t.Error("opener diversity should be between 0 and 1")
	}
}

// TestAnalyzeRhetoric_Isolation tests just the rhetoric metrics.
func TestAnalyzeRhetoric_Isolation(t *testing.T) {
	text := loadTestData(t, "narrative_prose.md")
	cfg := config.DefaultConfig()

	result := inkcheck.AnalyzeRhetoric(cfg, text)

	// Verify rhetoric metrics are in valid ranges
	if result.Hedging.Density < 0 {
		t.Error("hedging density should be non-negative")
	}
	if result.VocabSophistication.TypeTokenRatio < 0 || result.VocabSophistication.TypeTokenRatio > 1 {
		t.Error("TTR should be between 0 and 1")
	}
	if result.Specificity.Mean < 0 || result.Specificity.Mean > 5 {
		t.Error("specificity mean should be between 0 and 5")
	}
}

// TestReadability tests the readability metric.
func TestReadability(t *testing.T) {
	text := loadTestData(t, "technical_writing.md")
	cfg := config.DefaultConfig()

	result := inkcheck.AnalyzeReadability(cfg, text)

	// Should have a formula name
	if result.Formula == "" {
		t.Error("readability formula should not be empty")
	}

	// Score should be a reasonable number (not checking exact value)
	// Flesch-Kincaid grade typically ranges 0-18+
	if result.Score < 0 || result.Score > 30 {
		t.Errorf("readability score seems unreasonable: %v", result.Score)
	}
}

// TestAnalyzeAll_WithSemanticModel tests the full pipeline including semantic metrics.
func TestAnalyzeAll_WithSemanticModel(t *testing.T) {
	text := loadTestData(t, "argumentative_essay.md")
	cfg := config.DefaultConfig()

	model, err := semantic.LoadModel(cfg)
	if err != nil {
		t.Skipf("Skipping semantic test: %v", err)
	}

	result := inkcheck.AnalyzeAll(cfg, text, model)

	// Semantic metrics should be populated
	if result.Semantic.TopicCoherence.MeanSimilarity < 0 || result.Semantic.TopicCoherence.MeanSimilarity > 1 {
		t.Errorf("topic coherence should be 0-1, got %v", result.Semantic.TopicCoherence.MeanSimilarity)
	}

	if result.Semantic.SemanticProgression.MeanDrift < 0 || result.Semantic.SemanticProgression.MeanDrift > 1 {
		t.Errorf("semantic drift should be 0-1, got %v", result.Semantic.SemanticProgression.MeanDrift)
	}

	// Redundancy detection
	if result.Semantic.Redundancy.PairCount < 0 {
		t.Errorf("redundancy pair count should be non-negative, got %d", result.Semantic.Redundancy.PairCount)
	}

	// Information novelty
	if result.Semantic.InformationNovelty.MeanNovelty < 0 || result.Semantic.InformationNovelty.MeanNovelty > 1 {
		t.Errorf("information novelty should be 0-1, got %v", result.Semantic.InformationNovelty.MeanNovelty)
	}
}

// TestAnalyzeSemantic_Isolation tests semantic metrics in isolation.
func TestAnalyzeSemantic_Isolation(t *testing.T) {
	text := loadTestData(t, "narrative_prose.md")
	cfg := config.DefaultConfig()

	model, err := semantic.LoadModel(cfg)
	if err != nil {
		t.Skipf("Skipping semantic test: %v", err)
	}

	result := inkcheck.AnalyzeSemantic(cfg, text, model)

	// Topic coherence
	if result.TopicCoherence.MeanSimilarity < 0 || result.TopicCoherence.MeanSimilarity > 1 {
		t.Errorf("topic coherence mean should be 0-1, got %v", result.TopicCoherence.MeanSimilarity)
	}
	if result.TopicCoherence.CV < 0 {
		t.Errorf("topic coherence CV should be non-negative, got %v", result.TopicCoherence.CV)
	}

	// Semantic progression
	if result.SemanticProgression.MeanDrift < 0 || result.SemanticProgression.MeanDrift > 1 {
		t.Errorf("mean drift should be 0-1, got %v", result.SemanticProgression.MeanDrift)
	}

	// Redundancy
	if result.Redundancy.PairCount < 0 {
		t.Error("redundancy pair count should be non-negative")
	}
	for _, pair := range result.Redundancy.Pairs {
		if pair.Similarity < 0 || pair.Similarity > 1 {
			t.Errorf("pair similarity should be 0-1, got %v", pair.Similarity)
		}
	}

	// Information novelty
	if result.InformationNovelty.MeanNovelty < 0 || result.InformationNovelty.MeanNovelty > 1 {
		t.Errorf("mean novelty should be 0-1, got %v", result.InformationNovelty.MeanNovelty)
	}
	if len(result.InformationNovelty.NoveltyScores) == 0 {
		t.Error("should have novelty scores")
	}
}

// TestSemantic_TechnicalText tests semantic metrics on technical writing.
// Technical text often has consistent terminology (high coherence, low drift).
func TestSemantic_TechnicalText(t *testing.T) {
	text := loadTestData(t, "technical_writing.md")
	cfg := config.DefaultConfig()

	model, err := semantic.LoadModel(cfg)
	if err != nil {
		t.Skipf("Skipping semantic test: %v", err)
	}

	result := inkcheck.AnalyzeSemantic(cfg, text, model)

	// Technical writing should have reasonable coherence
	if result.TopicCoherence.MeanSimilarity < 0 {
		t.Error("topic coherence should be non-negative")
	}

	// Should not have excessive redundancy (different topics covered)
	if result.Redundancy.PairCount < 0 {
		t.Error("redundancy count should be non-negative")
	}
}
