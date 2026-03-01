package reporter

import (
	"encoding/json"
	"os"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain"
)

type JSONReporter struct{}

func (r *JSONReporter) GenerateReport(stats *domain.Statistics, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(stats); err != nil {
		return err
	}

	return nil
}
