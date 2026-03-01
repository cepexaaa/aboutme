package domain

import (
	"math/rand"
	"testing"
)

func BenchmarkFunctions(b *testing.B) {

	rng := rand.New(rand.NewSource(42))
	points := make([]Point, 1000)
	for i := range points {
		points[i] = Point{
			X: (rng.Float64() - 0.5) * 2.0,
			Y: (rng.Float64() - 0.5) * 2.0,
		}
	}

	tests := []struct {
		name string
		fn   TransformationFunction
	}{
		{"Linear", Linear{}},
		{"Swirl", Swirl{}},
		{"Horseshoe", Horseshoe{}},
		{"Disk", Disk{}},
		{"Heart", Heart{}},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {

				for j := 0; j < len(points); j++ {
					_ = tt.fn.Apply(points[j])
				}
			}
		})
	}
}

func BenchmarkAffineTransform_Apply(b *testing.B) {
	rng := rand.New(rand.NewSource(42))

	transform := AffineTransform{
		Params: AffineParams{
			A: 0.7, B: -0.3, C: 0.1,
			D: 0.3, E: 0.7, F: 0.1,
		},
	}

	points := make([]Point, 1000)
	for i := range points {
		points[i] = Point{
			X: (rng.Float64() - 0.5) * 2.0,
			Y: (rng.Float64() - 0.5) * 2.0,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(points); j++ {
			_ = transform.Apply(points[j])
		}
	}
}

func BenchmarkImageBuffer_AddPoint(b *testing.B) {
	width, height := 800, 600
	buffer := NewImageBuffer(width, height, false)

	rng := rand.New(rand.NewSource(42))
	color := Color{R: 0.5, G: 0.3, B: 0.2, A: 1.0}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := rng.Intn(width)
		y := rng.Intn(height)
		buffer.AddPoint(x, y, color, defaultSymmetry)
	}
}

func BenchmarkImageBuffer_Normalize(b *testing.B) {
	tests := []struct {
		name      string
		logScaled bool
		fillPct   float64
	}{
		{"Linear-Empty", false, 0.0},
		{"Linear-10%", false, 0.1},
		{"Linear-50%", false, 0.5},
		{"Linear-90%", false, 0.9},
		{"Log-Empty", true, 0.0},
		{"Log-10%", true, 0.1},
		{"Log-50%", true, 0.5},
		{"Log-90%", true, 0.9},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			width, height := 800, 600
			buffer := NewImageBuffer(width, height, tt.logScaled)

			rng := rand.New(rand.NewSource(42))
			numPoints := int(float64(width*height) * tt.fillPct)

			for i := 0; i < numPoints; i++ {
				x := rng.Intn(width)
				y := rng.Intn(height)
				buffer.AddPoint(x, y, Color{
					R: rng.Float64(),
					G: rng.Float64(),
					B: rng.Float64(),
					A: 1.0,
				}, defaultSymmetry)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {

				bufferCopy := *buffer
				b.StartTimer()
				bufferCopy.Normalize()
				b.StopTimer()
			}
		})
	}
}
