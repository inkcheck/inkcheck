package signature_test

import (
	"math"
	"testing"

	"github.com/inkcheck/signature"
)

func makeSignature(scores [signature.AxisCount]float64) signature.Signature {
	var sig signature.Signature
	for i, s := range scores {
		sig.Axes[i] = signature.AxisResult{
			Score:      s,
			SubMetrics: map[string]signature.SubMetricValue{},
		}
	}
	return sig
}

func TestCorpus_Centroid(t *testing.T) {
	corpus := &signature.Corpus{
		ID: "test",
		Signatures: []signature.Signature{
			makeSignature([signature.AxisCount]float64{0.2, 0.4, 0.6, 0.8, 1.0, 0.2, 0.4, 0.6, 0.8, 1.0}),
			makeSignature([signature.AxisCount]float64{0.8, 0.6, 0.4, 0.2, 0.0, 0.8, 0.6, 0.4, 0.2, 0.0}),
		},
	}

	centroid := corpus.Centroid()
	for i := range centroid {
		if math.Abs(centroid[i]-0.5) > 1e-9 {
			t.Errorf("centroid[%d] = %v, want 0.5", i, centroid[i])
		}
	}
}

func TestCorpus_EmptyCentroid(t *testing.T) {
	corpus := &signature.Corpus{ID: "empty"}
	centroid := corpus.Centroid()
	for i, v := range centroid {
		if v != 0 {
			t.Errorf("empty centroid[%d] = %v, want 0", i, v)
		}
	}
}

func TestCorpus_StdDev(t *testing.T) {
	// Two signatures with values 0.3 and 0.7 on all axes => mean=0.5, stddev=0.2
	corpus := &signature.Corpus{
		ID: "test",
		Signatures: []signature.Signature{
			makeSignature([signature.AxisCount]float64{0.3, 0.3, 0.3, 0.3, 0.3, 0.3, 0.3, 0.3, 0.3, 0.3}),
			makeSignature([signature.AxisCount]float64{0.7, 0.7, 0.7, 0.7, 0.7, 0.7, 0.7, 0.7, 0.7, 0.7}),
		},
	}

	sd := corpus.StdDev()
	for i, v := range sd {
		if math.Abs(v-0.2) > 1e-9 {
			t.Errorf("stddev[%d] = %v, want 0.2", i, v)
		}
	}
}

func TestCorpus_ConsistencyScore(t *testing.T) {
	// Identical signatures => stddev=0 on all axes => consistency=1.0
	corpus := &signature.Corpus{
		ID: "test",
		Signatures: []signature.Signature{
			makeSignature([signature.AxisCount]float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5}),
			makeSignature([signature.AxisCount]float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5}),
		},
	}

	score := corpus.ConsistencyScore()
	if math.Abs(score-1.0) > 1e-9 {
		t.Errorf("identical signatures: consistency = %v, want 1.0", score)
	}
}

func TestCompare(t *testing.T) {
	corpus := &signature.Corpus{
		ID: "brand-v1",
		Signatures: []signature.Signature{
			makeSignature([signature.AxisCount]float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5}),
			makeSignature([signature.AxisCount]float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5}),
		},
	}

	doc := makeSignature([signature.AxisCount]float64{0.6, 0.4, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5})
	comp := signature.Compare(doc, corpus)

	if comp.CorpusID != "brand-v1" {
		t.Errorf("CorpusID = %q, want %q", comp.CorpusID, "brand-v1")
	}
	if comp.CorpusSize != 2 {
		t.Errorf("CorpusSize = %d, want 2", comp.CorpusSize)
	}
	if comp.SimilarityScore <= 0 || comp.SimilarityScore > 1 {
		t.Errorf("SimilarityScore = %v, want (0, 1]", comp.SimilarityScore)
	}

	// Axis 0 (formality): doc=0.6, centroid=0.5, delta=0.1
	if math.Abs(comp.Axes[0].Delta-0.1) > 1e-9 {
		t.Errorf("formality delta = %v, want 0.1", comp.Axes[0].Delta)
	}
	// Axis 1 (confidence): doc=0.4, centroid=0.5, delta=-0.1
	if math.Abs(comp.Axes[1].Delta-(-0.1)) > 1e-9 {
		t.Errorf("confidence delta = %v, want -0.1", comp.Axes[1].Delta)
	}

	// WithinBand: stddev=0 for identical corpus, so delta != 0 => not within band
	if comp.Axes[0].WithinBand {
		t.Error("axis 0 should not be within band when stddev=0 and delta>0")
	}
}
