package rhetoric

import "github.com/inkcheck/shared"

// VoiceConsistencyResult holds passive voice analysis results.
type VoiceConsistencyResult struct {
	// PassiveRatio is the overall ratio of passive sentences.
	PassiveRatio float64
	// ParagraphPassiveRatios is the passive ratio per paragraph.
	ParagraphPassiveRatios []float64
	// CV is the coefficient of variation of paragraph passive ratios.
	CV float64
}

// VoiceConsistency analyzes passive voice usage consistency across paragraphs.
// High CV suggests varied voice usage; low CV suggests uniform voice across paragraphs.
// TODO: LLM judge for nuanced passive voice detection
func VoiceConsistency(text string) VoiceConsistencyResult {
	paragraphs := shared.SplitParagraphs(text)
	if len(paragraphs) < 2 {
		return VoiceConsistencyResult{}
	}

	totalSentences := 0
	totalPassive := 0
	ratios := make([]float64, len(paragraphs))

	for i, p := range paragraphs {
		sentences := shared.SplitSentences(p)
		if len(sentences) == 0 {
			continue
		}
		passive := 0
		for _, s := range sentences {
			if isPassive(s) {
				passive++
			}
		}
		totalSentences += len(sentences)
		totalPassive += passive
		ratios[i] = float64(passive) / float64(len(sentences))
	}

	overallRatio := 0.0
	if totalSentences > 0 {
		overallRatio = float64(totalPassive) / float64(totalSentences)
	}

	return VoiceConsistencyResult{
		PassiveRatio:           overallRatio,
		ParagraphPassiveRatios: ratios,
		CV:                     shared.CoefficientOfVariation(ratios),
	}
}

// isPassive detects passive voice using POS tags.
// Looks for patterns: "be" verb (VB*) followed by past participle (VBN).
func isPassive(sentence string) bool {
	tokens := shared.Tokenize(sentence)
	for i := 0; i < len(tokens)-1; i++ {
		tag := tokens[i].Tag
		// Check for forms of "be": am, is, are, was, were, been, being
		if isBeVerb(tokens[i].Text) && (tag == "VBZ" || tag == "VBP" || tag == "VBD" || tag == "VBN" || tag == "VBG" || tag == "VB") {
			// Look for past participle after be verb (allow adverbs in between)
			for j := i + 1; j < len(tokens) && j <= i+3; j++ {
				if tokens[j].Tag == "VBN" {
					return true
				}
				if tokens[j].Tag != "RB" && tokens[j].Tag != "RBR" && tokens[j].Tag != "RBS" {
					break
				}
			}
		}
	}
	return false
}

func isBeVerb(word string) bool {
	switch word {
	case "am", "is", "are", "was", "were", "been", "being", "be":
		return true
	}
	return false
}
