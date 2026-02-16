package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/inkcheck/config"
	"github.com/inkcheck/reading"
	"github.com/inkcheck/rhetoric"
	"github.com/inkcheck/semantic"
	"github.com/inkcheck/structure"
)

type printerFunc func(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager)

// metricPrinters maps metric names to their printer functions.
var metricPrinters = map[string]printerFunc{
	"paragraph_variance":          printParagraphVariance,
	"sentence_length_variance":    printSentenceLengthVariance,
	"sentence_opener_diversity":   printSentenceOpenerDiversity,
	"paragraph_position_analysis": printParagraphPosition,
	"punctuation_profile":         printPunctuationProfile,
	"transition_word_density":     printTransitionWordDensity,
	"vocabulary_sophistication":   printVocabSophistication,
	"hedging_analysis":            printHedgingAnalysis,
	"specificity_score":           printSpecificityScore,
	"voice_consistency":           printVoiceConsistency,
	"figurative_language":         printFigurativeLanguage,
	"rhetorical_diversity":        printRhetoricalDiversity,
	"claim_support_ratio":         printClaimSupportRatio,
	"counterargument_engagement":  printCounterargument,
	"audience_awareness":          printAudienceAwareness,
	"argument_structure":          printArgumentStructure,
	"tension_and_resolution":      printTensionAndResolution,
	"topic_coherence":             printTopicCoherence,
	"semantic_progression":        printSemanticProgression,
	"redundancy_detection":        printRedundancyDetection,
	"information_novelty":         printInformationNovelty,
	"readability":                 printReadability,
}

// printMetric prints a single metric using the registered printer function.
func printMetric(w io.Writer, text, prefix, metric string, model *semantic.ModelManager, cfg config.Config) {
	label := metric
	if prefix != "" {
		label = prefix + "\t" + metric
	}

	printer, ok := metricPrinters[metric]
	if !ok {
		return
	}
	printer(w, label, text, cfg, model)
}

// Individual metric printers

func printParagraphVariance(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	fmt.Fprintf(w, "%s\t%.4f\n", label, structure.ParagraphVariance(text))
}

func printSentenceLengthVariance(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	fmt.Fprintf(w, "%s\t%.4f\n", label, structure.SentenceLengthVariance(text))
}

func printSentenceOpenerDiversity(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	fmt.Fprintf(w, "%s\t%.4f\n", label, structure.SentenceOpenerDiversity(cfg, text))
}

func printParagraphPosition(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	r := structure.ParagraphPositionAnalysis(cfg, text)
	uniform := "no"
	if r.Uniform {
		uniform = "yes"
	}
	fmt.Fprintf(w, "%s\tuniform=%s\t(open=%d close=%d body_mean=%.1f open_dev=%.2f close_dev=%.2f)\n",
		label, uniform,
		r.OpeningLength, r.ClosingLength, r.BodyMeanLength,
		r.OpeningDeviation, r.ClosingDeviation)
}

func printPunctuationProfile(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	p := structure.PunctuationAnalysis(text)
	fmt.Fprintf(w, "%s\t%d/%d types\t(. %d , %d ; %d : %d — %d () %d ? %d ! %d)\n",
		label, p.Variety(), 8,
		p.Periods, p.Commas, p.Semicolons, p.Colons,
		p.Dashes, p.Parentheses, p.Questions, p.Exclamations)
}

func printTransitionWordDensity(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	r := rhetoric.TransitionWordDensity(cfg, text)
	fmt.Fprintf(w, "%s\t%.4f\t(total=%d distinct=%d", label, r.Variety, r.Total, r.Distinct)
	if len(r.Repeated) > 0 {
		fmt.Fprintf(w, " repeated=%s", strings.Join(r.Repeated, ","))
	}
	fmt.Fprintf(w, ")\n")
}

func printVocabSophistication(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	r := rhetoric.VocabSophisticationDistribution(cfg, text)
	fmt.Fprintf(w, "%s\tttr=%.4f\tmattr=%.4f\tband_cv=%.4f\tformal_words=%d (%.4f)", label, r.TypeTokenRatio, r.MATTR, r.BandCV, r.FormalWordCount, r.FormalWordRatio)
	if len(r.FormalWords) > 0 {
		words := make([]string, 0, len(r.FormalWords))
		for word, c := range r.FormalWords {
			words = append(words, fmt.Sprintf("%s:%d", word, c))
		}
		fmt.Fprintf(w, "\t(%s)", strings.Join(words, " "))
	}
	fmt.Fprintf(w, "\n")
}

func printHedgingAnalysis(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	r := rhetoric.HedgingAnalysis(text)
	fmt.Fprintf(w, "%s\tdensity=%.2f/100w\t(total=%d distinct=%d variety=%.4f modal=%d approx=%d plausibility=%d attribution=%d frequency=%d)\n",
		label, r.Density, r.Total, r.Distinct, r.Variety,
		r.Categories.Modal, r.Categories.Approximator, r.Categories.Plausibility,
		r.Categories.Attribution, r.Categories.Frequency)
}

func printSpecificityScore(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	r := rhetoric.SpecificityScore(text)
	fmt.Fprintf(w, "%s\tmean=%.2f\trange=%.0f\tcv=%.4f\n", label, r.Mean, r.Range, r.CV)
}

func printVoiceConsistency(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	r := rhetoric.VoiceConsistency(text)
	fmt.Fprintf(w, "%s\tpassive_ratio=%.4f\tcv=%.4f\n", label, r.PassiveRatio, r.CV)
}

func printFigurativeLanguage(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	r := rhetoric.FigurativeLanguagePresence(text)
	fmt.Fprintf(w, "%s\t%d instances\t(similes=%d rhetorical_q=%d alliteration=%d density=%.2f/100w)\n",
		label, r.TotalInstances, r.SimileCount, r.RhetoricalQuestionCount, r.AlliterationCount, r.DensityPer100Words)
}

func printRhetoricalDiversity(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	r := rhetoric.RhetoricalDiversity(text)
	fmt.Fprintf(w, "%s\tentropy=%.4f\t(questions=%d exclamations=%d imperatives=%d conditionals=%d declaratives=%d)\n",
		label, r.Entropy, r.Questions, r.Exclamations, r.Imperatives, r.Conditionals, r.Declaratives)
}

func printClaimSupportRatio(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	r := rhetoric.ClaimSupportRatio(text)
	fmt.Fprintf(w, "%s\tratio=%.4f\t(claims=%d support=%d neutral=%d)\n",
		label, r.Ratio, r.ClaimCount, r.SupportCount, r.NeutralCount)
}

func printCounterargument(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	r := rhetoric.CounterargumentEngagement(text)
	fmt.Fprintf(w, "%s\t%d instances\t(density=%.2f/100w", label, r.Instances, r.DensityPer100)
	if len(r.Phrases) > 0 {
		fmt.Fprintf(w, " phrases=%s", strings.Join(r.Phrases, ","))
	}
	fmt.Fprintf(w, ")\n")
}

func printAudienceAwareness(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	r := rhetoric.AudienceAwareness(cfg, text)
	fmt.Fprintf(w, "%s\tscore=%.4f\t(2nd_person=%d questions=%d parentheticals=%d jargon=%.4f)\n",
		label, r.EngagementScore, r.SecondPersonCount, r.DirectQuestionCount, r.ParentheticalCount, r.JargonDensity)
}

func printArgumentStructure(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	r := rhetoric.ArgumentStructureCoherence(text)
	fmt.Fprintf(w, "%s\tcoherence=%.4f\t(thesis=%v evidence=%v conclusion=%v thesis_pos=%.2f conclusion_pos=%.2f)\n",
		label, r.CoherenceScore, r.HasThesisMarker, r.HasEvidenceMarkers, r.HasConclusionMarker, r.ThesisPosition, r.ConclusionPosition)
}

func printTensionAndResolution(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	r := rhetoric.TensionAndResolution(text)
	arc := "no"
	if r.HasArc {
		arc = "yes"
	}
	fmt.Fprintf(w, "%s\tarc=%s\tscore=%.4f\t(tension=%d resolution=%d)\n",
		label, arc, r.ArcScore, r.TensionMarkers, r.ResolutionMarkers)
}

func printTopicCoherence(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	r := semantic.TopicCoherence(cfg, model, text)
	fmt.Fprintf(w, "%s\tmean=%.4f\tcv=%.4f\n", label, r.MeanSimilarity, r.CV)
}

func printSemanticProgression(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	r := semantic.SemanticProgression(cfg, model, text)
	fmt.Fprintf(w, "%s\tmean_drift=%.4f\tcv=%.4f\n", label, r.MeanDrift, r.CV)
}

func printRedundancyDetection(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	r := semantic.RedundancyDetection(cfg, model, text)
	fmt.Fprintf(w, "%s\t%d pairs", label, r.PairCount)
	for _, p := range r.Pairs {
		fmt.Fprintf(w, "\t(%d,%d)=%.4f", p.IndexA, p.IndexB, p.Similarity)
	}
	fmt.Fprintf(w, "\n")
}

func printInformationNovelty(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	r := semantic.InformationNoveltyCurve(cfg, model, text)
	fmt.Fprintf(w, "%s\tmean=%.4f\tcv=%.4f\n", label, r.MeanNovelty, r.CV)
}

func printReadability(w io.Writer, label, text string, cfg config.Config, model *semantic.ModelManager) {
	r := reading.ReadabilityAnalysis(cfg, text)
	grade := r.Grade
	if grade != "" {
		fmt.Fprintf(w, "%s\t%s\t%.2f\t(%s)\n", label, r.Formula, r.Score, grade)
	} else {
		fmt.Fprintf(w, "%s\t%s\t%.2f\n", label, r.Formula, r.Score)
	}
}
