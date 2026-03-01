package reporter

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain"
)

type MarkdownReporter struct{}

func (r *MarkdownReporter) GenerateReport(stats *domain.Statistics, outputPath string) error {
	var builder strings.Builder

	builder.WriteString("# Log analysis NGINX\n\n")
	builder.WriteString(fmt.Sprintf("> Generated: %s\n\n", time.Now().Format("02.01.2006 15:04:05")))

	builder.WriteString("## General information\n\n")
	builder.WriteString("| Metric | Value |\n")
	builder.WriteString("|---------|----------|\n")
	builder.WriteString(fmt.Sprintf("| Files | `%s` |\n", strings.Join(stats.Files, ", ")))
	builder.WriteString(fmt.Sprintf("| Number of requests | %d |\n", stats.TotalRequestsCount))
	builder.WriteString(fmt.Sprintf("| Average response size | %.2fb |\n", stats.ResponseSizeInBytes.Average))
	builder.WriteString(fmt.Sprintf("| Maximum response size | %.1fb |\n", stats.ResponseSizeInBytes.Max))
	builder.WriteString(fmt.Sprintf("| 95p response size | %.2fb |\n", stats.ResponseSizeInBytes.P95))

	if len(stats.Resources) > 0 {
		builder.WriteString("\n## Requested resources\n\n")
		builder.WriteString("| Resource | Count |\n")
		builder.WriteString("|--------|------------|\n")
		for _, resource := range stats.Resources {
			builder.WriteString(fmt.Sprintf("| `%s` | %d |\n", resource.Resource, resource.TotalRequestsCount))
		}
	}

	if len(stats.ResponseCodes) > 0 {
		builder.WriteString("\n## Response codes\n\n")
		builder.WriteString("| Code | Count |\n")
		builder.WriteString("|-----|------------|\n")
		for _, code := range stats.ResponseCodes {
			builder.WriteString(fmt.Sprintf("| %d | %d |\n", code.Code, code.TotalResponsesCount))
		}
	}

	if len(stats.RequestsPerDate) > 0 {
		builder.WriteString("\n## Distribution of requests by day\n\n")
		builder.WriteString("| Date | Day of week | Requests | Percent |\n")
		builder.WriteString("|------|-------------|---------|------|\n")
		for _, date := range stats.RequestsPerDate {
			builder.WriteString(fmt.Sprintf("| %s | %s | %d | %.2f%% |\n",
				date.Date, date.Weekday, date.TotalRequestsCount, date.TotalRequestsPercentage))
		}
	}

	if len(stats.UniqueProtocols) > 0 {
		builder.WriteString("\n## Unique Protocols\n\n")
		builder.WriteString(strings.Join(stats.UniqueProtocols, ", "))
		builder.WriteString("\n")
	}

	return os.WriteFile(outputPath, []byte(builder.String()), 0644)
}
