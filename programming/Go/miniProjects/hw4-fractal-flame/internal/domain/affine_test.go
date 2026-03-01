package domain

import (
	"math"
	"math/rand"
	"testing"
)

func TestAffineTransform_Apply(t *testing.T) {
	tests := []struct {
		name     string
		params   AffineParams
		point    Point
		expected Point
	}{
		{
			name: "identity transform",
			params: AffineParams{
				A: 1, B: 0, C: 0,
				D: 0, E: 1, F: 0,
			},
			point:    Point{X: 1, Y: 2},
			expected: Point{X: 1, Y: 2},
		},
		{
			name: "translation only",
			params: AffineParams{
				A: 1, B: 0, C: 5,
				D: 0, E: 1, F: 10,
			},
			point:    Point{X: 1, Y: 2},
			expected: Point{X: 6, Y: 12},
		},
		{
			name: "scale only",
			params: AffineParams{
				A: 2, B: 0, C: 0,
				D: 0, E: 3, F: 0,
			},
			point:    Point{X: 1, Y: 2},
			expected: Point{X: 2, Y: 6},
		},
		{
			name: "rotation 90 degrees",
			params: AffineParams{
				A: 0, B: -1, C: 0,
				D: 1, E: 0, F: 0,
			},
			point:    Point{X: 1, Y: 0},
			expected: Point{X: 0, Y: 1},
		},
		{
			name: "combined transform",
			params: AffineParams{
				A: 2, B: 1, C: 5,
				D: 3, E: 4, F: 6,
			},
			point:    Point{X: 1, Y: 2},
			expected: Point{X: 2*1 + 1*2 + 5, Y: 3*1 + 4*2 + 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transform := AffineTransform{Params: tt.params}
			result := transform.Apply(tt.point)

			if math.Abs(result.X-tt.expected.X) > 1e-10 || math.Abs(result.Y-tt.expected.Y) > 1e-10 {
				t.Errorf("Apply() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestNewRandomAffine(t *testing.T) {
	tests := []struct {
		name string
		seed float64
	}{
		{"seed 0", 0.0},
		{"seed 1", 1.0},
		{"seed pi", math.Pi},
		{"seed negative", -1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rnd := rand.Rand(*rand.New(rand.NewSource(int64(tt.seed))))
			transform := NewRandomAffine(tt.seed, &rnd)
			if transform.Params.A == 0 && transform.Params.B == 0 &&
				transform.Params.C == 0 && transform.Params.D == 0 &&
				transform.Params.E == 0 && transform.Params.F == 0 {
				t.Error("NewRandomAffine() returned zero transform")
			}
			if transform.Color.R == 0 && transform.Color.G == 0 && transform.Color.B == 0 {
				t.Error("NewRandomAffine() returned zero color")
			}
		})
	}
}
