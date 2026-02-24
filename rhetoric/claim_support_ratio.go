package rhetoric

import (
	"strings"

	"github.com/inkcheck/shared"
)

// ClaimSupportResult holds claim vs evidence classification results.
type ClaimSupportResult struct {
	ClaimCount   int
	SupportCount int
	NeutralCount int
	Total        int
	Ratio        float64 // Support / Claim ratio; 0 if no claims
}

var claimSignals = []string{
	"should", "must", "need to", "it is essential", "it is important",
	"we believe", "i believe", "i argue", "i would argue", "we argue",
	"the key", "the main", "clearly", "obviously", "undoubtedly",
	"i think", "in my opinion", "in my view", "the most compelling",
	"what sets it apart", "the real question", "it is about",
}

var supportSignals = []string{
	"for example", "for instance", "according to", "research shows",
	"studies show", "data shows", "evidence suggests", "as shown",
	"demonstrated by", "specifically", "in particular", "such as",
	"the study", "the research", "percent", "%",
	"consider ", "take the case", "you can see", "you can read",
	"here are", "here is a", "a similar pattern",
	"describes in", "in his book", "in her book", "on wikipedia",
}

// ClaimSupportRatio classifies sentences as claims, support, or neutral
// based on keyword signals, and computes the support-to-claim ratio.
// TODO: LLM judge for context-aware claim/evidence classification
func ClaimSupportRatio(text string) ClaimSupportResult {
	prose := shared.ExtractProseText(text)
	sentences := shared.SplitSentences(prose)
	if len(sentences) == 0 {
		return ClaimSupportResult{}
	}

	claims, supports, neutral := 0, 0, 0
	for _, s := range sentences {
		lower := strings.ToLower(s)
		isClaim := shared.ContainsAny(lower, claimSignals)
		isSupport := shared.ContainsAny(lower, supportSignals)

		switch {
		case isClaim && !isSupport:
			claims++
		case isSupport && !isClaim:
			supports++
		case isClaim && isSupport:
			supports++ // evidence with claim language still counts as support
		default:
			neutral++
		}
	}

	// Ratio: support / claim, capped at 5.0 to keep it bounded.
	// Zero claims with some support → 1.0 (all evidence, no bare assertions).
	// Zero claims and zero support → 0.0.
	ratio := 0.0
	if claims > 0 {
		ratio = float64(supports) / float64(claims)
		if ratio > 5.0 {
			ratio = 5.0
		}
	} else if supports > 0 {
		ratio = 1.0
	}

	return ClaimSupportResult{
		ClaimCount:   claims,
		SupportCount: supports,
		NeutralCount: neutral,
		Total:        len(sentences),
		Ratio:        ratio,
	}
}
