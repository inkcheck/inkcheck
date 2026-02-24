package rhetoric

import (
	"strings"

	"github.com/inkcheck/config"
	"github.com/inkcheck/shared"
	"github.com/inkcheck/wordlist"
)

// AudienceAwarenessResult holds audience engagement analysis.
type AudienceAwarenessResult struct {
	SecondPersonCount    int
	DirectQuestionCount  int
	ParentheticalCount   int
	JargonDensity        float64
	EngagementScore      float64
}

var secondPersonWords = []string{"you", "your", "yours", "yourself", "yourselves"}

// AudienceAwareness detects signals of audience engagement: second-person
// pronouns, direct questions, parenthetical asides, and jargon density.
// TODO: LLM judge for domain-specific jargon detection
func AudienceAwareness(cfg config.Config, text string) AudienceAwarenessResult {
	prose := shared.ExtractProseText(text)
	words := shared.ListWords(prose)
	if len(words) == 0 {
		return AudienceAwarenessResult{}
	}

	// Count second-person pronouns
	secondPerson := 0
	for _, w := range words {
		lower := strings.ToLower(shared.StripPunctuation(w))
		for _, sp := range secondPersonWords {
			if lower == sp {
				secondPerson++
				break
			}
		}
	}

	// Count direct questions
	sentences := shared.SplitSentences(prose)
	questions := 0
	for _, s := range sentences {
		if strings.HasSuffix(strings.TrimSpace(s), "?") {
			questions++
		}
	}

	// Count parenthetical asides
	parentheticals := strings.Count(prose, "(")

	// Jargon density: words outside configured rank threshold
	jargonCount := 0
	for _, w := range words {
		cleaned := shared.StripPunctuation(strings.ToLower(w))
		if cleaned == "" {
			continue
		}
		rank, found := wordlist.FrequencyRank()[cleaned]
		if !found || rank > cfg.JargonRankThreshold {
			jargonCount++
		}
	}
	jargonDensity := float64(jargonCount) / float64(len(words))

	// Engagement score: weighted combination
	score := float64(secondPerson)*0.3 + float64(questions)*0.3 +
		float64(parentheticals)*0.2 + (1.0-jargonDensity)*0.2
	// Normalize roughly to 0-1 range
	if totalSents := float64(len(sentences)); totalSents > 0 {
		score = score / totalSents
	}

	return AudienceAwarenessResult{
		SecondPersonCount:   secondPerson,
		DirectQuestionCount: questions,
		ParentheticalCount:  parentheticals,
		JargonDensity:       jargonDensity,
		EngagementScore:     score,
	}
}
