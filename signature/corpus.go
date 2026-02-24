package signature

import "math"

// Corpus holds a collection of signatures for computing centroids and comparisons.
type Corpus struct {
	ID         string
	Signatures []Signature
}

// Centroid computes the mean signature vector across all corpus signatures.
// Returns a zero signature if the corpus is empty.
func (c *Corpus) Centroid() [AxisCount]float64 {
	var centroid [AxisCount]float64
	n := len(c.Signatures)
	if n == 0 {
		return centroid
	}
	for _, sig := range c.Signatures {
		v := sig.Vector()
		for i := range centroid {
			centroid[i] += v[i]
		}
	}
	for i := range centroid {
		centroid[i] /= float64(n)
	}
	return centroid
}

// StdDev computes the per-axis standard deviation across corpus signatures.
func (c *Corpus) StdDev() [AxisCount]float64 {
	var sd [AxisCount]float64
	n := len(c.Signatures)
	if n < 2 {
		return sd
	}
	centroid := c.Centroid()
	for _, sig := range c.Signatures {
		v := sig.Vector()
		for i := range sd {
			d := v[i] - centroid[i]
			sd[i] += d * d
		}
	}
	for i := range sd {
		sd[i] = math.Sqrt(sd[i] / float64(n))
	}
	return sd
}

// ConsistencyScore computes the overall brand voice consistency.
// It is the mean of (1 - stddev) across all axes; a tight cluster = high consistency.
func (c *Corpus) ConsistencyScore() float64 {
	sd := c.StdDev()
	sum := 0.0
	for _, s := range sd {
		sum += 1 - s
	}
	return sum / float64(AxisCount)
}

// AxisComparison holds the comparison of a single axis against the corpus centroid.
type AxisComparison struct {
	DocumentScore float64
	CentroidScore float64
	CentroidStdDev float64
	Delta          float64
	WithinBand     bool
}

// Comparison holds the full comparison of a document signature against a corpus.
type Comparison struct {
	CorpusID        string
	CorpusSize      int
	SimilarityScore float64
	Axes            [AxisCount]AxisComparison
}

// Compare compares a document signature against a corpus, producing per-axis
// deltas, within-band flags, and an overall cosine similarity score.
func Compare(doc Signature, corpus *Corpus) Comparison {
	centroid := corpus.Centroid()
	sd := corpus.StdDev()
	docVec := doc.Vector()

	comp := Comparison{
		CorpusID:        corpus.ID,
		CorpusSize:      len(corpus.Signatures),
		SimilarityScore: CosineSimilarity(docVec, centroid),
	}

	for i := range comp.Axes {
		delta := docVec[i] - centroid[i]
		comp.Axes[i] = AxisComparison{
			DocumentScore:  docVec[i],
			CentroidScore:  centroid[i],
			CentroidStdDev: sd[i],
			Delta:          delta,
			WithinBand:     math.Abs(delta) <= sd[i],
		}
	}

	return comp
}
