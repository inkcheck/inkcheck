package rhetoric

import (
	"strings"

	"github.com/inkcheck/shared"
)

// ArgumentStructureResult holds argument structure coherence analysis.
type ArgumentStructureResult struct {
	HasThesisMarker     bool
	HasEvidenceMarkers  bool
	HasConclusionMarker bool
	ThesisPosition      float64 // 0.0=start, 1.0=end
	ConclusionPosition  float64
	CoherenceScore      float64 // 0-1: higher means better structured
}

var thesisMarkers = []string{
	"i argue", "we argue", "this paper", "this essay", "the thesis",
	"the purpose", "the aim", "i believe", "we believe",
	"the central", "the main argument", "in this article",
}

var evidenceMarkers = []string{
	"for example", "for instance", "according to", "research shows",
	"studies indicate", "data suggests", "evidence shows",
	"as demonstrated", "the findings", "results show",
}

var conclusionMarkers = []string{
	"in conclusion", "to summarize", "in summary", "therefore",
	"thus", "overall", "to conclude", "in closing",
	"taken together", "all things considered",
}

// ArgumentStructureCoherence checks whether the text follows a standard
// thesis-evidence-conclusion structure by looking for markers at expected positions.
// TODO: LLM judge for semantic argument structure analysis
func ArgumentStructureCoherence(text string) ArgumentStructureResult {
	paragraphs := shared.SplitParagraphs(text)
	if len(paragraphs) < 3 {
		return ArgumentStructureResult{}
	}

	totalParas := float64(len(paragraphs))
	hasThesis := false
	hasEvidence := false
	hasConclusion := false
	thesisPos := 0.0
	conclusionPos := 0.0

	for i, p := range paragraphs {
		lower := strings.ToLower(p)
		pos := float64(i) / (totalParas - 1)

		if !hasThesis && shared.ContainsAny(lower, thesisMarkers) {
			hasThesis = true
			thesisPos = pos
		}
		if shared.ContainsAny(lower, evidenceMarkers) {
			hasEvidence = true
		}
		if shared.ContainsAny(lower, conclusionMarkers) {
			hasConclusion = true
			conclusionPos = pos
		}
	}

	// Coherence score: reward thesis early, conclusion late, evidence present
	score := 0.0
	if hasThesis {
		score += 0.33 * (1.0 - thesisPos) // higher if thesis is early
	}
	if hasEvidence {
		score += 0.34
	}
	if hasConclusion {
		score += 0.33 * conclusionPos // higher if conclusion is late
	}

	return ArgumentStructureResult{
		HasThesisMarker:     hasThesis,
		HasEvidenceMarkers:  hasEvidence,
		HasConclusionMarker: hasConclusion,
		ThesisPosition:      thesisPos,
		ConclusionPosition:  conclusionPos,
		CoherenceScore:      score,
	}
}
