package rhetoric_test

import (
	"testing"

	"github.com/inkcheck/rhetoric"
)

func TestSpecificityScore_VagueDensity_HighVagueness(t *testing.T) {
	// Dense vague language: important, significant, crucial, various, etc.
	text := "It is important to note that various significant factors are crucial. " +
		"Numerous fundamental aspects are effective. " +
		"Overall, the essential elements remain crucial and important. " +
		"Various significant impacts matter overall."
	result := rhetoric.SpecificityScore(text)

	if result.VagueDensity <= 0 {
		t.Errorf("expected VagueDensity > 0 for vague text, got %v", result.VagueDensity)
	}
}

func TestSpecificityScore_VagueDensity_Concrete(t *testing.T) {
	// Text with numbers, proper nouns, evidence phrases — low vague density.
	text := "In 2023, Apple reported $383 billion in revenue. " +
		"According to Nielsen, 72% of consumers prefer sustainable packaging. " +
		"Specifically, the study involved 1,200 participants across 15 countries. " +
		"For example, the Chicago plant reduced emissions by 34%."
	result := rhetoric.SpecificityScore(text)

	if result.VagueDensity > 2.0 {
		t.Errorf("expected low VagueDensity for concrete text, got %v", result.VagueDensity)
	}
}

func TestSpecificityScore_VagueDensity_NonNegative(t *testing.T) {
	texts := []string{
		"Hello world. This is a test sentence.",
		"",
		"Single.",
	}
	for _, text := range texts {
		result := rhetoric.SpecificityScore(text)
		if result.VagueDensity < 0 {
			t.Errorf("VagueDensity should never be negative, got %v for %q", result.VagueDensity, text)
		}
	}
}

func TestSpecificityScore_ScoresInRange(t *testing.T) {
	text := "The study found significant results. According to researchers, 80% of cases " +
		"showed improvement. Important factors include various environmental conditions."
	result := rhetoric.SpecificityScore(text)

	for i, s := range result.Scores {
		if s < 1 || s > 5 {
			t.Errorf("Scores[%d] = %d, want 1–5", i, s)
		}
	}
}
