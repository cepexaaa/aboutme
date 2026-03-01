package application

import (
	"strings"
	"testing"

	"fractalflame/internal/domain"
)

func TestValidateConfig(t *testing.T) {
	validConfig := &domain.Config{
		Size:            domain.Size{Width: 800, Height: 600},
		Seed:            42.0,
		IterationCount:  100000,
		OutputPath:      "test.png",
		Threads:         4,
		GammaCorrection: false,
		Gamma:           2.2,
		SymmetryLevel:   1,
		Functions: []domain.FunctionConfig{
			{Name: "swirl", Weight: 1.0},
			{Name: "horseshoe", Weight: 0.5},
		},
		AffineParams: []domain.AffineParams{
			{A: 0.7, B: -0.3, C: 0.1, D: 0.3, E: 0.7, F: 0.1},
			{A: 0.4, B: 0.6, C: -0.3, D: -0.6, E: 0.4, F: 0.3},
		},
	}

	tests := []struct {
		name        string
		modifyFunc  func(*domain.Config)
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid config",
			modifyFunc: func(c *domain.Config) {},
			wantErr:    false,
		},
		{
			name: "invalid width",
			modifyFunc: func(c *domain.Config) {
				c.Size.Width = -100
			},
			wantErr:     true,
			errContains: "invalid image size",
		},
		{
			name: "invalid height",
			modifyFunc: func(c *domain.Config) {
				c.Size.Height = 0
			},
			wantErr:     true,
			errContains: "invalid image size",
		},
		{
			name: "size too large",
			modifyFunc: func(c *domain.Config) {
				c.Size.Width = 20000
				c.Size.Height = 20000
			},
			wantErr:     true,
			errContains: "image size too large",
		},
		{
			name: "invalid iteration count",
			modifyFunc: func(c *domain.Config) {
				c.IterationCount = 0
			},
			wantErr:     true,
			errContains: "iteration count must be positive",
		},
		{
			name: "iteration count too large",
			modifyFunc: func(c *domain.Config) {
				c.IterationCount = 1000000000
			},
			wantErr:     true,
			errContains: "iteration count too large",
		},
		{
			name: "invalid threads",
			modifyFunc: func(c *domain.Config) {
				c.Threads = 0
			},
			wantErr:     true,
			errContains: "thread count must be positive",
		},
		{
			name: "threads too large",
			modifyFunc: func(c *domain.Config) {
				c.Threads = 100
			},
			wantErr:     true,
			errContains: "thread count too large",
		},
		{
			name: "empty output path",
			modifyFunc: func(c *domain.Config) {
				c.OutputPath = ""
			},
			wantErr:     true,
			errContains: "output path cannot be empty",
		},
		{
			name: "wrong file extension",
			modifyFunc: func(c *domain.Config) {
				c.OutputPath = "result.jpg"
			},
			wantErr:     true,
			errContains: "output file must have .png extension",
		},
		{
			name: "no functions",
			modifyFunc: func(c *domain.Config) {
				c.Functions = []domain.FunctionConfig{}
			},
			wantErr:     true,
			errContains: "at least one transformation function",
		},
		{
			name: "invalid function weight",
			modifyFunc: func(c *domain.Config) {
				c.Functions[0].Weight = -1.0
			},
			wantErr:     true,
			errContains: "function weight must be positive",
		},
		{
			name: "unknown function",
			modifyFunc: func(c *domain.Config) {
				c.Functions[0].Name = "unknown_function"
			},
			wantErr:     true,
			errContains: "unknown_function = -1.000000",
		},
		{
			name: "total weight zero",
			modifyFunc: func(c *domain.Config) {
				c.Functions = []domain.FunctionConfig{
					{Name: "swirl", Weight: 0.0},
					{Name: "horseshoe", Weight: 0.0},
				}
			},
			wantErr:     true,
			errContains: "function weight must be positive",
		},
		{
			name: "no affine params",
			modifyFunc: func(c *domain.Config) {
				c.AffineParams = []domain.AffineParams{}
			},
			wantErr:     true,
			errContains: "function weight must be positive",
		},
		{
			name: "invalid gamma",
			modifyFunc: func(c *domain.Config) {
				c.Gamma = -1.0
			},
			wantErr:     true,
			errContains: "gamma value must be positive",
		},
		{
			name: "gamma too large",
			modifyFunc: func(c *domain.Config) {
				c.Gamma = 20.0
			},
			wantErr:     true,
			errContains: "gamma value too large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &domain.Config{}
			*config = *validConfig
			tt.modifyFunc(config)

			err := ValidateConfig(config)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error message %q doesn't contain %q", err.Error(), tt.errContains)
				}
			}
		})
	}
}
