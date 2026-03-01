package reporter

import (
	"log/slog"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/infrastructure"
)

type Reporter interface {
	GenerateReport(stats *domain.Statistics, outputPath string) error
}

func NewReporter(format infrastructure.OutputFormat) Reporter {
	switch format {
	case infrastructure.JSONFormat:
		return &JSONReporter{}
	case infrastructure.MarkdownFormat:
		return &MarkdownReporter{}
	case infrastructure.AsciiDocFormat:
		return &AsciiDocReporter{}
	default:
		slog.Warn("Unknown format, use markdown", "format", format)
		return &MarkdownReporter{}
	}
}
