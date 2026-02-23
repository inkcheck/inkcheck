package structure_test

import (
	"testing"

	"github.com/inkcheck/config"
	"github.com/inkcheck/structure"
)

func TestSentenceOpenerDiversity_HighDiversity(t *testing.T) {
	// Each sentence starts with a distinct word.
	text := "The cat sat down. A dog ran fast. His owner smiled warmly. " +
		"Running quickly, she left. Suddenly everything changed. " +
		"Many people noticed it. Few understood why."
	cfg := config.DefaultConfig()
	result := structure.SentenceOpenerDiversity(cfg, text)

	if result.Ratio <= 0 {
		t.Errorf("expected Ratio > 0, got %v", result.Ratio)
	}
	if result.Ratio > 1 {
		t.Errorf("expected Ratio <= 1, got %v", result.Ratio)
	}
	if result.Entropy <= 0 {
		t.Errorf("expected Entropy > 0 for diverse openers, got %v", result.Entropy)
	}
}

func TestSentenceOpenerDiversity_LowDiversity(t *testing.T) {
	// All sentences start with the same two-word opener "Furthermore the".
	// Default OpenerWordCount is 2, so only one distinct opener exists.
	text := "Furthermore the results show improvement. Furthermore the data confirms the trend. " +
		"Furthermore the evidence supports the claim. Furthermore the analysis validates our approach. " +
		"Furthermore the findings align with expectations."
	cfg := config.DefaultConfig()
	result := structure.SentenceOpenerDiversity(cfg, text)

	if result.Ratio >= 0.5 {
		t.Errorf("expected Ratio < 0.5 for repetitive openers, got %v", result.Ratio)
	}
	// Entropy should be low (all openers are the same → entropy = 0)
	if result.Entropy > 0.1 {
		t.Errorf("expected Entropy ≈ 0 for identical openers, got %v", result.Entropy)
	}
}

func TestSentenceOpenerDiversity_Empty(t *testing.T) {
	cfg := config.DefaultConfig()
	result := structure.SentenceOpenerDiversity(cfg, "")

	if result.Ratio != 0 {
		t.Errorf("expected Ratio = 0 for empty text, got %v", result.Ratio)
	}
	if result.Entropy != 0 {
		t.Errorf("expected Entropy = 0 for empty text, got %v", result.Entropy)
	}
}

func TestSentenceOpenerDiversity_SingleSentence(t *testing.T) {
	cfg := config.DefaultConfig()
	result := structure.SentenceOpenerDiversity(cfg, "Just one sentence here.")

	if result.Ratio != 1.0 {
		t.Errorf("expected Ratio = 1.0 for single sentence, got %v", result.Ratio)
	}
	if result.Entropy != 0.0 {
		t.Errorf("expected Entropy = 0.0 for single sentence, got %v", result.Entropy)
	}
}

func TestSentenceOpenerDiversity_EntropyRange(t *testing.T) {
	// Entropy should always be non-negative.
	texts := []string{
		"Hello world. Goodbye moon.",
		"Single sentence here.",
		"One. Two. Three. Four. Five.",
	}
	cfg := config.DefaultConfig()
	for _, text := range texts {
		result := structure.SentenceOpenerDiversity(cfg, text)
		if result.Entropy < 0 {
			t.Errorf("Entropy should never be negative, got %v for %q", result.Entropy, text)
		}
	}
}
