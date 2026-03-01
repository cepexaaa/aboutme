package reporter

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain"
)

type AsciiDocReporter struct{}

func (r *AsciiDocReporter) GenerateReport(stats *domain.Statistics, outputPath string) error {
	var builder strings.Builder

	builder.WriteString("= Log analysis NGINX\n\n")
	builder.WriteString(fmt.Sprintf("_Generated: %s_\n\n", time.Now().Format("02.01.2006 15:04:05")))

	builder.WriteString("== General information\n\n")
	builder.WriteString("[cols=\"1,1\"]\n")
	builder.WriteString("|===\n")
	builder.WriteString("| Metric | Value\n")
	builder.WriteString(fmt.Sprintf("| Files | `%s`\n", strings.Join(stats.Files, ", ")))
	builder.WriteString(fmt.Sprintf("| Number of requests | %d\n", stats.TotalRequestsCount))
	builder.WriteString(fmt.Sprintf("| Average response size | %.2fb\n", stats.ResponseSizeInBytes.Average))
	builder.WriteString(fmt.Sprintf("| Maximum response size | %.1fb\n", stats.ResponseSizeInBytes.Max))
	builder.WriteString(fmt.Sprintf("| 95p response size | %.2fb\n", stats.ResponseSizeInBytes.P95))
	builder.WriteString("|===\n\n")

	if len(stats.Resources) > 0 {
		builder.WriteString("== Requested resources\n\n")
		builder.WriteString("[cols=\"1,1\"]\n")
		builder.WriteString("|===\n")
		builder.WriteString("| Resource | Count\n")
		for _, resource := range stats.Resources {
			builder.WriteString(fmt.Sprintf("| `%s` | %d\n", resource.Resource, resource.TotalRequestsCount))
		}
		builder.WriteString("|===\n\n")
	}

	if len(stats.ResponseCodes) > 0 {
		builder.WriteString("== Response codes\n\n")
		builder.WriteString("[cols=\"1,1\"]\n")
		builder.WriteString("|===\n")
		builder.WriteString("| Code | Count\n")
		for _, code := range stats.ResponseCodes {
			builder.WriteString(fmt.Sprintf("| %d | %d\n", code.Code, code.TotalResponsesCount))
		}
		builder.WriteString("|===\n\n")
	}

	if len(stats.RequestsPerDate) > 0 {
		builder.WriteString("== Requests Distribution by Date\n\n")
		builder.WriteString("[cols=\"2,2,2,2\", options=\"header\"]\n")
		builder.WriteString("|===\n")
		builder.WriteString("| Date | Weekday | Requests | Percentage\n")
		for _, date := range stats.RequestsPerDate {
			builder.WriteString(fmt.Sprintf("| %s\n| %s\n| %d\n| %.2f%%\n",
				date.Date, date.Weekday, date.TotalRequestsCount, date.TotalRequestsPercentage))
		}
		builder.WriteString("|===\n\n")
	}

	if len(stats.UniqueProtocols) > 0 {
		builder.WriteString("== Unique Protocols\n\n")
		builder.WriteString("[cols=\"1\"]\n")
		builder.WriteString("|===\n")
		for _, protocol := range stats.UniqueProtocols {
			builder.WriteString(fmt.Sprintf("| %s\n", protocol))
		}
		builder.WriteString("|===\n\n")
	}

	return os.WriteFile(outputPath, []byte(builder.String()), 0644)
}
