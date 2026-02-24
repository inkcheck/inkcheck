package signature

// subMetricDef defines the normalization parameters for a single sub-metric.
type subMetricDef struct {
	Name   string
	Curve  Curve
	Lo, Hi float64
	Invert bool
	Weight float64
}

// axisDef defines a radar axis as a weighted composite of sub-metrics.
type axisDef struct {
	Name       string
	SubMetrics []subMetricDef
}

func equalWeight(n int) float64 {
	if n <= 0 {
		return 0
	}
	return 1.0 / float64(n)
}

// axesDefs returns the 10 axis definitions with sub-metric bounds, curves,
// and weights as specified in requirements.md. Default: equal weights.
func axesDefs() [10]axisDef {
	return [10]axisDef{
		formalityDef(),
		confidenceDef(),
		rhythmDef(),
		economyDef(),
		precisionDef(),
		coherenceDef(),
		vocabularyDef(),
		stanceDef(),
		emotionalToneDef(),
		temporalOrientationDef(),
	}
}

func formalityDef() axisDef {
	w := equalWeight(3)
	return axisDef{
		Name: "formality",
		SubMetrics: []subMetricDef{
			{Name: "formal_word_ratio", Curve: Linear, Lo: 0, Hi: 1, Invert: false, Weight: w},
			{Name: "passive_voice_ratio", Curve: Linear, Lo: 0, Hi: 1, Invert: false, Weight: w},
			{Name: "contraction_rate", Curve: Linear, Lo: 0, Hi: 0.15, Invert: true, Weight: w},
		},
	}
}

func confidenceDef() axisDef {
	w := equalWeight(3)
	return axisDef{
		Name: "confidence",
		SubMetrics: []subMetricDef{
			{Name: "hedging_density", Curve: Sqrt, Lo: 0, Hi: 5, Invert: true, Weight: w},
			{Name: "modal_verb_density", Curve: Sqrt, Lo: 0, Hi: 4, Invert: false, Weight: w},
			{Name: "active_voice_ratio", Curve: Linear, Lo: 0, Hi: 1, Invert: false, Weight: w},
		},
	}
}

func rhythmDef() axisDef {
	w := equalWeight(4)
	return axisDef{
		Name: "rhythm",
		SubMetrics: []subMetricDef{
			{Name: "sentence_length_cv", Curve: Linear, Lo: 0, Hi: 0.8, Invert: false, Weight: w},
			{Name: "paragraph_length_cv", Curve: Linear, Lo: 0, Hi: 1.5, Invert: false, Weight: w},
			{Name: "opener_diversity", Curve: Linear, Lo: 0, Hi: 2.5, Invert: false, Weight: w},
			{Name: "sentence_type_entropy", Curve: Linear, Lo: 0, Hi: 2.0, Invert: false, Weight: w},
		},
	}
}

func economyDef() axisDef {
	w := equalWeight(4)
	return axisDef{
		Name: "economy",
		SubMetrics: []subMetricDef{
			{Name: "avg_sentence_length", Curve: Log, Lo: 5, Hi: 40, Invert: true, Weight: w},
			{Name: "wordy_phrase_density", Curve: Sqrt, Lo: 0, Hi: 5, Invert: true, Weight: w},
			{Name: "words_per_clause", Curve: Linear, Lo: 1, Hi: 20, Invert: true, Weight: w},
			{Name: "syntactic_complexity", Curve: Sqrt, Lo: 0, Hi: 0.6, Invert: true, Weight: w},
		},
	}
}

func precisionDef() axisDef {
	w := equalWeight(3)
	return axisDef{
		Name: "precision",
		SubMetrics: []subMetricDef{
			{Name: "specificity_score", Curve: Linear, Lo: 0, Hi: 1, Invert: false, Weight: w},
			{Name: "vague_word_density", Curve: Sqrt, Lo: 0, Hi: 8, Invert: true, Weight: w},
			{Name: "redundancy_score", Curve: Linear, Lo: 0, Hi: 1, Invert: true, Weight: w},
		},
	}
}

func coherenceDef() axisDef {
	w := equalWeight(4)
	return axisDef{
		Name: "coherence",
		SubMetrics: []subMetricDef{
			{Name: "transition_density", Curve: Sqrt, Lo: 0, Hi: 6, Invert: false, Weight: w},
			{Name: "topic_coherence", Curve: Linear, Lo: 0, Hi: 1, Invert: false, Weight: w},
			{Name: "argument_structure", Curve: Linear, Lo: 0, Hi: 1, Invert: false, Weight: w},
			{Name: "claim_support_ratio", Curve: Linear, Lo: 0, Hi: 1, Invert: false, Weight: w},
		},
	}
}

func vocabularyDef() axisDef {
	w := equalWeight(3)
	return axisDef{
		Name: "vocabulary",
		SubMetrics: []subMetricDef{
			{Name: "mattr", Curve: Linear, Lo: 0, Hi: 1, Invert: false, Weight: w},
			{Name: "lexical_density", Curve: Linear, Lo: 0.3, Hi: 0.7, Invert: false, Weight: w},
			{Name: "low_freq_word_ratio", Curve: Linear, Lo: 0, Hi: 0.4, Invert: false, Weight: w},
		},
	}
}

func stanceDef() axisDef {
	return axisDef{
		Name: "stance",
		SubMetrics: []subMetricDef{
			{Name: "reader_centricity", Curve: Sqrt, Lo: 0, Hi: 0.05, Invert: false, Weight: 1.0},
		},
	}
}

func emotionalToneDef() axisDef {
	w := equalWeight(4)
	return axisDef{
		Name: "emotional_tone",
		SubMetrics: []subMetricDef{
			{Name: "positive_affect_ratio", Curve: Linear, Lo: 0, Hi: 0.15, Invert: false, Weight: w},
			{Name: "emotional_intensity", Curve: Sqrt, Lo: 0, Hi: 3, Invert: false, Weight: w},
			{Name: "empathy_marker_density", Curve: Sqrt, Lo: 0, Hi: 2, Invert: false, Weight: w},
			{Name: "arousal_level", Curve: Linear, Lo: 0, Hi: 1, Invert: false, Weight: w},
		},
	}
}

func temporalOrientationDef() axisDef {
	w := equalWeight(4)
	return axisDef{
		Name: "temporal_orientation",
		SubMetrics: []subMetricDef{
			{Name: "future_modal_density", Curve: Sqrt, Lo: 0, Hi: 4, Invert: false, Weight: w},
			{Name: "past_tense_ratio", Curve: Linear, Lo: 0, Hi: 0.6, Invert: true, Weight: w},
			{Name: "evidential_marker_density", Curve: Sqrt, Lo: 0, Hi: 3, Invert: true, Weight: w},
			{Name: "aspiration_marker_density", Curve: Sqrt, Lo: 0, Hi: 3, Invert: false, Weight: w},
		},
	}
}
