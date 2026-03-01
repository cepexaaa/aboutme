package application

import (
	"image"
	"math/rand"
	"testing"

	"fractalflame/internal/domain"
	"fractalflame/internal/infrastructure"
)

func TestGenerator_CreateAffineTransforms(t *testing.T) {
	logger := infrastructure.NewLogger()
	generator := NewGenerator(logger)

	tests := []struct {
		name          string
		config        *domain.Config
		expectedCount int
	}{
		{
			name: "with predefined params",
			config: &domain.Config{
				AffineParams: []domain.AffineParams{
					{A: 0.7, B: -0.3, C: 0.1, D: 0.3, E: 0.7, F: 0.1},
					{A: 0.4, B: 0.6, C: -0.3, D: -0.6, E: 0.4, F: 0.3},
				},
			},
			expectedCount: 2,
		},
		{
			name: "empty params - should create random",
			config: &domain.Config{
				Seed:         42.0,
				AffineParams: []domain.AffineParams{},
			},
			expectedCount: 0,
		},
		{
			name: "zero params - should create random",
			config: &domain.Config{
				Seed: 42.0,
				AffineParams: []domain.AffineParams{
					{A: 0, B: 0, C: 0, D: 0, E: 0, F: 0},
				},
			},
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transforms := generator.createAffineTransforms(tt.config, rand.New(rand.NewSource(int64(tt.config.Seed))))

			if len(transforms) != tt.expectedCount {
				t.Errorf("got %d transforms, want %d", len(transforms), tt.expectedCount)
			}

			for i, transform := range transforms {
				if transform.Color.R == 0 && transform.Color.G == 0 && transform.Color.B == 0 {
					t.Errorf("transform %d has zero color", i)
				}
			}
		})
	}
}

func TestGenerator_CreateFunctions(t *testing.T) {
	logger := infrastructure.NewLogger()
	generator := NewGenerator(logger)

	tests := []struct {
		name          string
		config        *domain.Config
		expectedCount int
		expectedNames []string
	}{
		{
			name: "single function",
			config: &domain.Config{
				Functions: []domain.FunctionConfig{
					{Name: "swirl", Weight: 1.0},
				},
			},
			expectedCount: 1,
			expectedNames: []string{"swirl"},
		},
		{
			name: "multiple functions",
			config: &domain.Config{
				Functions: []domain.FunctionConfig{
					{Name: "swirl", Weight: 1.0},
					{Name: "horseshoe", Weight: 0.5},
					{Name: "linear", Weight: 0.3},
				},
			},
			expectedCount: 3,
			expectedNames: []string{"swirl", "horseshoe", "linear"},
		},
		{
			name: "unknown function should be skipped",
			config: &domain.Config{
				Functions: []domain.FunctionConfig{
					{Name: "swirl", Weight: 1.0},
					{Name: "unknown", Weight: 0.5},
					{Name: "linear", Weight: 0.3},
				},
			},
			expectedCount: 2,
			expectedNames: []string{"swirl", "linear"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			functions := generator.createFunctions(tt.config)

			if len(functions) != tt.expectedCount {
				t.Errorf("got %d functions, want %d", len(functions), tt.expectedCount)
			}

			for i, fn := range functions {
				if fn.weight <= 0 {
					t.Errorf("function[%d] weight = %f, should be positive", i, fn.weight)
				}
			}
		})
	}
}

func TestGenerator_NormalizeFunctionWeights(t *testing.T) {
	logger := infrastructure.NewLogger()
	generator := NewGenerator(logger)

	tests := []struct {
		name        string
		functions   []weightedFunc
		configFuncs []domain.FunctionConfig
		checkFunc   func([]weightedFunc) bool
	}{
		{
			name: "single function normalization",
			functions: []weightedFunc{
				{fn: domain.Swirl{}, weight: 1.0},
			},
			configFuncs: []domain.FunctionConfig{
				{Name: "swirl", Weight: 1.0},
			},
			checkFunc: func(funcs []weightedFunc) bool {
				return len(funcs) == 1 && funcs[0].weight == 1.0
			},
		},
		{
			name: "multiple functions normalization",
			functions: []weightedFunc{
				{fn: domain.Swirl{}, weight: 2.0},
				{fn: domain.Horseshoe{}, weight: 1.0},
				{fn: domain.Linear{}, weight: 1.0},
			},
			configFuncs: []domain.FunctionConfig{
				{Name: "swirl", Weight: 2.0},
				{Name: "horseshoe", Weight: 1.0},
				{Name: "linear", Weight: 1.0},
			},
			checkFunc: func(funcs []weightedFunc) bool {
				if len(funcs) != 3 {
					return false
				}

				total := funcs[0].weight + funcs[1].weight + funcs[2].weight
				return total > 0.999 && total < 1.001
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normalized := generator.normalizeFunctionWeights(tt.functions, tt.configFuncs)

			if !tt.checkFunc(normalized) {
				t.Errorf("NormalizeFunctionWeights() failed check for test: %s", tt.name)
			}
		})
	}
}

func TestGenerator_SelectFunctionIndex(t *testing.T) {
	logger := infrastructure.NewLogger()
	generator := NewGenerator(logger)

	functions := []weightedFunc{
		{fn: domain.Swirl{}, weight: 0.5},
		{fn: domain.Horseshoe{}, weight: 0.3},
		{fn: domain.Linear{}, weight: 0.2},
	}

	normalized := make([]weightedFunc, len(functions))
	copy(normalized, functions)

	tests := []struct {
		name       string
		mockRandom float64
		expected   int
	}{
		{"select first (0-0.5)", 0.25, 0},
		{"select second (0.5-0.8)", 0.6, 1},
		{"select third (0.8-1.0)", 0.9, 2},
		{"exact boundary 0.5", 0.5, 1},
		{"exact boundary 0.8", 0.8, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			oldRand := RandFloat64
			RandFloat64 = func() float64 { return tt.mockRandom }
			defer func() { RandFloat64 = oldRand }()

			idx := generator.selectFunctionIndex(normalized)

			if idx >= len(normalized) {
				t.Errorf("SelectFunctionIndex() = %d, want %d", idx, tt.expected)
			}
		})
	}
}

var RandFloat64 = func() float64 { return 0.5 }

func TestGenerator_ScaleToImage(t *testing.T) {
	logger := infrastructure.NewLogger()
	generator := NewGenerator(logger)

	tests := []struct {
		name                 string
		point                domain.Point
		width, height        int
		expectedX, expectedY int
	}{
		{
			name:  "center point",
			point: domain.Point{X: 0, Y: 0},
			width: 100, height: 100,
			expectedX: 49, expectedY: 49,
		},
		{
			name:  "top left",
			point: domain.Point{X: -3, Y: -3},
			width: 100, height: 100,
			expectedX: 0, expectedY: 0,
		},
		{
			name:  "bottom right",
			point: domain.Point{X: 3, Y: 3},
			width: 100, height: 100,
			expectedX: 99, expectedY: 99,
		},
		{
			name:  "outside bounds - should clamp",
			point: domain.Point{X: 10, Y: -10},
			width: 100, height: 100,
			expectedX: 99, expectedY: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := generator.scaleToImage(tt.point, tt.width, tt.height)

			if x != tt.expectedX {
				t.Errorf("ScaleToImage() X = %d, want %d", x, tt.expectedX)
			}
			if y != tt.expectedY {
				t.Errorf("ScaleToImage() Y = %d, want %d", y, tt.expectedY)
			}
		})
	}
}

type mockImageWriter struct {
	saveCalled bool
	lastImage  image.Image
	lastPath   string
	saveError  error
}

func (m *mockImageWriter) SaveImage(img image.Image, path string) error {
	m.saveCalled = true
	m.lastImage = img
	m.lastPath = path
	return m.saveError
}

func TestGenerator_Integration(t *testing.T) {
	logger := infrastructure.NewLogger()

	config := &domain.Config{
		Size:            domain.Size{Width: 10, Height: 10},
		Seed:            42.0,
		IterationCount:  1000,
		OutputPath:      "test_output.png",
		Threads:         1,
		GammaCorrection: true,
		Gamma:           2.2,
		Functions: []domain.FunctionConfig{
			{Name: "linear", Weight: 1.0},
		},
		AffineParams: []domain.AffineParams{
			{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.0},
			{A: 0.5, B: 0.0, C: 0.5, D: 0.0, E: 0.5, F: 0.0},
		},
	}

	generator := NewGenerator(logger)

	t.Run("generate single thread", func(t *testing.T) {
		img, err := generator.Generate(config)

		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		if img == nil {
			t.Fatal("Generate() returned nil image")
		}

		bounds := img.Bounds()
		if bounds.Dx() != config.Size.Width || bounds.Dy() != config.Size.Height {
			t.Errorf("Image size = %v, want %dx%d", bounds, config.Size.Width, config.Size.Height)
		}

		hasNonBlack := false
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				if r > 0 || g > 0 || b > 0 {
					hasNonBlack = true
					break
				}
			}
			if hasNonBlack {
				break
			}
		}

		if !hasNonBlack {
			t.Error("Generated image is completely black")
		}
	})

	t.Run("generate multi thread", func(t *testing.T) {
		config.Threads = 2
		img, err := generator.Generate(config)

		if err != nil {
			t.Fatalf("Generate() with threads error = %v", err)
		}

		if img == nil {
			t.Fatal("Generate() with threads returned nil image")
		}
	})
}
