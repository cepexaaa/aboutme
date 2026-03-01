package application

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/infrastructure"
)

func TestAnalyzer_CalculateStatistics(t *testing.T) {
	a := &Analyzer{}

	entries := []domain.LogEntry{
		{
			TimeLocal:     parseTime("17/May/2015:08:05:32 +0000"),
			Status:        200,
			BodyBytesSent: 100,
			Resource:      "/product/1",
			Protocol:      "HTTP/1.1",
		},
		{
			TimeLocal:     parseTime("17/May/2015:08:05:33 +0000"),
			Status:        200,
			BodyBytesSent: 200,
			Resource:      "/product/1",
			Protocol:      "HTTP/1.1",
		},
		{
			TimeLocal:     parseTime("17/May/2015:08:05:34 +0000"),
			Status:        404,
			BodyBytesSent: 50,
			Resource:      "/product/2",
			Protocol:      "HTTP/1.1",
		},
		{
			TimeLocal:     parseTime("17/May/2015:08:05:35 +0000"),
			Status:        500,
			BodyBytesSent: 300,
			Resource:      "/product/1",
			Protocol:      "HTTP/2.0",
		},
		{
			TimeLocal:     parseTime("18/May/2015:08:05:36 +0000"),
			Status:        200,
			BodyBytesSent: 150,
			Resource:      "/product/3",
			Protocol:      "HTTP/1.1",
		},
	}

	filenames := []string{"test.log"}
	stats := a.calculateStatistics(entries, filenames)

	require.NotNil(t, stats)
	assert.Equal(t, 5, stats.TotalRequestsCount)
	assert.Equal(t, filenames, stats.Files)

	assert.InDelta(t, 160.0, stats.ResponseSizeInBytes.Average, 0.01)
	assert.Equal(t, float64(300), stats.ResponseSizeInBytes.Max)
	assert.InDelta(t, 300.0, stats.ResponseSizeInBytes.P95, 0.01)

	require.Len(t, stats.Resources, 3)
	assert.Equal(t, "/product/1", stats.Resources[0].Resource)
	assert.Equal(t, 3, stats.Resources[0].TotalRequestsCount)
	assert.Equal(t, "/product/2", stats.Resources[1].Resource)
	assert.Equal(t, 1, stats.Resources[1].TotalRequestsCount)
	assert.Equal(t, "/product/3", stats.Resources[2].Resource)
	assert.Equal(t, 1, stats.Resources[2].TotalRequestsCount)

	require.Len(t, stats.ResponseCodes, 3)
	assert.Equal(t, 200, stats.ResponseCodes[0].Code)
	assert.Equal(t, 3, stats.ResponseCodes[0].TotalResponsesCount)
	assert.Equal(t, 404, stats.ResponseCodes[1].Code)
	assert.Equal(t, 1, stats.ResponseCodes[1].TotalResponsesCount)
	assert.Equal(t, 500, stats.ResponseCodes[2].Code)
	assert.Equal(t, 1, stats.ResponseCodes[2].TotalResponsesCount)

	require.Len(t, stats.RequestsPerDate, 2)

	require.Len(t, stats.UniqueProtocols, 2)
	assert.Contains(t, stats.UniqueProtocols, "HTTP/1.1")
	assert.Contains(t, stats.UniqueProtocols, "HTTP/2.0")
}

func TestAnalyzer_CalculateResponseSizes_Percentile(t *testing.T) {
	a := &Analyzer{}

	entries := make([]domain.LogEntry, 100)
	for i := 0; i < 100; i++ {
		entries[i] = domain.LogEntry{
			BodyBytesSent: int64(i + 1),
		}
	}

	filenames := []string{"test.log"}
	stats := a.calculateStatistics(entries, filenames)

	assert.InDelta(t, 95.0, stats.ResponseSizeInBytes.P95, 1.0)
}

func TestAnalyzer_CalculateStatistics_EmptyData(t *testing.T) {
	a := &Analyzer{}

	stats := a.calculateStatistics([]domain.LogEntry{}, []string{"test.log"})

	require.NotNil(t, stats)
	assert.Equal(t, 0, stats.TotalRequestsCount)
	assert.Equal(t, []string{"test.log"}, stats.Files)
	assert.Empty(t, stats.Resources)
	assert.Empty(t, stats.ResponseCodes)
	assert.Empty(t, stats.RequestsPerDate)
	assert.Empty(t, stats.UniqueProtocols)
}

func TestAnalyzer_FilterByDate(t *testing.T) {
	a := &Analyzer{
		config: &infrastructure.Config{},
	}

	entries := []domain.LogEntry{
		{TimeLocal: parseTime("01/May/2015:08:05:32 +0000")},
		{TimeLocal: parseTime("15/May/2015:08:05:32 +0000")},
		{TimeLocal: parseTime("17/May/2015:08:05:32 +0000")},
		{TimeLocal: parseTime("20/May/2015:08:05:32 +0000")},
		{TimeLocal: parseTime("25/May/2015:08:05:32 +0000")},
	}

	fromDate := parseTime("15/May/2015:00:00:00 +0000")
	a.config.From = &fromDate
	filtered := a.filterByDate(entries)
	assert.Len(t, filtered, 4)

	a.config.From = nil
	toDate := parseTime("20/May/2015:23:59:59 +0000")
	a.config.To = &toDate
	filtered = a.filterByDate(entries)
	assert.Len(t, filtered, 4)

	a.config.From = &fromDate
	a.config.To = &toDate
	filtered = a.filterByDate(entries)
	assert.Len(t, filtered, 3)

	a.config.From = nil
	a.config.To = nil
	filtered = a.filterByDate(entries)
	assert.Len(t, filtered, 5)
}

func parseTime(timeStr string) time.Time {
	t, err := time.Parse("02/Jan/2006:15:04:05 -0700", timeStr)
	if err != nil {
		panic(err)
	}
	return t
}
