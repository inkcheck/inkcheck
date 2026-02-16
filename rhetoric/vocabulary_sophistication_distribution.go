package rhetoric

import (
	"strings"

	"github.com/inkcheck/config"
	"github.com/inkcheck/shared"
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
	rank, found := wordFrequencyRank[lower]
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

	return VocabSophisticationResult{
		TotalWords:      totalWords,
		UniqueWords:     len(unique),
		TypeTokenRatio:  float64(len(unique)) / float64(totalWords),
		MATTR:           mattr,
		BandCounts:      bandCounts,
		BandRatios:      bandRatios,
		BandCV:          shared.CoefficientOfVariation(ratioValues),
		FormalWordCount: formalCount,
		FormalWordRatio: float64(formalCount) / float64(totalWords),
		FormalWords:     formalFound,
	}
}
