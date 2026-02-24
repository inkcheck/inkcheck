package signature

import (
	"math"
	"testing"
)

func TestNormalize_Linear(t *testing.T) {
	tests := []struct {
		name   string
		raw    float64
		lo, hi float64
		invert bool
		want   float64
	}{
		{"midpoint", 0.5, 0, 1, false, 0.5},
		{"at lo", 0, 0, 1, false, 0},
		{"at hi", 1, 0, 1, false, 1},
		{"below lo clamps", -0.5, 0, 1, false, 0},
		{"above hi clamps", 2.0, 0, 1, false, 1},
		{"inverted midpoint", 0.5, 0, 1, true, 0.5},
		{"inverted at lo", 0, 0, 1, true, 1},
		{"inverted at hi", 1, 0, 1, true, 0},
		{"custom range", 10, 5, 15, false, 0.5},
		{"contraction hi=0.15", 0.075, 0, 0.15, true, 0.5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalize(tc.raw, tc.lo, tc.hi, Linear, tc.invert)
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("normalize(%v, %v, %v, Linear, %v) = %v, want %v",
					tc.raw, tc.lo, tc.hi, tc.invert, got, tc.want)
			}
		})
	}
}

func TestNormalize_Sqrt(t *testing.T) {
	// sqrt(0.25) = 0.5
	got := normalize(1.25, 0, 5, Sqrt, false)
	want := math.Sqrt(0.25)
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("normalize(1.25, 0, 5, Sqrt, false) = %v, want %v", got, want)
	}

	// inverted
	got = normalize(1.25, 0, 5, Sqrt, true)
	want = 1 - math.Sqrt(0.25)
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("inverted sqrt: got %v, want %v", got, want)
	}
}

func TestNormalize_Log(t *testing.T) {
	got := normalize(22.5, 5, 40, Log, false)
	want := math.Log1p(17.5) / math.Log1p(35)
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("normalize(22.5, 5, 40, Log, false) = %v, want %v", got, want)
	}
}

func TestNormalize_EqualLoHi(t *testing.T) {
	got := normalize(5, 5, 5, Linear, false)
	if got != 0 {
		t.Errorf("equal lo/hi should return 0, got %v", got)
	}
}

func TestComposite(t *testing.T) {
	values := []float64{0.4, 0.6, 0.8}
	weights := []float64{1, 1, 1}
	got := composite(values, weights)
	want := 0.6
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("composite = %v, want %v", got, want)
	}
}

func TestComposite_Weighted(t *testing.T) {
	values := []float64{0.0, 1.0}
	weights := []float64{0.25, 0.75}
	got := composite(values, weights)
	want := 0.75
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("weighted composite = %v, want %v", got, want)
	}
}

func TestComposite_Empty(t *testing.T) {
	got := composite(nil, nil)
	if got != 0 {
		t.Errorf("empty composite should be 0, got %v", got)
	}
}
