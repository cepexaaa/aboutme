package application

import (
	"testing"

	"fractalflame/internal/domain"
)

func TestParseFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []domain.FunctionConfig
	}{
		{
			name:  "single function",
			input: "swirl:1.0",
			expected: []domain.FunctionConfig{
				{Name: "swirl", Weight: 1.0},
			},
		},
		{
			name:  "multiple functions",
			input: "swirl:1.0,horseshoe:0.5,linear:0.3",
			expected: []domain.FunctionConfig{
				{Name: "swirl", Weight: 1.0},
				{Name: "horseshoe", Weight: 0.5},
				{Name: "linear", Weight: 0.3},
			},
		},
		{
			name:  "with spaces",
			input: "swirl: 1.0 , horseshoe: 0.5 ",
			expected: []domain.FunctionConfig{
				{Name: "swirl", Weight: 1.0},
				{Name: "horseshoe", Weight: 0.5},
			},
		},
		{
			name:  "invalid weight",
			input: "swirl:invalid,horseshoe:0.5",
			expected: []domain.FunctionConfig{
				{Name: "swirl", Weight: 1.0},
				{Name: "horseshoe", Weight: 0.5},
			},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []domain.FunctionConfig{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseFunctions(tt.input)

			if len(result) != len(tt.expected) {
				t.Fatalf("got %d functions, want %d", len(result), len(tt.expected))
			}

			for i, fn := range result {
				if fn.Name != tt.expected[i].Name {
					t.Errorf("function[%d].Name = %s, want %s", i, fn.Name, tt.expected[i].Name)
				}
				if fn.Weight != tt.expected[i].Weight {
					t.Errorf("function[%d].Weight = %f, want %f", i, fn.Weight, tt.expected[i].Weight)
				}
			}
		})
	}
}

func TestParseAffineParams(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []domain.AffineParams
	}{
		{
			name:  "single affine",
			input: "1.0,2.0,3.0,4.0,5.0,6.0",
			expected: []domain.AffineParams{
				{A: 1.0, B: 2.0, C: 3.0, D: 4.0, E: 5.0, F: 6.0},
			},
		},
		{
			name:  "multiple affines",
			input: "1,2,3,4,5,6/0.1,0.2,0.3,0.4,0.5,0.6",
			expected: []domain.AffineParams{
				{A: 1.0, B: 2.0, C: 3.0, D: 4.0, E: 5.0, F: 6.0},
				{A: 0.1, B: 0.2, C: 0.3, D: 0.4, E: 0.5, F: 0.6},
			},
		},
		{
			name:  "with spaces",
			input: "1 , 2 , 3 , 4 , 5 , 6 / 0.1 , 0.2 , 0.3 , 0.4 , 0.5 , 0.6",
			expected: []domain.AffineParams{
				{A: 1.0, B: 2.0, C: 3.0, D: 4.0, E: 5.0, F: 6.0},
				{A: 0.1, B: 0.2, C: 0.3, D: 0.4, E: 0.5, F: 0.6},
			},
		},
		{
			name:     "invalid numbers",
			input:    "a,b,c,d,e,f",
			expected: []domain.AffineParams{{A: 0.0, B: 0.0, C: 0.0, D: 0.0, E: 0.0, F: 0.0}},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []domain.AffineParams{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAffineParams(tt.input)

			if len(result) != len(tt.expected) {
				t.Fatalf("got %d affine params, want %d", len(result), len(tt.expected))
			}

			for i, param := range result {
				if param != tt.expected[i] {
					t.Errorf("param[%d] = %v, want %v", i, param, tt.expected[i])
				}
			}
		})
	}
}

func TestApplyFileConfig(t *testing.T) {
	defaultConfig := &domain.Config{
		Size:            domain.Size{Width: 100, Height: 100},
		Seed:            1.0,
		IterationCount:  1000,
		OutputPath:      "default.png",
		Threads:         1,
		GammaCorrection: false,
		Gamma:           2.2,
		Functions:       []domain.FunctionConfig{{Name: "swirl", Weight: 1.0}},
		AffineParams:    []domain.AffineParams{{A: 1.0}},
	}

	tests := []struct {
		name           string
		fileConfig     *domain.Config
		expectedConfig *domain.Config
	}{
		{
			name: "override all fields",
			fileConfig: &domain.Config{
				Size:            domain.Size{Width: 200, Height: 200},
				Seed:            2.0,
				IterationCount:  2000,
				OutputPath:      "override.png",
				Threads:         4,
				GammaCorrection: true,
				Gamma:           1.8,
				Functions:       []domain.FunctionConfig{{Name: "horseshoe", Weight: 1.0}},
				AffineParams:    []domain.AffineParams{{A: 2.0}},
			},
			expectedConfig: &domain.Config{
				Size:            domain.Size{Width: 200, Height: 200},
				Seed:            2.0,
				IterationCount:  2000,
				OutputPath:      "override.png",
				Threads:         4,
				GammaCorrection: true,
				Gamma:           1.8,
				Functions:       []domain.FunctionConfig{{Name: "horseshoe", Weight: 1.0}},
				AffineParams:    []domain.AffineParams{{A: 2.0}},
			},
		},
		{
			name: "partial override",
			fileConfig: &domain.Config{
				Size:            domain.Size{Width: 150, Height: 150},
				IterationCount:  1500,
				OutputPath:      "partial.png",
				GammaCorrection: true,
			},
			expectedConfig: &domain.Config{
				Size:            domain.Size{Width: 150, Height: 150},
				Seed:            1.0,
				IterationCount:  1500,
				OutputPath:      "partial.png",
				Threads:         1,
				GammaCorrection: true,
				Gamma:           2.2,
				Functions:       []domain.FunctionConfig{{Name: "swirl", Weight: 1.0}},
				AffineParams:    []domain.AffineParams{{A: 1.0}},
			},
		},
		{
			name: "zero values not applied",
			fileConfig: &domain.Config{
				Size:           domain.Size{Width: 0, Height: 0},
				Seed:           0.0,
				IterationCount: 0,
				Threads:        0,
				Gamma:          0.0,
			},
			expectedConfig: defaultConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &domain.Config{}
			*config = *defaultConfig

			applyFileConfig(config, tt.fileConfig)

			if config.Size != tt.expectedConfig.Size {
				t.Errorf("Size = %v, want %v", config.Size, tt.expectedConfig.Size)
			}
			if config.Seed != tt.expectedConfig.Seed {
				t.Errorf("Seed = %f, want %f", config.Seed, tt.expectedConfig.Seed)
			}
			if config.IterationCount != tt.expectedConfig.IterationCount {
				t.Errorf("IterationCount = %d, want %d", config.IterationCount, tt.expectedConfig.IterationCount)
			}
			if config.OutputPath != tt.expectedConfig.OutputPath {
				t.Errorf("OutputPath = %s, want %s", config.OutputPath, tt.expectedConfig.OutputPath)
			}
			if config.Threads != tt.expectedConfig.Threads {
				t.Errorf("Threads = %d, want %d", config.Threads, tt.expectedConfig.Threads)
			}
			if config.GammaCorrection != tt.expectedConfig.GammaCorrection {
				t.Errorf("GammaCorrection = %v, want %v", config.GammaCorrection, tt.expectedConfig.GammaCorrection)
			}
			if config.Gamma != tt.expectedConfig.Gamma {
				t.Errorf("Gamma = %f, want %f", config.Gamma, tt.expectedConfig.Gamma)
			}
			if len(config.Functions) != len(tt.expectedConfig.Functions) {
				t.Errorf("Functions length = %d, want %d", len(config.Functions), len(tt.expectedConfig.Functions))
			}
			if len(config.AffineParams) != len(tt.expectedConfig.AffineParams) {
				t.Errorf("AffineParams length = %d, want %d", len(config.AffineParams), len(tt.expectedConfig.AffineParams))
			}
		})
	}
}
