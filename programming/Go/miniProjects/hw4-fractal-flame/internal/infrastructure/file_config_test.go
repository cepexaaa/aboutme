package infrastructure

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"fractalflame/internal/domain"
)

func TestLoadConfigFromFile(t *testing.T) {

	tempDir := t.TempDir()

	tests := []struct {
		name        string
		configJSON  string
		wantConfig  *domain.Config
		wantErr     bool
		errContains string
	}{
		{
			name: "valid full config",
			configJSON: `{
                "size": {"width": 800, "height": 600},
                "seed": 42.0,
                "iteration_count": 100000,
                "output_path": "test.png",
                "threads": 4,
                "gamma_correction": true,
                "gamma": 2.4,
                "functions": [
                    {"name": "swirl", "weight": 1.0},
                    {"name": "horseshoe", "weight": 0.5}
                ],
                "affine_params": [
                    {"a": 0.7, "b": -0.3, "c": 0.1, "d": 0.3, "e": 0.7, "f": 0.1},
                    {"a": 0.4, "b": 0.6, "c": -0.3, "d": -0.6, "e": 0.4, "f": 0.3}
                ]
            }`,
			wantConfig: &domain.Config{
				Size:            domain.Size{Width: 800, Height: 600},
				Seed:            42.0,
				IterationCount:  100000,
				OutputPath:      "test.png",
				Threads:         4,
				GammaCorrection: true,
				Gamma:           2.4,
				Functions: []domain.FunctionConfig{
					{Name: "swirl", Weight: 1.0},
					{Name: "horseshoe", Weight: 0.5},
				},
				AffineParams: []domain.AffineParams{
					{A: 0.7, B: -0.3, C: 0.1, D: 0.3, E: 0.7, F: 0.1},
					{A: 0.4, B: 0.6, C: -0.3, D: -0.6, E: 0.4, F: 0.3},
				},
			},
			wantErr: false,
		},
		{
			name: "valid minimal config",
			configJSON: `{
                "size": {"width": 1920, "height": 1080},
                "iteration_count": 1000,
                "output_path": "minimal.png"
            }`,
			wantConfig: &domain.Config{
				Size:           domain.Size{Width: 1920, Height: 1080},
				IterationCount: 1000,
				OutputPath:     "minimal.png",
				Gamma:          0.0,
			},
			wantErr: false,
		},
		{
			name:        "invalid JSON",
			configJSON:  `{ invalid json }`,
			wantErr:     true,
			errContains: "failed to parse config file",
		},
		{
			name:       "missing required fields - still valid",
			configJSON: `{}`,
			wantConfig: &domain.Config{Gamma: 0.0},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tmpFile := filepath.Join(tempDir, "test_config.json")
			err := os.WriteFile(tmpFile, []byte(tt.configJSON), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			config, err := LoadConfigFromFile(tmpFile)

			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfigFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error message %q doesn't contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if config == nil {
				t.Fatal("LoadConfigFromFile() returned nil config")
			}

			if config.Size != tt.wantConfig.Size {
				t.Errorf("Size = %v, want %v", config.Size, tt.wantConfig.Size)
			}
			if config.Seed != tt.wantConfig.Seed {
				t.Errorf("Seed = %f, want %f", config.Seed, tt.wantConfig.Seed)
			}
			if config.IterationCount != tt.wantConfig.IterationCount {
				t.Errorf("IterationCount = %d, want %d", config.IterationCount, tt.wantConfig.IterationCount)
			}
			if config.OutputPath != tt.wantConfig.OutputPath {
				t.Errorf("OutputPath = %s, want %s", config.OutputPath, tt.wantConfig.OutputPath)
			}
			if config.Threads != tt.wantConfig.Threads {
				t.Errorf("Threads = %d, want %d", config.Threads, tt.wantConfig.Threads)
			}
			if config.GammaCorrection != tt.wantConfig.GammaCorrection {
				t.Errorf("GammaCorrection = %v, want %v", config.GammaCorrection, tt.wantConfig.GammaCorrection)
			}
			if config.Gamma != tt.wantConfig.Gamma {
				t.Errorf("Gamma = %f, want %f", config.Gamma, tt.wantConfig.Gamma)
			}
			if len(config.Functions) != len(tt.wantConfig.Functions) {
				t.Errorf("Functions length = %d, want %d", len(config.Functions), len(tt.wantConfig.Functions))
			}
			if len(config.AffineParams) != len(tt.wantConfig.AffineParams) {
				t.Errorf("AffineParams length = %d, want %d", len(config.AffineParams), len(tt.wantConfig.AffineParams))
			}
		})
	}
}

func TestLoadConfigFromFile_NonExistentFile(t *testing.T) {
	_, err := LoadConfigFromFile("non_existent_file.json")

	if err == nil {
		t.Fatal("LoadConfigFromFile() should return error for non-existent file")
	}

	if !strings.Contains(err.Error(), "failed to read config file") {
		t.Errorf("error message %q should contain 'failed to read config file'", err.Error())
	}
}

func TestLoadConfigFromFile_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	tmpFile := filepath.Join(tempDir, "empty.json")

	err := os.WriteFile(tmpFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	_, err = LoadConfigFromFile(tmpFile)

	if err == nil {
		t.Fatal("LoadConfigFromFile() should return error for empty file")
	}
}

func TestConfigSerialization(t *testing.T) {

	original := &domain.Config{
		Size:            domain.Size{Width: 800, Height: 600},
		Seed:            42.0,
		IterationCount:  100000,
		OutputPath:      "test.png",
		Threads:         4,
		GammaCorrection: true,
		Gamma:           2.4,
		Functions: []domain.FunctionConfig{
			{Name: "swirl", Weight: 1.0},
			{Name: "horseshoe", Weight: 0.5},
		},
		AffineParams: []domain.AffineParams{
			{A: 0.7, B: -0.3, C: 0.1, D: 0.3, E: 0.7, F: 0.1},
			{A: 0.4, B: 0.6, C: -0.3, D: -0.6, E: 0.4, F: 0.3},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	tempDir := t.TempDir()
	tmpFile := filepath.Join(tempDir, "config.json")
	err = os.WriteFile(tmpFile, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	loaded, err := LoadConfigFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loaded.Size != original.Size {
		t.Errorf("Size mismatch: got %v, want %v", loaded.Size, original.Size)
	}
	if loaded.Seed != original.Seed {
		t.Errorf("Seed mismatch: got %f, want %f", loaded.Seed, original.Seed)
	}
	if loaded.IterationCount != original.IterationCount {
		t.Errorf("IterationCount mismatch: got %d, want %d", loaded.IterationCount, original.IterationCount)
	}
	if loaded.OutputPath != original.OutputPath {
		t.Errorf("OutputPath mismatch: got %s, want %s", loaded.OutputPath, original.OutputPath)
	}
	if loaded.Threads != original.Threads {
		t.Errorf("Threads mismatch: got %d, want %d", loaded.Threads, original.Threads)
	}
	if loaded.GammaCorrection != original.GammaCorrection {
		t.Errorf("GammaCorrection mismatch: got %v, want %v", loaded.GammaCorrection, original.GammaCorrection)
	}
	if loaded.Gamma != original.Gamma {
		t.Errorf("Gamma mismatch: got %f, want %f", loaded.Gamma, original.Gamma)
	}
	if len(loaded.Functions) != len(original.Functions) {
		t.Errorf("Functions length mismatch: got %d, want %d", len(loaded.Functions), len(original.Functions))
	}
	if len(loaded.AffineParams) != len(original.AffineParams) {
		t.Errorf("AffineParams length mismatch: got %d, want %d", len(loaded.AffineParams), len(original.AffineParams))
	}
}
