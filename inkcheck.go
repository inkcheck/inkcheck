package inkcheck

import (
	"github.com/inkcheck/config"
	"github.com/inkcheck/reading"
	"github.com/inkcheck/rhetoric"
	"github.com/inkcheck/semantic"
	"github.com/inkcheck/structure"
)

// Result holds the output of all inkcheck metrics for a text.
type Result struct {
	Structure   StructureResult
	Rhetoric    RhetoricResult
	Semantic    SemanticResult
	Readability ReadabilityResult
}

// StructureResult holds all structural metric outputs.
type StructureResult struct {
	ParagraphVariance       float64
	ParagraphLengths        []int
	SentenceLengthVariance  float64
	SentenceOpenerDiversity float64
	ParagraphPosition       structure.ParagraphPositionResult
	Punctuation             structure.PunctuationProfile
}

// RhetoricResult holds all rhetorical metric outputs.
type RhetoricResult struct {
	TransitionWordDensity   rhetoric.TransitionResult
	VocabSophistication     rhetoric.VocabSophisticationResult
	Hedging                 rhetoric.HedgingResult
	Specificity             rhetoric.SpecificityResult
	VoiceConsistency        rhetoric.VoiceConsistencyResult
	FigurativeLanguage      rhetoric.FigurativeLanguageResult
	RhetoricalDiversity     rhetoric.RhetoricalDiversityResult
	ClaimSupport            rhetoric.ClaimSupportResult
	Counterargument         rhetoric.CounterargumentResult
	AudienceAwareness       rhetoric.AudienceAwarenessResult
	ArgumentStructure       rhetoric.ArgumentStructureResult
	TensionAndResolution    rhetoric.TensionResolutionResult
}

// SemanticResult holds all semantic metric outputs.
type SemanticResult struct {
	TopicCoherence      semantic.TopicCoherenceResult
	SemanticProgression semantic.SemanticProgressionResult
	Redundancy          semantic.RedundancyResult
	InformationNovelty  semantic.InformationNoveltyCurveResult
}

// ReadabilityResult holds readability metric output.
type ReadabilityResult struct {
	Formula string
	Score   float64
	Grade   string
}

// Analyze runs all structure, rhetoric, and readability metrics on the given text.
// Semantic metrics are skipped (use AnalyzeAll with a ModelManager instead).
func Analyze(cfg config.Config, text string) Result {
	return Result{
		Structure:   AnalyzeStructure(cfg, text),
		Rhetoric:    AnalyzeRhetoric(cfg, text),
		Readability: AnalyzeReadability(cfg, text),
	}
}

// AnalyzeAll runs all metrics including semantic analysis.
func AnalyzeAll(cfg config.Config, text string, model *semantic.ModelManager) Result {
	return Result{
		Structure:   AnalyzeStructure(cfg, text),
		Rhetoric:    AnalyzeRhetoric(cfg, text),
		Semantic:    AnalyzeSemantic(cfg, text, model),
		Readability: AnalyzeReadability(cfg, text),
	}
}

// AnalyzeStructure runs only the 5 structural metrics.
func AnalyzeStructure(cfg config.Config, text string) StructureResult {
	return StructureResult{
		ParagraphVariance:       structure.ParagraphVariance(text),
		ParagraphLengths:        structure.ParagraphLengths(text),
		SentenceLengthVariance:  structure.SentenceLengthVariance(text),
		SentenceOpenerDiversity: structure.SentenceOpenerDiversity(cfg, text),
		ParagraphPosition:       structure.ParagraphPositionAnalysis(cfg, text),
		Punctuation:             structure.PunctuationAnalysis(text),
	}
}

// AnalyzeRhetoric runs only the 12 rhetorical metrics.
func AnalyzeRhetoric(cfg config.Config, text string) RhetoricResult {
	return RhetoricResult{
		TransitionWordDensity:   rhetoric.TransitionWordDensity(cfg, text),
		VocabSophistication:     rhetoric.VocabSophisticationDistribution(cfg, text),
		Hedging:                 rhetoric.HedgingAnalysis(text),
		Specificity:             rhetoric.SpecificityScore(text),
		VoiceConsistency:        rhetoric.VoiceConsistency(text),
		FigurativeLanguage:      rhetoric.FigurativeLanguagePresence(text),
		RhetoricalDiversity:     rhetoric.RhetoricalDiversity(text),
		ClaimSupport:            rhetoric.ClaimSupportRatio(text),
		Counterargument:         rhetoric.CounterargumentEngagement(text),
		AudienceAwareness:       rhetoric.AudienceAwareness(cfg, text),
		ArgumentStructure:       rhetoric.ArgumentStructureCoherence(text),
		TensionAndResolution:    rhetoric.TensionAndResolution(text),
	}
}

// AnalyzeSemantic runs only the 4 semantic metrics.
func AnalyzeSemantic(cfg config.Config, text string, model *semantic.ModelManager) SemanticResult {
	return SemanticResult{
		TopicCoherence:      semantic.TopicCoherence(cfg, model, text),
		SemanticProgression: semantic.SemanticProgression(cfg, model, text),
		Redundancy:          semantic.RedundancyDetection(cfg, model, text),
		InformationNovelty:  semantic.InformationNoveltyCurve(cfg, model, text),
	}
}

// AnalyzeReadability runs the readability metric.
func AnalyzeReadability(cfg config.Config, text string) ReadabilityResult {
	r := reading.ReadabilityAnalysis(cfg, text)
	return ReadabilityResult{
		Formula: r.Formula,
		Score:   r.Score,
		Grade:   r.Grade,
	}
}
