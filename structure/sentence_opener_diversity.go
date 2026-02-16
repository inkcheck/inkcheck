package structure

import (
	"strings"

	"github.com/inkcheck/config"
	"github.com/inkcheck/shared"
)

func SentenceOpenerDiversity(cfg config.Config, text string) float64 {
	sentences := shared.SplitSentences(shared.ExtractProseText(text))
	if len(sentences) == 0 {
		return 0
	}
	seen := make(map[string]struct{})
	for _, s := range sentences {
		opener := extractOpener(s, cfg.OpenerWordCount)
		if opener != "" {
			seen[opener] = struct{}{}
		}
	}
	return float64(len(seen)) / float64(len(sentences))
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
