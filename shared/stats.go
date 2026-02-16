package shared

import "math"

const (
	// minSamplesForStdDev is the minimum number of samples required to compute
	// a meaningful standard deviation.
	minSamplesForStdDev = 2
)

// Mean returns the arithmetic mean of values.
func Mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// StdDev returns the population standard deviation of values.
// Returns 0 if there are fewer than 2 values.
func StdDev(values []float64) float64 {
	if len(values) < minSamplesForStdDev {
		return 0
	}
	m := Mean(values)
	sum := 0.0
	for _, v := range values {
		d := v - m
		sum += d * d
	}
	return math.Sqrt(sum / float64(len(values)))
}

// CoefficientOfVariation returns the CV (std dev / mean) of values.
// Returns 0 if all values are zero. Uses absolute value of mean to handle
// negative means properly.
func CoefficientOfVariation(values []float64) float64 {
	m := Mean(values)
	sd := StdDev(values)
	if m == 0 {
		return 0
	}
	return sd / math.Abs(m)
}

// CosineSimilarity computes the cosine similarity between two float32 vectors.
// Returns 0 if either vector is zero-length or all zeros.
func CosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

// Entropy computes the Shannon entropy of a probability distribution.
// Input values should be non-negative and sum to 1 (or at least be proportional).
// Returns bits (log base 2).
func Entropy(distribution []float64) float64 {
	h := 0.0
	for _, p := range distribution {
		if p > 0 {
			h -= p * math.Log2(p)
		}
	}
	return h
}
