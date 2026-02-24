package signature

import "math"

// Curve specifies the normalization curve applied to a raw sub-metric value.
type Curve int

const (
	Linear Curve = iota
	Sqrt
	Log
)

// normalize applies range normalization with the specified curve, then
// optionally inverts the result. The output is clamped to [0, 1].
//
// Steps (per requirements.md):
//
//	Step 1 — Range normalization with curve:
//	  linear: clamp((raw - lo) / (hi - lo), 0, 1)
//	  sqrt:   clamp(sqrt((raw - lo) / (hi - lo)), 0, 1)
//	  log:    clamp(log1p(raw - lo) / log1p(hi - lo), 0, 1)
//
//	Step 2 — Direction normalization:
//	  if invert: n = 1 - n
func normalize(raw, lo, hi float64, curve Curve, invert bool) float64 {
	if hi <= lo {
		return 0
	}

	var n float64
	switch curve {
	case Linear:
		n = (raw - lo) / (hi - lo)
	case Sqrt:
		n = math.Sqrt(clampF((raw - lo) / (hi - lo)))
	case Log:
		diff := raw - lo
		if diff <= 0 {
			n = 0
		} else {
			n = math.Log1p(diff) / math.Log1p(hi-lo)
		}
	}

	n = clampF(n)

	if invert {
		n = 1 - n
	}
	return n
}

// composite computes a weighted average of normalized sub-metric values.
//
//	axis_score = sum(weight_i * n_i) / sum(weight_i)
func composite(values []float64, weights []float64) float64 {
	if len(values) == 0 || len(values) != len(weights) {
		return 0
	}
	var num, den float64
	for i := range values {
		num += weights[i] * values[i]
		den += weights[i]
	}
	if den == 0 {
		return 0
	}
	return num / den
}

func clampF(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
