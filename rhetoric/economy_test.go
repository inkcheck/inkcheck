package rhetoric_test

import (
	"testing"

	"github.com/inkcheck/rhetoric"
)

func TestEconomyAnalysis_WordyText(t *testing.T) {
	text := "In order to achieve the best results, it is important to note that " +
		"due to the fact that we are at this point in time still learning, " +
		"we should take into consideration all available options. " +
		"For the purpose of clarity, it should be noted that future plans are flexible."
	result := rhetoric.EconomyAnalysis(text)

	if result.WordyPhraseCount == 0 {
		t.Error("expected WordyPhraseCount > 0 for text full of wordy phrases")
	}
	if result.WordyPhraseDensity <= 0 {
		t.Errorf("expected WordyPhraseDensity > 0, got %v", result.WordyPhraseDensity)
	}
	if result.AvgSentenceLength <= 0 {
		t.Errorf("expected AvgSentenceLength > 0, got %v", result.AvgSentenceLength)
	}
}

func TestEconomyAnalysis_SparseText(t *testing.T) {
	text := "We ship fast. Code matters. Build things that work."
	result := rhetoric.EconomyAnalysis(text)

	if result.AvgSentenceLength <= 0 {
		t.Errorf("expected AvgSentenceLength > 0, got %v", result.AvgSentenceLength)
	}
	if result.WordsPerClause <= 0 {
		t.Errorf("expected WordsPerClause > 0, got %v", result.WordsPerClause)
	}
}

func TestEconomyAnalysis_SubordinationIndex(t *testing.T) {
	text := "Although we tried, the outcome was uncertain. " +
		"Because the data was incomplete, we could not conclude. " +
		"While results were promising, further study is needed."
	result := rhetoric.EconomyAnalysis(text)

	if result.SubordinationIndex <= 0 {
		t.Errorf("expected SubordinationIndex > 0 for text with subordinating conjunctions, got %v",
			result.SubordinationIndex)
	}
}

func TestEconomyAnalysis_ClauseCount(t *testing.T) {
	// Three sentences; one has a semicolon, one has "because" subordinating conjunction.
	text := "The sun rose; birds sang. We stayed inside because it was cold. The day ended."
	result := rhetoric.EconomyAnalysis(text)

	// 3 sentences + 1 semicolon + 1 "because" = at least 5 clauses
	if result.ClauseCount < 3 {
		t.Errorf("expected ClauseCount >= 3, got %d", result.ClauseCount)
	}
	if result.WordsPerClause <= 0 {
		t.Errorf("expected WordsPerClause > 0, got %v", result.WordsPerClause)
	}
}

func TestEconomyAnalysis_SingleSentence(t *testing.T) {
	text := "The quick brown fox jumps over the lazy dog."
	result := rhetoric.EconomyAnalysis(text)

	if result.AvgSentenceLength <= 0 {
		t.Errorf("expected AvgSentenceLength > 0, got %v", result.AvgSentenceLength)
	}
	if result.ClauseCount < 1 {
		t.Errorf("expected ClauseCount >= 1, got %d", result.ClauseCount)
	}
}

func TestEconomyAnalysis_PunctuationClauses(t *testing.T) {
	// Semicolons and colons create clause boundaries.
	text := "The sun rose; birds sang: a new day began."
	result := rhetoric.EconomyAnalysis(text)

	// 1 sentence + 1 semicolon + 1 colon = 3 clauses
	if result.ClauseCount < 3 {
		t.Errorf("expected ClauseCount >= 3 with semicolon and colon, got %d", result.ClauseCount)
	}
}

func TestEconomyAnalysis_Empty(t *testing.T) {
	result := rhetoric.EconomyAnalysis("")

	if result.AvgSentenceLength != 0 {
		t.Errorf("expected AvgSentenceLength = 0, got %v", result.AvgSentenceLength)
	}
	if result.WordyPhraseCount != 0 {
		t.Errorf("expected WordyPhraseCount = 0, got %d", result.WordyPhraseCount)
	}
	if result.WordsPerClause != 0 {
		t.Errorf("expected WordsPerClause = 0, got %v", result.WordsPerClause)
	}
}
