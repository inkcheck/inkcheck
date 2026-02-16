package rhetoric

import (
	"strings"

	"github.com/inkcheck/config"
	"github.com/inkcheck/shared"
)

var transitionPhrases = []string{
	"additionally", "also", "as a result", "besides", "consequently",
	"conversely", "finally", "for example", "for instance", "furthermore",
	"hence", "however", "in addition", "in contrast", "in fact", "indeed",
	"instead", "likewise", "meanwhile", "moreover", "nevertheless",
	"nonetheless", "on the other hand", "otherwise", "overall", "similarly",
	"specifically", "still", "subsequently", "therefore", "thus", "yet",
}

type TransitionResult struct {
	Total    int
	Distinct int
	Variety  float64
	Repeated []string
}

func TransitionWordDensity(cfg config.Config, text string) TransitionResult {
	lower := strings.ToLower(shared.ExtractProseText(text))
	counts := make(map[string]int)

	for _, phrase := range transitionPhrases {
		n := countOccurrences(lower, phrase)
		if n > 0 {
			counts[phrase] = n
		}
	}

	total := 0
	for _, c := range counts {
		total += c
	}

	var repeated []string
	for phrase, c := range counts {
		if c >= cfg.TransitionRepeatThreshold {
			repeated = append(repeated, phrase)
		}
	}

	variety := 0.0
	if total > 0 {
		variety = float64(len(counts)) / float64(total)
	}

	return TransitionResult{
		Total:    total,
		Distinct: len(counts),
		Variety:  variety,
		Repeated: repeated,
	}
}

