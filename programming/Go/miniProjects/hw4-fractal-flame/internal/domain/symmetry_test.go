package domain

import (
	"math"
	"testing"
)

func TestNewSymmetryTransformer(t *testing.T) {
	tests := []struct {
		name  string
		level int
	}{
		{"level 1", 1},
		{"level 2", 2},
		{"level 3", 3},
		{"level 4", 4},
		{"level 6", 6},
		{"level 8", 8},
		{"invalid level 0", 0},
		{"invalid level -1", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := NewSymmetryTransformer(tt.level)

			expectedLevel := tt.level
			if expectedLevel < 1 {
				expectedLevel = 1
			}

			if st.GetLevel() != expectedLevel {
				t.Errorf("GetLevel() = %d, want %d", st.GetLevel(), expectedLevel)
			}

			expectedAngle := (2 * math.Pi) / float64(expectedLevel)
			if math.Abs(st.angle-expectedAngle) > 1e-10 {
				t.Errorf("Angle = %f, want %f",
					st.angle, expectedAngle)
			}
		})
	}
}

func TestSymmetryTransformer_Apply(t *testing.T) {
	tests := []struct {
		name      string
		level     int
		point     Point
		checkFunc func([]Point) bool
	}{
		{
			name:  "level 1 - no symmetry",
			level: 1,
			point: Point{X: 1, Y: 0},
			checkFunc: func(points []Point) bool {
				return len(points) == 1 &&
					points[0].X == 1 && points[0].Y == 0
			},
		},
		{
			name:  "level 2 - 180 degree symmetry",
			level: 2,
			point: Point{X: 1, Y: 0},
			checkFunc: func(points []Point) bool {
				if len(points) != 2 {
					return false
				}

				if points[0].X != 1 || points[0].Y != 0 {
					return false
				}

				if math.Abs(points[1].X-(-1)) > 1e-10 ||
					math.Abs(points[1].Y-0) > 1e-10 {
					return false
				}
				return true
			},
		},
		{
			name:  "level 4 - 90 degree symmetry",
			level: 4,
			point: Point{X: 1, Y: 0},
			checkFunc: func(points []Point) bool {
				if len(points) != 4 {
					return false
				}

				for _, p := range points {
					dist := math.Sqrt(p.X*p.X + p.Y*p.Y)
					if math.Abs(dist-1) > 1e-10 {
						return false
					}
				}
				return true
			},
		},
		{
			name:  "level 3 - 120 degree symmetry at origin",
			level: 3,
			point: Point{X: 0, Y: 0},
			checkFunc: func(points []Point) bool {
				if len(points) != 3 {
					return false
				}

				for _, p := range points {
					if p.X != 0 || p.Y != 0 {
						return false
					}
				}
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := NewSymmetryTransformer(tt.level)
			points := st.Apply(tt.point)

			if !tt.checkFunc(points) {
				t.Errorf("Apply() failed for test: %s", tt.name)
			}
		})
	}
}

func TestSymmetryTransformer_RotatePoint(t *testing.T) {
	st := NewSymmetryTransformer(4)

	tests := []struct {
		name     string
		point    Point
		angle    float64
		expected Point
	}{
		{
			name:     "rotate 90 degrees",
			point:    Point{X: 1, Y: 0},
			angle:    math.Pi / 2,
			expected: Point{X: 0, Y: 1},
		},
		{
			name:     "rotate 180 degrees",
			point:    Point{X: 1, Y: 0},
			angle:    math.Pi,
			expected: Point{X: -1, Y: 0},
		},
		{
			name:     "rotate 270 degrees",
			point:    Point{X: 1, Y: 0},
			angle:    3 * math.Pi / 2,
			expected: Point{X: 0, Y: -1},
		},
		{
			name:     "rotate 360 degrees",
			point:    Point{X: 1, Y: 0},
			angle:    2 * math.Pi,
			expected: Point{X: 1, Y: 0},
		},
		{
			name:  "rotate arbitrary point",
			point: Point{X: 0.5, Y: 0.3},
			angle: math.Pi / 4,
			expected: func() Point {
				sin := math.Sin(math.Pi / 4)
				cos := math.Cos(math.Pi / 4)
				return Point{
					X: 0.5*cos - 0.3*sin,
					Y: 0.5*sin + 0.3*cos,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := st.rotatePoint(tt.point, tt.angle)

			if math.Abs(result.X-tt.expected.X) > 1e-10 ||
				math.Abs(result.Y-tt.expected.Y) > 1e-10 {
				t.Errorf("rotatePoint() = %v, want %v", result, tt.expected)
			}
		})
	}
}
