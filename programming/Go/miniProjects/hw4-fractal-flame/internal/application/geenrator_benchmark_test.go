package application

import (
	"testing"

	"fractalflame/internal/domain"
	"fractalflame/internal/infrastructure"
)

func benchmarkConfig(width, height, iterations, threads int) *domain.Config {
	return &domain.Config{
		Size: domain.Size{
			Width:  width,
			Height: height,
		},
		Seed:            42.0,
		IterationCount:  iterations,
		OutputPath:      "benchmark.png",
		Threads:         threads,
		GammaCorrection: true,
		Gamma:           2.2,
		Functions: []domain.FunctionConfig{
			{Name: "swirl", Weight: 1.0},
			{Name: "horseshoe", Weight: 0.7},
			{Name: "linear", Weight: 0.3},
		},
		AffineParams: []domain.AffineParams{
			{A: 0.7, B: -0.3, C: 0.1, D: 0.3, E: 0.7, F: 0.1},
			{A: 0.4, B: 0.6, C: -0.3, D: -0.6, E: 0.4, F: 0.3},
			{A: 0.5, B: -0.5, C: 0.2, D: 0.5, E: 0.5, F: -0.2},
		},
	}
}

func BenchmarkGenerate_SingleThread(b *testing.B) {
	logger := infrastructure.NewLogger()
	generator := NewGenerator(logger)

	config := benchmarkConfig(800, 600, 100000, 1)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := generator.Generate(config)
		if err != nil {
			b.Fatalf("Generate failed: %v", err)
		}
	}
}

func BenchmarkGenerate_MultiThread2(b *testing.B) {
	logger := infrastructure.NewLogger()
	generator := NewGenerator(logger)

	config := benchmarkConfig(800, 600, 100000, 2)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := generator.Generate(config)
		if err != nil {
			b.Fatalf("Generate failed: %v", err)
		}
	}
}

func BenchmarkGenerate_MultiThread4(b *testing.B) {
	logger := infrastructure.NewLogger()
	generator := NewGenerator(logger)

	config := benchmarkConfig(800, 600, 100000, 4)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := generator.Generate(config)
		if err != nil {
			b.Fatalf("Generate failed: %v", err)
		}
	}
}

func BenchmarkGenerate_MultiThread8(b *testing.B) {
	logger := infrastructure.NewLogger()
	generator := NewGenerator(logger)

	config := benchmarkConfig(800, 600, 100000, 8)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := generator.Generate(config)
		if err != nil {
			b.Fatalf("Generate failed: %v", err)
		}
	}
}

func BenchmarkGenerate_MultiThread16(b *testing.B) {
	logger := infrastructure.NewLogger()
	generator := NewGenerator(logger)

	config := benchmarkConfig(800, 600, 100000, 16)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := generator.Generate(config)
		if err != nil {
			b.Fatalf("Generate failed: %v", err)
		}
	}
}

func BenchmarkGenerate_VariousSizes(b *testing.B) {
	logger := infrastructure.NewLogger()
	generator := NewGenerator(logger)

	tests := []struct {
		name    string
		width   int
		height  int
		threads int
	}{
		{"320x240-1", 320, 240, 1},
		{"640x480-1", 640, 480, 1},
		{"800x600-1", 800, 600, 1},
		{"1024x768-1", 1024, 768, 1},
		{"320x240-4", 320, 240, 4},
		{"640x480-4", 640, 480, 4},
		{"800x600-4", 800, 600, 4},
		{"1024x768-4", 1024, 768, 4},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			config := benchmarkConfig(tt.width, tt.height, 50000, tt.threads)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := generator.Generate(config)
				if err != nil {
					b.Fatalf("Generate failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkGenerate_VariousIterations(b *testing.B) {
	logger := infrastructure.NewLogger()
	generator := NewGenerator(logger)

	tests := []struct {
		name       string
		iterations int
		threads    int
	}{
		{"10k-1", 10000, 1},
		{"50k-1", 50000, 1},
		{"100k-1", 100000, 1},
		{"500k-1", 500000, 1},
		{"10k-4", 10000, 4},
		{"50k-4", 50000, 4},
		{"100k-4", 100000, 4},
		{"500k-4", 500000, 4},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			config := benchmarkConfig(800, 600, tt.iterations, tt.threads)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := generator.Generate(config)
				if err != nil {
					b.Fatalf("Generate failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkGenerate_DifferentFunctions(b *testing.B) {
	logger := infrastructure.NewLogger()
	generator := NewGenerator(logger)

	tests := []struct {
		name      string
		functions []domain.FunctionConfig
		threads   int
	}{
		{
			name: "single-linear-1",
			functions: []domain.FunctionConfig{
				{Name: "linear", Weight: 1.0},
			},
			threads: 1,
		},
		{
			name: "single-linear-4",
			functions: []domain.FunctionConfig{
				{Name: "linear", Weight: 1.0},
			},
			threads: 4,
		},
		{
			name: "swirl-only-1",
			functions: []domain.FunctionConfig{
				{Name: "swirl", Weight: 1.0},
			},
			threads: 1,
		},
		{
			name: "swirl-only-4",
			functions: []domain.FunctionConfig{
				{Name: "swirl", Weight: 1.0},
			},
			threads: 4,
		},
		{
			name: "complex-1",
			functions: []domain.FunctionConfig{
				{Name: "swirl", Weight: 1.0},
				{Name: "horseshoe", Weight: 0.7},
				{Name: "disk", Weight: 0.5},
				{Name: "heart", Weight: 0.3},
				{Name: "linear", Weight: 0.2},
			},
			threads: 1,
		},
		{
			name: "complex-4",
			functions: []domain.FunctionConfig{
				{Name: "swirl", Weight: 1.0},
				{Name: "horseshoe", Weight: 0.7},
				{Name: "disk", Weight: 0.5},
				{Name: "heart", Weight: 0.3},
				{Name: "linear", Weight: 0.2},
			},
			threads: 4,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			config := benchmarkConfig(800, 600, 100000, tt.threads)
			config.Functions = tt.functions

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := generator.Generate(config)
				if err != nil {
					b.Fatalf("Generate failed: %v", err)
				}
			}
		})
	}
}
