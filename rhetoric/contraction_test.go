package rhetoric_test

import (
	"testing"

	"github.com/inkcheck/rhetoric"
)

func TestContractionRate_WithContractions(t *testing.T) {
	text := "We don't always get it right. I'm not sure that's the best approach. " +
		"It isn't easy, but we're working on it. You'll see the difference."
	result := rhetoric.ContractionRate(text)

	if result.Count == 0 {
		t.Error("expected Count > 0 for text with contractions")
	}
	if result.Rate <= 0 {
		t.Errorf("expected Rate > 0, got %v", result.Rate)
	}
	if result.Rate > 1.0 {
		t.Errorf("expected Rate <= 1.0, got %v", result.Rate)
	}
}

func TestContractionRate_Formal(t *testing.T) {
	text := "The analysis demonstrates that the results are consistent with prior findings. " +
		"The organisation has committed to reviewing the policy in the next quarter."
	result := rhetoric.ContractionRate(text)

	if result.Count != 0 {
		t.Errorf("expected Count = 0 for formal text, got %d", result.Count)
	}
	if result.Rate != 0.0 {
		t.Errorf("expected Rate = 0.0, got %v", result.Rate)
	}
}

func TestContractionRate_SingleContraction(t *testing.T) {
	result := rhetoric.ContractionRate("don't")

	if result.Count != 1 {
		t.Errorf("expected Count = 1, got %d", result.Count)
	}
	if result.Rate != 1.0 {
		t.Errorf("expected Rate = 1.0, got %v", result.Rate)
	}
}

func TestContractionRate_Empty(t *testing.T) {
	result := rhetoric.ContractionRate("")

	if result.Count != 0 {
		t.Errorf("expected Count = 0, got %d", result.Count)
	}
	if result.Rate != 0.0 {
		t.Errorf("expected Rate = 0.0, got %v", result.Rate)
	}
}

func TestContractionRate_MarkdownCodeBlock(t *testing.T) {
	// Contractions inside code blocks should be ignored by ExtractProseText.
	text := "This is formal text.\n\n```\ndon't use contractions in code\n```\n\nThe analysis continues."
	result := rhetoric.ContractionRate(text)

	if result.Count != 0 {
		t.Errorf("expected Count = 0 when contractions only appear in code blocks, got %d", result.Count)
	}
}
