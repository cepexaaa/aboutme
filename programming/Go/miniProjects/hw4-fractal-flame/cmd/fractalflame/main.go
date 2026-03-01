package main

import (
	"fmt"
	"os"

	"fractalflame/internal/application"
	"fractalflame/internal/infrastructure"
)

func main() {
	logger := infrastructure.NewLogger()

	config, err := application.LoadConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load configuration: %v", err))
		os.Exit(1)
	}

	if err := application.ValidateConfig(config); err != nil {
		logger.Error(fmt.Sprintf("Configuration validation failed: %v", err))
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf("Starting fractal flame generation with %d threads", config.Threads))
	logger.Info(fmt.Sprintf("Image size: %dx%d", config.Size.Width, config.Size.Height))
	logger.Info(fmt.Sprintf("Iterations: %d", config.IterationCount))

	generator := application.NewGenerator(logger)

	img, err := generator.Generate(config)
	if err != nil {
		logger.Error(fmt.Sprintf("Generation failed: %v", err))
		os.Exit(1)
	}

	writer := infrastructure.NewImageWriter()
	if err := writer.SaveImage(img, config.OutputPath); err != nil {
		logger.Error(fmt.Sprintf("Failed to save image: %v", err))
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf("Image successfully saved to: %s", config.OutputPath))
}
