package signature

// Output types matching signature.schema.json for JSON serialization.

// SignatureOutput is the top-level JSON output matching InkwellSignature schema.
type SignatureOutput struct {
	Version    string              `json:"version"`
	Document   DocumentOutput      `json:"document"`
	Signature  SignatureAxesOutput `json:"signature"`
	Comparison *ComparisonOutput   `json:"comparison,omitempty"`
}

// DocumentOutput holds source document metadata and contextual annotations.
type DocumentOutput struct {
	WordCount          int                `json:"word_count"`
	SentenceCount      int                `json:"sentence_count"`
	ParagraphCount     int                `json:"paragraph_count"`
	ReadingTimeSeconds int                `json:"reading_time_seconds,omitempty"`
	Readability        *ReadabilityOutput `json:"readability,omitempty"`
}

// ReadabilityOutput holds contextual readability annotation.
type ReadabilityOutput struct {
	Formula string  `json:"formula"`
	Score   float64 `json:"score"`
	Grade   string  `json:"grade"`
}

// SignatureAxesOutput holds all 10 axis results.
type SignatureAxesOutput struct {
	Formality           AxisOutput `json:"formality"`
	Confidence          AxisOutput `json:"confidence"`
	Rhythm              AxisOutput `json:"rhythm"`
	Economy             AxisOutput `json:"economy"`
	Precision           AxisOutput `json:"precision"`
	Coherence           AxisOutput `json:"coherence"`
	Vocabulary          AxisOutput `json:"vocabulary"`
	Stance              AxisOutput `json:"stance"`
	EmotionalTone       AxisOutput `json:"emotional_tone"`
	TemporalOrientation AxisOutput `json:"temporal_orientation"`
}

// AxisOutput holds a single axis score and its sub-metric breakdown.
type AxisOutput struct {
	Score      float64                      `json:"score"`
	SubMetrics map[string]SubMetricOutput   `json:"sub_metrics"`
}

// SubMetricOutput holds a single sub-metric's raw, normalized, and weight values.
type SubMetricOutput struct {
	Raw        float64 `json:"raw"`
	Normalized float64 `json:"normalized"`
	Weight     float64 `json:"weight"`
}

// ComparisonOutput holds the optional corpus comparison.
type ComparisonOutput struct {
	CorpusID        string                         `json:"corpus_id"`
	CorpusSize      int                            `json:"corpus_size"`
	SimilarityScore float64                        `json:"similarity_score"`
	Axes            map[string]AxisComparisonOutput `json:"axes"`
}

// AxisComparisonOutput holds per-axis comparison against the corpus centroid.
type AxisComparisonOutput struct {
	DocumentScore  float64 `json:"document_score"`
	CentroidScore  float64 `json:"centroid_score"`
	CentroidStdDev float64 `json:"centroid_stddev"`
	Delta          float64 `json:"delta"`
	WithinBand     bool    `json:"within_band"`
}

// SchemaVersion is the current Inkwell signature schema version.
const SchemaVersion = "1.0.0"

// wordsPerMinute is the assumed reading speed for computing reading time.
const wordsPerMinute = 200

// ToOutput converts a Signature and DocumentInfo to the schema-compliant output.
func ToOutput(sig Signature, doc DocumentInfo) SignatureOutput {
	readingTime := doc.ReadingTimeSeconds
	if readingTime == 0 && doc.WordCount > 0 {
		readingTime = (doc.WordCount * 60) / wordsPerMinute
	}

	out := SignatureOutput{
		Version: SchemaVersion,
		Document: DocumentOutput{
			WordCount:          doc.WordCount,
			SentenceCount:      doc.SentenceCount,
			ParagraphCount:     doc.ParagraphCount,
			ReadingTimeSeconds: readingTime,
		},
		Signature: axesToOutput(sig),
	}

	if doc.ReadabilityFormula != "" {
		out.Document.Readability = &ReadabilityOutput{
			Formula: doc.ReadabilityFormula,
			Score:   doc.ReadabilityScore,
			Grade:   doc.ReadabilityGrade,
		}
	}

	return out
}

// ToOutputWithComparison adds corpus comparison to the output.
func ToOutputWithComparison(sig Signature, doc DocumentInfo, comp Comparison) SignatureOutput {
	out := ToOutput(sig, doc)
	out.Comparison = comparisonToOutput(comp)
	return out
}

func axesToOutput(sig Signature) SignatureAxesOutput {
	return SignatureAxesOutput{
		Formality:           axisToOutput(sig.Axes[Formality]),
		Confidence:          axisToOutput(sig.Axes[Confidence]),
		Rhythm:              axisToOutput(sig.Axes[Rhythm]),
		Economy:             axisToOutput(sig.Axes[Economy]),
		Precision:           axisToOutput(sig.Axes[Precision]),
		Coherence:           axisToOutput(sig.Axes[Coherence]),
		Vocabulary:          axisToOutput(sig.Axes[Vocabulary]),
		Stance:              axisToOutput(sig.Axes[Stance]),
		EmotionalTone:       axisToOutput(sig.Axes[EmotionalTone]),
		TemporalOrientation: axisToOutput(sig.Axes[TemporalOrientation]),
	}
}

func axisToOutput(a AxisResult) AxisOutput {
	subs := make(map[string]SubMetricOutput, len(a.SubMetrics))
	for name, sm := range a.SubMetrics {
		subs[name] = SubMetricOutput{
			Raw:        sm.Raw,
			Normalized: sm.Normalized,
			Weight:     sm.Weight,
		}
	}
	return AxisOutput{
		Score:      a.Score,
		SubMetrics: subs,
	}
}

func comparisonToOutput(comp Comparison) *ComparisonOutput {
	axes := make(map[string]AxisComparisonOutput, AxisCount)
	for i, ac := range comp.Axes {
		axes[AxisNames[i]] = AxisComparisonOutput{
			DocumentScore:  ac.DocumentScore,
			CentroidScore:  ac.CentroidScore,
			CentroidStdDev: ac.CentroidStdDev,
			Delta:          ac.Delta,
			WithinBand:     ac.WithinBand,
		}
	}
	return &ComparisonOutput{
		CorpusID:        comp.CorpusID,
		CorpusSize:      comp.CorpusSize,
		SimilarityScore: comp.SimilarityScore,
		Axes:            axes,
	}
}
