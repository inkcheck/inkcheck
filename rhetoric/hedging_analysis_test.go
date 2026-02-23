package rhetoric_test

import (
	"testing"

	"github.com/inkcheck/rhetoric"
)

func TestHedgingAnalysis_AssertiveModals_Confident(t *testing.T) {
	// Text heavy with assertive/confident modals.
	text := "We will deliver this project on time. You must complete the form before Friday. " +
		"The system shall comply with all regulations. Teams need to update their schedules. " +
		"All employees have to attend the briefing. The manager has to sign off on the plan. " +
		"We will achieve our goals. This approach shall yield results."
	result := rhetoric.HedgingAnalysis(text)

	if result.AssertiveModalCount == 0 {
		t.Error("expected AssertiveModalCount > 0 for confident-modal text")
	}
	if result.AssertiveModalDensity <= 0 {
		t.Errorf("expected AssertiveModalDensity > 0, got %v", result.AssertiveModalDensity)
	}
}

func TestHedgingAnalysis_AssertiveModals_Hedged(t *testing.T) {
	// Text with only uncertainty hedges, no assertive modals.
	text := "It may be possible that results could vary somewhat. " +
		"Perhaps the approach might work in some cases. " +
		"It seems this could be approximately correct. " +
		"Generally, this probably applies, though it is possible conditions differ."
	result := rhetoric.HedgingAnalysis(text)

	// Assertive modal count should be zero or very low.
	if result.AssertiveModalCount > 2 {
		t.Errorf("expected low AssertiveModalCount for hedged text, got %d", result.AssertiveModalCount)
	}
}

func TestHedgingAnalysis_AssertiveModalDensity_Proportional(t *testing.T) {
	text := "We will succeed. You must try harder. The plan shall work. " +
		"Everyone needs to cooperate. Teams have to coordinate well. " +
		"Leaders ought to listen carefully. Results will improve significantly."
	result := rhetoric.HedgingAnalysis(text)

	if result.AssertiveModalDensity <= 0 {
		t.Errorf("expected AssertiveModalDensity > 0, got %v", result.AssertiveModalDensity)
	}
	// Density should be assertiveCount / wordCount * 100
	if result.AssertiveModalCount > 0 && result.AssertiveModalDensity == 0 {
		t.Error("AssertiveModalDensity should be non-zero when AssertiveModalCount > 0")
	}
}

func TestHedgingAnalysis_BothHedgesAndAssertive(t *testing.T) {
	// Text with a mix — both hedges and assertive modals coexist.
	text := "We believe this may work, but we will commit to the plan regardless. " +
		"Results might vary; however, we must act now. " +
		"It is possible conditions change, yet you must stay the course. " +
		"Perhaps timing will shift, but the project shall launch in Q3."
	result := rhetoric.HedgingAnalysis(text)

	if result.Total == 0 {
		t.Error("expected hedges to be detected in mixed text")
	}
	if result.AssertiveModalCount == 0 {
		t.Error("expected assertive modals to be detected in mixed text")
	}
}

func TestHedgingAnalysis_Empty(t *testing.T) {
	result := rhetoric.HedgingAnalysis("")

	if result.AssertiveModalCount != 0 {
		t.Errorf("expected AssertiveModalCount = 0 for empty text, got %d", result.AssertiveModalCount)
	}
	if result.AssertiveModalDensity != 0 {
		t.Errorf("expected AssertiveModalDensity = 0 for empty text, got %v", result.AssertiveModalDensity)
	}
}
