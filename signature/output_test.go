package signature_test

import (
	"encoding/json"
	"testing"

	"github.com/inkcheck/signature"
)

func TestToOutput_ValidJSON(t *testing.T) {
	sig := signature.Compute(signature.RawMetrics{
		WordCount:     847,
		SentenceCount: 52,
		ParagraphCount: 9,
	})

	doc := signature.DocumentInfo{
		WordCount:          847,
		SentenceCount:      52,
		ParagraphCount:     9,
		ReadabilityFormula: "flesch_kincaid_grade",
		ReadabilityScore:   10.2,
		ReadabilityGrade:   "10th grade",
	}

	out := signature.ToOutput(sig, doc)

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal output: %v", err)
	}

	// Verify it round-trips
	var decoded signature.SignatureOutput
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if decoded.Version != signature.SchemaVersion {
		t.Errorf("version = %q, want %q", decoded.Version, signature.SchemaVersion)
	}
	if decoded.Document.WordCount != 847 {
		t.Errorf("word_count = %d, want 847", decoded.Document.WordCount)
	}
	if decoded.Document.Readability == nil {
		t.Fatal("readability should not be nil")
	}
	if decoded.Document.Readability.Formula != "flesch_kincaid_grade" {
		t.Errorf("formula = %q, want flesch_kincaid_grade", decoded.Document.Readability.Formula)
	}
}

func TestToOutput_ReadingTime(t *testing.T) {
	sig := signature.Compute(signature.RawMetrics{WordCount: 400})
	doc := signature.DocumentInfo{WordCount: 400}
	out := signature.ToOutput(sig, doc)

	// 400 words / 200 wpm = 120 seconds
	if out.Document.ReadingTimeSeconds != 120 {
		t.Errorf("reading time = %d, want 120", out.Document.ReadingTimeSeconds)
	}
}

func TestToOutputWithComparison(t *testing.T) {
	corpus := &signature.Corpus{
		ID: "brand-v2",
		Signatures: []signature.Signature{
			makeSignature([signature.AxisCount]float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5}),
		},
	}

	sig := makeSignature([signature.AxisCount]float64{0.6, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5})
	doc := signature.DocumentInfo{WordCount: 100, SentenceCount: 5, ParagraphCount: 2}
	comp := signature.Compare(sig, corpus)

	out := signature.ToOutputWithComparison(sig, doc, comp)

	if out.Comparison == nil {
		t.Fatal("comparison should not be nil")
	}
	if out.Comparison.CorpusID != "brand-v2" {
		t.Errorf("corpus_id = %q, want brand-v2", out.Comparison.CorpusID)
	}

	formalityComp, ok := out.Comparison.Axes["formality"]
	if !ok {
		t.Fatal("missing formality in comparison axes")
	}
	if formalityComp.DocumentScore != 0.6 {
		t.Errorf("formality document_score = %v, want 0.6", formalityComp.DocumentScore)
	}

	// Verify it marshals to valid JSON
	_, err := json.Marshal(out)
	if err != nil {
		t.Fatalf("failed to marshal output with comparison: %v", err)
	}
}

func TestToOutput_NoReadability(t *testing.T) {
	sig := signature.Compute(signature.RawMetrics{})
	doc := signature.DocumentInfo{}
	out := signature.ToOutput(sig, doc)

	if out.Document.Readability != nil {
		t.Error("readability should be nil when no formula provided")
	}
}
