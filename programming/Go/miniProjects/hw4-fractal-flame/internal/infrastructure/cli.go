package infrastructure

import (
	"github.com/spf13/pflag"
)

type CLIConfig struct {
	Width           int
	Height          int
	Seed            float64
	IterationCount  int
	OutputPath      string
	Threads         int
	AffineParamsStr string
	FunctionsStr    string
	ConfigPath      string
	GammaCorrection bool
	Gamma           float64
	SymmetryLevel   int
}

func ParseCLI() *CLIConfig {
	config := &CLIConfig{}

	pflag.IntVarP(&config.Width, "width", "w", 0, "Image width")
	pflag.IntVarP(&config.Height, "height", "h", 0, "Image height")
	config.Seed = *pflag.Float64("seed", 0.0, "Random seed")
	pflag.IntVarP(&config.IterationCount, "iteration-count", "i", 0, "Iteration count")
	pflag.StringVarP(&config.OutputPath, "output-path", "o", "", "Output path")
	pflag.IntVarP(&config.Threads, "threads", "t", 0, "Number of threads")
	pflag.StringVarP(&config.AffineParamsStr, "affine-params", "a", "", "Affine parameters")
	pflag.StringVarP(&config.FunctionsStr, "functions", "f", "", "Transformation functions")
	pflag.StringVar(&config.ConfigPath, "config", "", "Path to config file")
	pflag.BoolVarP(&config.GammaCorrection, "gamma-correction", "g", false, "Enable gamma correction")
	pflag.Float64Var(&config.Gamma, "gamma", 2.2, "Gamma value for correction")
	pflag.IntVarP(&config.SymmetryLevel, "symmetry-level", "s", 1, "Symmetry level (N >= 1, rotations around center)")

	pflag.Parse()

	return config
}
