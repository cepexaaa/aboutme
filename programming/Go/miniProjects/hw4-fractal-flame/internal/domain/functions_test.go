package domain

import (
	"math"
	"testing"
)

func TestSwirl_Apply(t *testing.T) {
	tests := []struct {
		name     string
		point    Point
		expected Point
	}{
		{
			name:     "origin",
			point:    Point{X: 0, Y: 0},
			expected: Point{X: 0, Y: 0},
		},
		{
			name:     "point on x-axis",
			point:    Point{X: 1, Y: 0},
			expected: Point{X: math.Sin(1), Y: math.Cos(1)},
		},
		{
			name:     "point on y-axis",
			point:    Point{X: 0, Y: 1},
			expected: Point{X: -math.Cos(1), Y: math.Sin(1)},
		},
		{
			name:  "arbitrary point",
			point: Point{X: 0.5, Y: 0.3},
			expected: func() Point {
				r2 := 0.5*0.5 + 0.3*0.3
				sinr2 := math.Sin(r2)
				cosr2 := math.Cos(r2)
				return Point{
					X: 0.5*sinr2 - 0.3*cosr2,
					Y: 0.5*cosr2 + 0.3*sinr2,
				}
			}(),
		},
	}

	swirl := Swirl{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := swirl.Apply(tt.point)

			if math.Abs(result.X-tt.expected.X) > 1e-10 || math.Abs(result.Y-tt.expected.Y) > 1e-10 {
				t.Errorf("Swirl.Apply() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestHorseshoe_Apply(t *testing.T) {
	tests := []struct {
		name     string
		point    Point
		expected Point
	}{
		{
			name:     "point on x-axis",
			point:    Point{X: 1, Y: 0},
			expected: Point{X: 1, Y: 0},
		},
		{
			name:     "point on y-axis",
			point:    Point{X: 0, Y: 1},
			expected: Point{X: -1, Y: 0},
		},
		{
			name:     "point at 45 degrees",
			point:    Point{X: math.Sqrt(2) / 2, Y: math.Sqrt(2) / 2},
			expected: Point{X: 0, Y: 1},
		},
	}

	horseshoe := Horseshoe{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := horseshoe.Apply(tt.point)

			if math.Abs(result.X-tt.expected.X) > 1e-10 || math.Abs(result.Y-tt.expected.Y) > 1e-10 {
				t.Errorf("Horseshoe.Apply() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestLinear_Apply(t *testing.T) {
	tests := []struct {
		name  string
		point Point
	}{
		{"origin", Point{X: 0, Y: 0}},
		{"positive", Point{X: 1.5, Y: 2.5}},
		{"negative", Point{X: -1.5, Y: -2.5}},
		{"mixed", Point{X: 1.5, Y: -2.5}},
	}

	linear := Linear{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := linear.Apply(tt.point)

			if result.X != tt.point.X || result.Y != tt.point.Y {
				t.Errorf("Linear.Apply() = %v, expected %v", result, tt.point)
			}
		})
	}
}

func TestGetFunctionByName(t *testing.T) {
	tests := []struct {
		name         string
		functionName string
		wantErr      bool
	}{
		{"swirl", "swirl", false},
		{"horseshoe", "horseshoe", false},
		{"disk", "disk", false},
		{"heart", "heart", false},
		{"linear", "linear", false},
		{"unknown", "unknown_function", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := GetFunctionByName(tt.functionName)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetFunctionByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				p := Point{X: 2, Y: 3}
				p2 := f.Apply(p)
				if p == p2 && tt.functionName != "linear" {
					t.Errorf("function %v doesn't work", tt.functionName)
				}
			}
		})
	}
}
