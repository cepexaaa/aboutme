package infrastructure

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
)

func resetFlags() {
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
}

func TestParseCLI_ValidFlags(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedConfig *CLIConfig
	}{
		{
			name: "all short flags",
			args: []string{
				"fractalflame",
				"-w", "800",
				"-h", "600",
				"-i", "100000",
				"-o", "output.png",
				"-t", "4",
				"-f", "swirl:1.0",
				"-a", "1,2,3,4,5,6",
				"-g",
				"--gamma", "2.4",
			},
			expectedConfig: &CLIConfig{
				Width:           800,
				Height:          600,
				IterationCount:  100000,
				OutputPath:      "output.png",
				Threads:         4,
				FunctionsStr:    "swirl:1.0",
				AffineParamsStr: "1,2,3,4,5,6",
				GammaCorrection: true,
				Gamma:           2.4,
			},
		},
		{
			name: "all long flags",
			args: []string{
				"fractalflame",
				"--width", "1920",
				"--height", "1080",
				"--iteration-count", "500000",
				"--output-path", "result.png",
				"--threads", "8",
				"--functions", "swirl:1.0,horseshoe:0.5",
				"--affine-params", "1,2,3,4,5,6/7,8,9,10,11,12",
				"--gamma-correction",
				"--gamma", "1.8",
			},
			expectedConfig: &CLIConfig{
				Width:           1920,
				Height:          1080,
				IterationCount:  500000,
				OutputPath:      "result.png",
				Threads:         8,
				FunctionsStr:    "swirl:1.0,horseshoe:0.5",
				AffineParamsStr: "1,2,3,4,5,6/7,8,9,10,11,12",
				GammaCorrection: true,
				Gamma:           1.8,
			},
		},
		{
			name: "mixed flags",
			args: []string{
				"fractalflame",
				"-w", "800",
				"--height", "600",
				"-i", "100000",
				"--output-path", "mixed.png",
				"-t", "2",
				"--gamma-correction",
			},
			expectedConfig: &CLIConfig{
				Width:           800,
				Height:          600,
				IterationCount:  100000,
				OutputPath:      "mixed.png",
				Threads:         2,
				GammaCorrection: true,
				Gamma:           2.2,
			},
		},
		{
			name: "config file flag",
			args: []string{
				"fractalflame",
				"--config", "config.json",
			},
			expectedConfig: &CLIConfig{
				ConfigPath: "config.json",
				Gamma:      2.2,
			},
		},
		{
			name: "no flags - defaults",
			args: []string{"test"},
			expectedConfig: &CLIConfig{
				Gamma: 2.2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlags()

			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			os.Args = tt.args

			config := ParseCLI()

			if config.Width != tt.expectedConfig.Width {
				t.Errorf("Width = %d, want %d", config.Width, tt.expectedConfig.Width)
			}
			if config.Height != tt.expectedConfig.Height {
				t.Errorf("Height = %d, want %d", config.Height, tt.expectedConfig.Height)
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
			if config.FunctionsStr != tt.expectedConfig.FunctionsStr {
				t.Errorf("FunctionsStr = %s, want %s", config.FunctionsStr, tt.expectedConfig.FunctionsStr)
			}
			if config.AffineParamsStr != tt.expectedConfig.AffineParamsStr {
				t.Errorf("AffineParamsStr = %s, want %s", config.AffineParamsStr, tt.expectedConfig.AffineParamsStr)
			}
			if config.ConfigPath != tt.expectedConfig.ConfigPath {
				t.Errorf("ConfigPath = %s, want %s", config.ConfigPath, tt.expectedConfig.ConfigPath)
			}
			if config.GammaCorrection != tt.expectedConfig.GammaCorrection {
				t.Errorf("GammaCorrection = %v, want %v", config.GammaCorrection, tt.expectedConfig.GammaCorrection)
			}
			if config.Gamma != tt.expectedConfig.Gamma {
				t.Errorf("Gamma = %f, want %f", config.Gamma, tt.expectedConfig.Gamma)
			}
		})
	}
}
