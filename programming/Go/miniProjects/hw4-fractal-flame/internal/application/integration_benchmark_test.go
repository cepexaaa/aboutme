package application

import (
	"testing"

	"fractalflame/internal/domain"
	"fractalflame/internal/infrastructure"
)

func BenchmarkComparison(b *testing.B) {
	logger := infrastructure.NewLogger()

	scenarios := []struct {
		name       string
		width      int
		height     int
		iterations int
		functions  []domain.FunctionConfig
	}{
		{
			name:       "Small-Simple",
			width:      320,
			height:     240,
			iterations: 50000,
			functions: []domain.FunctionConfig{
				{Name: "linear", Weight: 1.0},
			},
		},
		{
			name:       "Medium-Complex",
			width:      800,
			height:     600,
			iterations: 100000,
			functions: []domain.FunctionConfig{
				{Name: "swirl", Weight: 1.0},
				{Name: "horseshoe", Weight: 0.7},
			},
		},
		{
			name:       "Large-VeryComplex",
			width:      1920,
			height:     1080,
			iterations: 500000,
			functions: []domain.FunctionConfig{
				{Name: "swirl", Weight: 1.0},
				{Name: "horseshoe", Weight: 0.7},
				{Name: "disk", Weight: 0.5},
				{Name: "heart", Weight: 0.3},
				{Name: "linear", Weight: 0.2},
			},
		},
	}

	threadConfigs := []struct {
		threads int
		label   string
	}{
		{1, "1-thread"},
		{2, "2-threads"},
		{4, "4-threads"},
		{8, "8-threads"},
		{16, "16-threads"},
	}

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			for _, threadConfig := range threadConfigs {
				b.Run(threadConfig.label, func(b *testing.B) {
					generator := NewGenerator(logger)

					config := &domain.Config{
						Size: domain.Size{
							Width:  scenario.width,
							Height: scenario.height,
						},
						Seed:            42.0,
						IterationCount:  scenario.iterations,
						OutputPath:      "benchmark.png",
						Threads:         threadConfig.threads,
						GammaCorrection: true,
						Gamma:           2.2,
						Functions:       scenario.functions,
						AffineParams: []domain.AffineParams{
							{A: 0.7, B: -0.3, C: 0.1, D: 0.3, E: 0.7, F: 0.1},
							{A: 0.4, B: 0.6, C: -0.3, D: -0.6, E: 0.4, F: 0.3},
						},
					}

					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						_, err := generator.Generate(config)
						if err != nil {
							b.Fatalf("Generate failed: %v", err)
						}
					}

					b.ReportMetric(
						float64(scenario.iterations)/b.Elapsed().Seconds(),
						"iterations/sec",
					)
				})
			}
		})
	}
}

func BenchmarkMemoryUsage(b *testing.B) {
	logger := infrastructure.NewLogger()
	generator := NewGenerator(logger)

	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"320x240", 320, 240},
		{"640x480", 640, 480},
		{"800x600", 800, 600},
		{"1024x768", 1024, 768},
		{"1920x1080", 1920, 1080},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {

			buffer := domain.NewImageBuffer(tt.width, tt.height, true)

			config := benchmarkConfig(tt.width, tt.height, 100000, 4)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {

				*buffer = *domain.NewImageBuffer(tt.width, tt.height, true)

				img, err := generator.Generate(config)
				if err != nil {
					b.Fatalf("Generate failed: %v", err)
				}

				_ = img.Bounds()
			}

			memUsage := float64(tt.width*tt.height*16) / (1024 * 1024)
			b.ReportMetric(memUsage, "MB/buffer")
		})
	}
}
