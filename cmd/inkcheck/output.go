package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/inkcheck/config"
	"github.com/inkcheck/reading"
	"github.com/inkcheck/rhetoric"
	"github.com/inkcheck/semantic"
	"github.com/inkcheck/structure"
)

// OutputFormat represents the output format type.
type OutputFormat string

const (
	// FormatText outputs metrics in human-readable text format.
	FormatText OutputFormat = "text"
	// FormatJSON outputs metrics in JSON format.
	FormatJSON OutputFormat = "json"
)

// MetricOutput holds the output for a single metric.
type MetricOutput struct {
	Name   string      `json:"name"`
	Value  interface{} `json:"value"`
	Source string      `json:"source,omitempty"`
}

// computeMetricValue computes and returns the value for a given metric.
func computeMetricValue(metric, text string, cfg config.Config, model *semantic.ModelManager) interface{} {
	switch metric {
	// Structure metrics
	case "paragraph_variance":
		return structure.ParagraphVariance(text)
	case "sentence_length_variance":
		return structure.SentenceLengthVariance(text)
	case "sentence_opener_diversity":
		return structure.SentenceOpenerDiversity(cfg, text)
	case "paragraph_position_analysis":
		return structure.ParagraphPositionAnalysis(cfg, text)
	case "punctuation_profile":
		return structure.PunctuationAnalysis(text)

	// Rhetoric metrics
	case "transition_word_density":
		return rhetoric.TransitionWordDensity(cfg, text)
	case "vocabulary_sophistication":
		return rhetoric.VocabSophisticationDistribution(cfg, text)
	case "hedging_analysis":
		return rhetoric.HedgingAnalysis(text)
	case "specificity_score":
		return rhetoric.SpecificityScore(text)
	case "voice_consistency":
		return rhetoric.VoiceConsistency(text)
	case "figurative_language":
		return rhetoric.FigurativeLanguagePresence(text)
	case "rhetorical_diversity":
		return rhetoric.RhetoricalDiversity(text)
	case "claim_support_ratio":
		return rhetoric.ClaimSupportRatio(text)
	case "counterargument_engagement":
		return rhetoric.CounterargumentEngagement(text)
	case "audience_awareness":
		return rhetoric.AudienceAwareness(cfg, text)
	case "argument_structure":
		return rhetoric.ArgumentStructureCoherence(text)
	case "tension_and_resolution":
		return rhetoric.TensionAndResolution(text)
	case "stance_analysis":
		return rhetoric.StanceAnalysis(text)
	case "contraction_rate":
		return rhetoric.ContractionRate(text)
	case "temporal_orientation":
		return rhetoric.TemporalOrientation(text)
	case "economy_analysis":
		return rhetoric.EconomyAnalysis(text)

	// Semantic metrics
	case "topic_coherence":
		return semantic.TopicCoherence(cfg, model, text)
	case "semantic_progression":
		return semantic.SemanticProgression(cfg, model, text)
	case "redundancy_detection":
		return semantic.RedundancyDetection(cfg, model, text)
	case "information_novelty":
		return semantic.InformationNoveltyCurve(cfg, model, text)
	case "emotional_tone":
		return semantic.EmotionalTone(cfg, model, text)

	// Readability metric
	case "readability":
		return reading.ReadabilityAnalysis(cfg, text)

	default:
		return nil
	}
}

// printMetricsJSON outputs metrics in JSON format.
func printMetricsJSON(w io.Writer, text, prefix, metric string, model *semantic.ModelManager, cfg config.Config) error {
	var metrics []MetricOutput

	if metric == "" {
		// Output all metrics
		for _, m := range allMetricNames {
			value := computeMetricValue(m, text, cfg, model)
			metrics = append(metrics, MetricOutput{
				Name:   m,
				Value:  value,
				Source: prefix,
			})
		}
	} else {
		// Output single metric
		value := computeMetricValue(metric, text, cfg, model)
		metrics = append(metrics, MetricOutput{
			Name:   metric,
			Value:  value,
			Source: prefix,
		})
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(metrics)
}

// printMetricsWithFormat outputs metrics in the specified format.
func printMetricsWithFormat(format OutputFormat, w io.Writer, text, prefix, metric string, model *semantic.ModelManager, cfg config.Config) error {
	switch format {
	case FormatJSON:
		return printMetricsJSON(w, text, prefix, metric, model, cfg)
	case FormatText:
		if metric == "" {
			for _, m := range allMetricNames {
				printMetric(w, text, prefix, m, model, cfg)
			}
		} else {
			printMetric(w, text, prefix, metric, model, cfg)
		}
		return nil
	default:
		return fmt.Errorf("unknown output format: %s", format)
	}
}
