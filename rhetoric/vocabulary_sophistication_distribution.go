package rhetoric

import (
	"strings"

	"github.com/inkcheck/config"
	"github.com/inkcheck/shared"
	"github.com/inkcheck/wordlist"
)

type VocabSophisticationResult struct {
	TotalWords      int
	UniqueWords     int
	TypeTokenRatio  float64
	MATTR           float64
	BandCounts      [5]int
	BandRatios      [5]float64
	BandCV          float64
	FormalWordCount int
	FormalWordRatio float64
	FormalWords     map[string]int
	LexicalDensity  float64 // content words / total words
	LowFreqWordRatio float64 // BandRatios[2] + [3] + [4] (6k+ rank or unknown)
}

// functionWords is a closed-class set used to identify content words.
// Content words = total words − function words; drives LexicalDensity.
var functionWords = map[string]bool{
	// articles
	"a": true, "an": true, "the": true,
	// prepositions
	"of": true, "in": true, "to": true, "for": true, "on": true, "with": true,
	"at": true, "by": true, "from": true, "as": true, "into": true, "through": true,
	"during": true, "before": true, "after": true, "above": true, "below": true,
	"between": true, "out": true, "off": true, "over": true, "under": true,
	"again": true, "further": true, "then": true, "once": true, "about": true,
	"against": true, "along": true, "among": true, "around": true, "upon": true,
	"without": true, "within": true, "throughout": true, "toward": true, "towards": true,
	"onto": true, "except": true, "per": true, "via": true, "versus": true,
	// conjunctions
	"and": true, "but": true, "or": true, "nor": true, "so": true, "yet": true,
	"if": true, "although": true, "because": true, "since": true, "while": true,
	"whereas": true, "whether": true, "unless": true, "until": true, "when": true,
	"where": true, "that": true, "which": true, "who": true, "whom": true,
	"though": true, "even": true, "both": true, "either": true, "neither": true,
	// auxiliary verbs
	"is": true, "are": true, "was": true, "were": true, "be": true, "been": true,
	"being": true, "have": true, "has": true, "had": true, "do": true, "does": true,
	"did": true, "will": true, "would": true, "shall": true, "should": true,
	"may": true, "might": true, "must": true, "can": true, "could": true,
	// pronouns
	"i": true, "me": true, "my": true, "myself": true,
	"we": true, "us": true, "our": true, "ours": true, "ourselves": true,
	"you": true, "your": true, "yours": true, "yourself": true, "yourselves": true,
	"he": true, "him": true, "his": true, "himself": true,
	"she": true, "her": true, "hers": true, "herself": true,
	"it": true, "its": true, "itself": true,
	"they": true, "them": true, "their": true, "theirs": true, "themselves": true,
	"what": true, "this": true, "these": true, "those": true, "there": true,
	// common adverbs
	"not": true, "no": true, "very": true, "just": true, "more": true, "also": true,
	"up": true, "down": true, "here": true, "how": true, "all": true, "some": true,
	"any": true, "each": true, "few": true, "most": true, "other": true,
	"such": true, "only": true, "same": true, "than": true, "too": true, "now": true,
}

var formalWords = map[string]bool{
	"utilize": true, "utilizes": true, "utilizing": true, "utilized": true,
	"leverage": true, "leverages": true, "leveraging": true, "leveraged": true,
	"nuanced": true, "multifaceted": true, "comprehensive": true,
	"furthermore": true, "moreover": true, "additionally": true,
	"facilitate": true, "facilitates": true, "facilitating": true,
	"paramount": true, "endeavor": true, "endeavors": true,
	"encompasses": true, "encompassing": true,
	"delve": true, "delves": true, "delving": true,
	"intricate": true, "intricacies": true, "pivotal": true,
	"robust": true, "streamline": true, "streamlines": true, "streamlining": true,
	"holistic": true, "synergy": true, "paradigm": true,
	"landscape": true, "ecosystem": true, "tapestry": true,
	"underscore": true, "underscores": true, "underscoring": true,
	"navigating": true, "navigate": true, "realm": true,
	"foster": true, "fosters": true, "fostering": true,
	"embark": true, "embarks": true, "embarking": true,
	"testament": true, "cornerstone": true,
	"groundbreaking": true, "cutting-edge": true,
}

func classifyWord(word string) int {
	lower := strings.ToLower(word)
	lower = shared.StripPunctuation(lower)
	if lower == "" {
		return 0
	}
	rank, found := wordlist.FrequencyRank()[lower]
	if !found {
		return 4
	}
	switch {
	case rank <= 1000:
		return 0
	case rank <= 3000:
		return 1
	case rank <= 6000:
		return 2
	default:
		return 3
	}
}

func computeMATTR(words []string, windowSize int) float64 {
	if len(words) < windowSize {
		if len(words) == 0 {
			return 0
		}
		unique := make(map[string]bool, len(words))
		for _, w := range words {
			unique[w] = true
		}
		return float64(len(unique)) / float64(len(words))
	}

	freq := make(map[string]int, windowSize)
	for i := 0; i < windowSize; i++ {
		freq[words[i]]++
	}

	sum := float64(len(freq)) / float64(windowSize)
	windows := 1

	for i := windowSize; i < len(words); i++ {
		leaving := words[i-windowSize]
		freq[leaving]--
		if freq[leaving] == 0 {
			delete(freq, leaving)
		}
		freq[words[i]]++
		sum += float64(len(freq)) / float64(windowSize)
		windows++
	}

	return sum / float64(windows)
}

func VocabSophisticationDistribution(cfg config.Config, text string) VocabSophisticationResult {
	prose := shared.ExtractProseText(text)
	words := shared.ListWords(prose)
	if len(words) == 0 {
		return VocabSophisticationResult{}
	}

	unique := make(map[string]bool)
	var bandCounts [5]int
	formalFound := make(map[string]int)
	cleanedWords := make([]string, 0, len(words))
	contentWordCount := 0

	for _, w := range words {
		lower := strings.ToLower(w)
		cleaned := shared.StripPunctuation(lower)
		if cleaned == "" {
			continue
		}
		unique[cleaned] = true
		cleanedWords = append(cleanedWords, cleaned)
		band := classifyWord(w)
		bandCounts[band]++
		if formalWords[cleaned] {
			formalFound[cleaned]++
		}
		if !functionWords[cleaned] {
			contentWordCount++
		}
	}

	totalWords := 0
	for _, c := range bandCounts {
		totalWords += c
	}
	if totalWords == 0 {
		return VocabSophisticationResult{}
	}

	var bandRatios [5]float64
	ratioValues := make([]float64, 5)
	for i, c := range bandCounts {
		bandRatios[i] = float64(c) / float64(totalWords)
		ratioValues[i] = bandRatios[i]
	}

	formalCount := 0
	for _, c := range formalFound {
		formalCount += c
	}

	mattr := computeMATTR(cleanedWords, cfg.MATTRWindowSize)

	lexicalDensity := float64(contentWordCount) / float64(totalWords)
	lowFreqWordRatio := bandRatios[2] + bandRatios[3] + bandRatios[4]

	return VocabSophisticationResult{
		TotalWords:       totalWords,
		UniqueWords:      len(unique),
		TypeTokenRatio:   float64(len(unique)) / float64(totalWords),
		MATTR:            mattr,
		BandCounts:       bandCounts,
		BandRatios:       bandRatios,
		BandCV:           shared.CoefficientOfVariation(ratioValues),
		FormalWordCount:  formalCount,
		FormalWordRatio:  float64(formalCount) / float64(totalWords),
		FormalWords:      formalFound,
		LexicalDensity:   lexicalDensity,
		LowFreqWordRatio: lowFreqWordRatio,
	}
}
