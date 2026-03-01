package main

import (
	"log/slog"
	"os"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/application"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/infrastructure"
)

func main() {
	cfg, err := infrastructure.ParseFlags()
	if err != nil {
		slog.Error("Error: parsing parametrs", "error", err)
		os.Exit(2)
	}

	slog.Info("Launching the Log Analyzer", "path", cfg.Path, "format", cfg.Format)

	analyzer := application.NewAnalyzer(cfg)
	if err := analyzer.Analyze(); err != nil {
		slog.Error("Analysis error", "error", err)
		os.Exit(1)
	}

	slog.Info("The analysis was completed successfully")
}
