package reporter

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain"
)

func TestJSONReporter_GenerateReport(t *testing.T) {
	r := &JSONReporter{}

	stats := &domain.Statistics{
		Files:              []string{"file1.log", "file2.log"},
		TotalRequestsCount: 100,
		ResponseSizeInBytes: domain.ResponseSizeStats{
			Average: 150.5,
			Max:     1000,
			P95:     250.75,
		},
		Resources: []domain.ResourceStats{
			{Resource: "/api/v1/users", TotalRequestsCount: 50},
			{Resource: "/api/v1/products", TotalRequestsCount: 30},
		},
		ResponseCodes: []domain.ResponseCodeStats{
			{Code: 200, TotalResponsesCount: 80},
			{Code: 404, TotalResponsesCount: 20},
		},
	}

	tmpFile := "/tmp/test_report.json"
	defer os.Remove(tmpFile)

	err := r.GenerateReport(stats, tmpFile)
	require.NoError(t, err)

	_, err = os.Stat(tmpFile)
	assert.NoError(t, err)

	content, err := os.ReadFile(tmpFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), `"totalRequestsCount": 100`)
	assert.Contains(t, string(content), `"resource": "/api/v1/users"`)
}

func TestMarkdownReporter_GenerateReport(t *testing.T) {
	r := &MarkdownReporter{}

	stats := &domain.Statistics{
		Files:              []string{"access.log"},
		TotalRequestsCount: 50,
		ResponseSizeInBytes: domain.ResponseSizeStats{
			Average: 200.0,
			Max:     1500,
			P95:     400.0,
		},
		RequestsPerDate: []domain.DateStats{
			{
				Date:                    "2024-03-15",
				Weekday:                 "Friday",
				TotalRequestsCount:      25,
				TotalRequestsPercentage: 50.0,
			},
		},
		UniqueProtocols: []string{"HTTP/1.1", "HTTP/2.0"},
	}

	tmpFile := "/tmp/test_report.md"
	defer os.Remove(tmpFile)

	err := r.GenerateReport(stats, tmpFile)
	require.NoError(t, err)

	content, err := os.ReadFile(tmpFile)
	require.NoError(t, err)

	contentStr := string(content)
	assert.Contains(t, contentStr, "# Log analysis NGINX")
	assert.Contains(t, contentStr, "Number of requests")
	assert.Contains(t, contentStr, "HTTP/1.1")
	assert.Contains(t, contentStr, "2024-03-15")
}

func TestAsciiDocReporter_GenerateReport(t *testing.T) {
	r := &AsciiDocReporter{}

	stats := &domain.Statistics{
		Files:              []string{"test.log"},
		TotalRequestsCount: 10,
		ResponseSizeInBytes: domain.ResponseSizeStats{
			Average: 100.0,
			Max:     500,
			P95:     200.0,
		},
		UniqueProtocols: []string{"HTTP/1.1"},
	}

	tmpFile := "/tmp/test_report.adoc"
	defer os.Remove(tmpFile)

	err := r.GenerateReport(stats, tmpFile)
	require.NoError(t, err)

	content, err := os.ReadFile(tmpFile)
	require.NoError(t, err)

	contentStr := string(content)
	assert.Contains(t, contentStr, "= Log analysis NGINX")
	assert.Contains(t, contentStr, "General information")
	assert.Contains(t, contentStr, "Unique Protocols")
}
