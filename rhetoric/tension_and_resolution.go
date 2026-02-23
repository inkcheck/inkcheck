package rhetoric

import (
	"strings"

	"github.com/inkcheck/shared"
)

// TensionResolutionResult holds narrative tension arc analysis.
type TensionResolutionResult struct {
	TensionMarkers    int
	ResolutionMarkers int
	HasArc            bool    // true if tension precedes resolution
	ArcScore          float64 // 0-1: measures narrative arc quality
}

var tensionPhrases = []string{
	"but", "however", "yet", "although", "despite",
	"the problem", "the challenge", "the issue", "the question",
	"remains unclear", "is debated", "controversial",
	"on the other hand", "in contrast", "nevertheless",
	"complicates", "raises concerns", "critics",
}

var resolutionPhrases = []string{
	"the solution", "the answer", "this resolves",
	"as a result", "consequently", "therefore",
	"this shows", "this demonstrates", "this proves",
	"in conclusion", "ultimately", "finally",
	"the evidence supports", "this confirms",
}

// TensionAndResolution tracks tension and resolution marker positions
// to detect narrative arc structure.
// TODO: LLM judge for semantic tension/resolution analysis
func TensionAndResolution(text string) TensionResolutionResult {
	paragraphs := shared.SplitParagraphs(text)
	if len(paragraphs) < 2 {
		return TensionResolutionResult{}
	}

	tensionCount := 0
	resolutionCount := 0
	lastTensionPos := -1
	lastResolutionPos := -1

	for i, p := range paragraphs {
		lower := strings.ToLower(p)
		hasTension := shared.ContainsAny(lower, tensionPhrases)
		hasResolution := shared.ContainsAny(lower, resolutionPhrases)

		if hasTension {
			tensionCount++
			lastTensionPos = i
		}
		if hasResolution {
			resolutionCount++
			lastResolutionPos = i
		}
	}

	// Has arc if tension markers appear before resolution markers
	hasArc := tensionCount > 0 && resolutionCount > 0 && lastTensionPos < lastResolutionPos

	// Arc score: 0-1 based on presence of both and ordering
	arcScore := 0.0
	if tensionCount > 0 {
		arcScore += 0.3
	}
	if resolutionCount > 0 {
		arcScore += 0.3
	}
	if hasArc {
		arcScore += 0.4
	}

	return TensionResolutionResult{
		TensionMarkers:    tensionCount,
		ResolutionMarkers: resolutionCount,
		HasArc:            hasArc,
		ArcScore:          arcScore,
	}
}
