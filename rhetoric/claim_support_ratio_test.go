package rhetoric_test

import (
	"testing"

	"github.com/inkcheck/rhetoric"
)

func TestClaimSupportRatio_BalancedText(t *testing.T) {
	// Mix of claims and supporting evidence.
	text := "We believe innovation is essential for growth. " +
		"For example, companies that invested in R&D grew 30% faster according to research. " +
		"It is important to act quickly. " +
		"Studies show that early movers capture 60% of market share. " +
		"We argue that the evidence is clear."
	result := rhetoric.ClaimSupportRatio(text)

	if result.ClaimCount == 0 {
		t.Error("expected ClaimCount > 0")
	}
	if result.SupportCount == 0 {
		t.Error("expected SupportCount > 0")
	}
	if result.Ratio < 0 {
		t.Errorf("Ratio should be non-negative, got %v", result.Ratio)
	}
	if result.Ratio > 5.0 {
		t.Errorf("Ratio should be capped at 5.0, got %v", result.Ratio)
	}
}

func TestClaimSupportRatio_ZeroClaims_WithSupport(t *testing.T) {
	// Only evidence sentences, no bare assertions.
	text := "For example, the data shows a 20% increase. " +
		"According to the report, revenue grew by $10 million. " +
		"Research shows a clear correlation. " +
		"As shown in the graph, the trend is upward. " +
		"Specifically, three studies demonstrated this effect."
	result := rhetoric.ClaimSupportRatio(text)

	if result.ClaimCount != 0 {
		t.Errorf("expected ClaimCount = 0, got %d", result.ClaimCount)
	}
	if result.SupportCount == 0 {
		t.Error("expected SupportCount > 0")
	}
	if result.Ratio != 1.0 {
		t.Errorf("expected Ratio = 1.0 for zero-claim-with-support, got %v", result.Ratio)
	}
}

func TestClaimSupportRatio_ZeroClaimsZeroSupport(t *testing.T) {
	// Neutral sentences with no claim or support signals.
	text := "The cat walked down the street. A bird landed on the fence. " +
		"The sky was blue. Time passed slowly."
	result := rhetoric.ClaimSupportRatio(text)

	if result.Ratio != 0.0 {
		t.Errorf("expected Ratio = 0.0 for no claims and no support, got %v", result.Ratio)
	}
}

func TestClaimSupportRatio_RatioCap(t *testing.T) {
	// Massive evidence relative to claims — ratio should be capped at 5.0.
	text := "I believe this is true. " +
		"For example, study 1 shows X. For example, study 2 shows Y. " +
		"Research shows Z. Data shows W. Evidence suggests V. " +
		"As shown in report A. Specifically, finding B. " +
		"According to source C. The research confirms D. " +
		"The study found E. Percent of cases show F."
	result := rhetoric.ClaimSupportRatio(text)

	if result.Ratio > 5.0 {
		t.Errorf("Ratio should be capped at 5.0, got %v", result.Ratio)
	}
}

func TestClaimSupportRatio_Empty(t *testing.T) {
	result := rhetoric.ClaimSupportRatio("")

	if result.Total != 0 {
		t.Errorf("expected Total = 0 for empty text, got %d", result.Total)
	}
	if result.Ratio != 0 {
		t.Errorf("expected Ratio = 0 for empty text, got %v", result.Ratio)
	}
}
