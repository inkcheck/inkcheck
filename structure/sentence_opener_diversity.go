package structure

import (
	"strings"

	"github.com/inkcheck/config"
	"github.com/inkcheck/shared"
)

// SentenceOpenerDiversityResult holds both the ratio and entropy of sentence openers.
type SentenceOpenerDiversityResult struct {
	Ratio   float64 // distinct openers / total sentences (0–1)
	Entropy float64 // Shannon entropy of opener distribution (bits)
}

// SentenceOpenerDiversity analyses the variety of sentence-opening patterns.
// It returns both a simple ratio (unique/total) and Shannon entropy of the
// opener frequency distribution for richer normalization options.
func SentenceOpenerDiversity(cfg config.Config, text string) SentenceOpenerDiversityResult {
	sentences := shared.SplitSentences(shared.ExtractProseText(text))
	if len(sentences) == 0 {
		return SentenceOpenerDiversityResult{}
	}

	freq := make(map[string]int)
	for _, s := range sentences {
		opener := extractOpener(s, cfg.OpenerWordCount)
		if opener != "" {
			freq[opener]++
		}
	}

	total := len(sentences)
	distinct := len(freq)

	ratio := float64(distinct) / float64(total)

	// Shannon entropy of the opener frequency distribution
	dist := make([]float64, 0, len(freq))
	for _, count := range freq {
		dist = append(dist, float64(count)/float64(total))
	}
	entropy := shared.Entropy(dist)

	return SentenceOpenerDiversityResult{
		Ratio:   ratio,
		Entropy: entropy,
	}
}

func extractOpener(text string, n int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}
	if len(words) < n {
		n = len(words)
	}
	parts := make([]string, n)
	for i := 0; i < n; i++ {
		parts[i] = strings.ToLower(words[i])
	}
	return strings.Join(parts, " ")
}
