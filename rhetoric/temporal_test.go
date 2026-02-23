package rhetoric_test

import (
	"testing"

	"github.com/inkcheck/rhetoric"
)

func TestTemporalOrientation_Prospective(t *testing.T) {
	text := "We will deliver a better experience next year. " +
		"Our goal is to expand into new markets. We aim to serve more customers. " +
		"We're committed to building something great. We plan to release in Q3."
	result := rhetoric.TemporalOrientation(text)

	if result.FutureModalCount == 0 {
		t.Error("expected FutureModalCount > 0 for prospective text")
	}
	if result.FutureModalDensity <= 0 {
		t.Errorf("expected FutureModalDensity > 0, got %v", result.FutureModalDensity)
	}
	if result.AspirationCount == 0 {
		t.Error("expected AspirationCount > 0 for goal/aim language")
	}
	if result.AspirationDensity <= 0 {
		t.Errorf("expected AspirationDensity > 0, got %v", result.AspirationDensity)
	}
}

func TestTemporalOrientation_Retrospective(t *testing.T) {
	text := "Studies showed that reading improves cognition. " +
		"According to the research, participants responded positively. " +
		"It was found that regular exercise reduced stress. " +
		"The team analysed the data and found consistent patterns."
	result := rhetoric.TemporalOrientation(text)

	if result.EvidentialCount == 0 {
		t.Error("expected EvidentialCount > 0 for evidential text")
	}
	if result.EvidentialDensity <= 0 {
		t.Errorf("expected EvidentialDensity > 0, got %v", result.EvidentialDensity)
	}
	if result.PastTenseCount == 0 {
		t.Error("expected PastTenseCount > 0 for past-tense text")
	}
}

func TestTemporalOrientation_RegularPastTense(t *testing.T) {
	// Tests the -ed suffix heuristic specifically.
	text := "The researchers investigated the phenomenon. They discovered patterns " +
		"and published their findings. The results confirmed earlier predictions."
	result := rhetoric.TemporalOrientation(text)

	// "investigated", "discovered", "published", "confirmed" should all be caught
	if result.PastTenseCount < 4 {
		t.Errorf("expected PastTenseCount >= 4 for text with regular -ed verbs, got %d", result.PastTenseCount)
	}
}

func TestTemporalOrientation_NoDoubleCount(t *testing.T) {
	// "We will be ready" has one future modal token ("will").
	// It should NOT be double-counted by phrase matching.
	text := "We will be ready."
	result := rhetoric.TemporalOrientation(text)

	if result.FutureModalCount != 1 {
		t.Errorf("expected FutureModalCount = 1 (only 'will' token), got %d", result.FutureModalCount)
	}
}

func TestTemporalOrientation_GoingToPhrase(t *testing.T) {
	text := "We are going to succeed. They are going to try."
	result := rhetoric.TemporalOrientation(text)

	// "going to" appears twice; "are going to" also matches twice.
	// Both should contribute to future count.
	if result.FutureModalCount < 2 {
		t.Errorf("expected FutureModalCount >= 2, got %d", result.FutureModalCount)
	}
}

func TestTemporalOrientation_FalsePositiveEd(t *testing.T) {
	// Words in the exclusion list should not be counted as past tense.
	text := "The sacred naked wicked hundred kindred seed need feed breed speed bleed."
	result := rhetoric.TemporalOrientation(text)

	if result.PastTenseCount != 0 {
		t.Errorf("expected PastTenseCount = 0 for false-positive -ed words, got %d", result.PastTenseCount)
	}
}

func TestTemporalOrientation_Empty(t *testing.T) {
	result := rhetoric.TemporalOrientation("")

	if result.FutureModalCount != 0 || result.PastTenseCount != 0 ||
		result.EvidentialCount != 0 || result.AspirationCount != 0 {
		t.Error("expected all counts to be 0 for empty text")
	}
	if result.FutureModalDensity != 0 || result.PastTenseDensity != 0 ||
		result.EvidentialDensity != 0 || result.AspirationDensity != 0 {
		t.Error("expected all densities to be 0 for empty text")
	}
}
