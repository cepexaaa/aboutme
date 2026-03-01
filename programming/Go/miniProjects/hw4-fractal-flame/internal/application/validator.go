package application

import (
	"fmt"
	"path/filepath"

	"fractalflame/internal/domain"
)

func ValidateConfig(config *domain.Config) error {
	if config.Size.Width <= 0 || config.Size.Height <= 0 {
		return fmt.Errorf("invalid image size: %dx%d", config.Size.Width, config.Size.Height)
	}
	if config.Size.Width > 10000 || config.Size.Height > 10000 {
		return fmt.Errorf("image size too large: %dx%d", config.Size.Width, config.Size.Height)
	}
	if config.IterationCount <= 0 {
		return fmt.Errorf("iteration count must be positive: %d", config.IterationCount)
	}
	if config.IterationCount > 100000000 {
		return fmt.Errorf("iteration count too large: %d", config.IterationCount)
	}
	if config.Threads <= 0 {
		return fmt.Errorf("thread count must be positive: %d", config.Threads)
	}
	if config.Threads > 32 {
		return fmt.Errorf("thread count too large: %d", config.Threads)
	}
	if config.OutputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}
	ext := filepath.Ext(config.OutputPath)
	if ext != ".png" && ext != ".PNG" {
		return fmt.Errorf("output file must have .png extension: %s", config.OutputPath)
	}
	if len(config.Functions) == 0 {
		return fmt.Errorf("at least one transformation function must be specified")
	}
	if config.Gamma <= 0 {
		return fmt.Errorf("gamma value must be positive: %f", config.Gamma)
	}

	if config.Gamma > 10 {
		return fmt.Errorf("gamma value too large: %f", config.Gamma)
	}
	if config.SymmetryLevel < 1 {
		return fmt.Errorf("symmetry level must be >= 1: %d", config.SymmetryLevel)
	}

	if config.SymmetryLevel > 360 {
		return fmt.Errorf("symmetry level too high: %d", config.SymmetryLevel)
	}
	totalWeight := 0.0
	for _, f := range config.Functions {
		if f.Weight <= 0 {
			return fmt.Errorf("function weight must be positive: %s = %f", f.Name, f.Weight)
		}
		if _, err := domain.GetFunctionByName(f.Name); err != nil {
			return fmt.Errorf("unknown function: %s", f.Name)
		}
		totalWeight += f.Weight
	}
	if totalWeight <= 0 {
		return fmt.Errorf("total function weight must be positive")
	}
	if len(config.AffineParams) == 0 {
		return fmt.Errorf("at least one affine parameter set must be specified")
	}
	return nil
}
