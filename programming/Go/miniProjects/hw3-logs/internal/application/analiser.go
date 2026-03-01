package application

import (
	"bufio"
	"log/slog"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/infrastructure"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/infrastructure/reporter"
)

type Analyzer struct {
	config   *infrastructure.Config
	reader   *infrastructure.FileReader
	parser   *Parser
	reporter reporter.Reporter
}

func NewAnalyzer(cfg *infrastructure.Config) *Analyzer {
	return &Analyzer{
		config:   cfg,
		reader:   infrastructure.NewFileReader(),
		parser:   NewParser(),
		reporter: reporter.NewReporter(cfg.Format),
	}
}

func (a *Analyzer) Analyze() error {
	readers, filenames, err := a.reader.ReadFiles(a.config.Path)
	if err != nil {
		return err
	}

	var allEntries []domain.LogEntry

	for i, reader := range readers {
		slog.Info("File Analysis", "file", filenames[i])
		scanner := bufio.NewScanner(reader)

		entries, err := a.parser.ParseFile(scanner)
		if err != nil {
			return err
		}

		allEntries = append(allEntries, entries...)
	}

	filteredEntries := a.filterByDate(allEntries)

	stats := a.calculateStatistics(filteredEntries, filenames)

	if err := a.reporter.GenerateReport(stats, a.config.Output); err != nil {
		return err
	}

	return nil
}

func (a *Analyzer) filterByDate(entries []domain.LogEntry) []domain.LogEntry {
	if a.config.From == nil && a.config.To == nil {
		return entries
	}

	var filtered []domain.LogEntry

	for _, entry := range entries {
		entryTime := entry.TimeLocal

		if a.config.From != nil && entryTime.Before(*a.config.From) {
			continue
		}

		if a.config.To != nil && entryTime.After(*a.config.To) {
			continue
		}

		filtered = append(filtered, entry)
	}

	slog.Info("Filtering by date",
		"total entries", len(entries),
		"filtered", len(filtered))
	return filtered
}
