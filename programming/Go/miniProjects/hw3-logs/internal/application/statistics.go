package application

import (
	"log/slog"
	"math"
	"path/filepath"
	"sort"
	"time"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain"
)

func (a *Analyzer) calculateStatistics(entries []domain.LogEntry, filenames []string) *domain.Statistics {
	if len(entries) == 0 {
		slog.Warn("No data available for analysis")
		return &domain.Statistics{
			Files:               filenames,
			TotalRequestsCount:  0,
			ResponseSizeInBytes: domain.ResponseSizeStats{},
			Resources:           []domain.ResourceStats{},
			ResponseCodes:       []domain.ResponseCodeStats{},
			RequestsPerDate:     []domain.DateStats{},
			UniqueProtocols:     []string{},
		}
	}

	var displayFiles []string
	for _, file := range filenames {
		displayFiles = append(displayFiles, filepath.Base(file))
	}

	stats := &domain.Statistics{
		Files:              displayFiles,
		TotalRequestsCount: len(entries),
	}

	a.calculateResponseSizes(stats, entries)
	a.calculateResources(stats, entries)
	a.calculateResponseCodes(stats, entries)
	a.calculateRequestsPerDate(stats, entries)
	a.calculateUniqueProtocols(stats, entries)

	return stats
}

func (a *Analyzer) calculateResponseSizes(stats *domain.Statistics, entries []domain.LogEntry) {
	var sizes []int64
	var sum int64
	var max int64

	for _, entry := range entries {
		sizes = append(sizes, entry.BodyBytesSent)
		sum += entry.BodyBytesSent
		if entry.BodyBytesSent > max {
			max = entry.BodyBytesSent
		}
	}

	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] < sizes[j]
	})

	// Calculation of 95 percentile
	p95Index := int(float64(len(sizes)) * 0.95)
	var p95 float64
	if p95Index < len(sizes) {
		p95 = float64(sizes[p95Index])
	}

	stats.ResponseSizeInBytes = domain.ResponseSizeStats{
		Average: roundFloat(float64(sum)/float64(len(entries)), 2),
		Max:     float64(max),
		P95:     roundFloat(p95, 2),
	}
}

func (a *Analyzer) calculateResources(stats *domain.Statistics, entries []domain.LogEntry) {
	resourceCounts := make(map[string]int)

	for _, entry := range entries {
		if entry.Resource != "" {
			resourceCounts[entry.Resource]++
		}
	}

	type resourceCount struct {
		resource string
		count    int
	}

	var resources []resourceCount
	for resource, count := range resourceCounts {
		resources = append(resources, resourceCount{resource, count})
	}

	sort.Slice(resources, func(i, j int) bool {
		return resources[i].count > resources[j].count
	})

	// Top-10
	topN := min(10, len(resources))
	for i := 0; i < topN; i++ {
		stats.Resources = append(stats.Resources, domain.ResourceStats{
			Resource:           resources[i].resource,
			TotalRequestsCount: resources[i].count,
		})
	}
}

func (a *Analyzer) calculateResponseCodes(stats *domain.Statistics, entries []domain.LogEntry) {
	codeCounts := make(map[int]int)

	for _, entry := range entries {
		codeCounts[entry.Status]++
	}

	for code, count := range codeCounts {
		stats.ResponseCodes = append(stats.ResponseCodes, domain.ResponseCodeStats{
			Code:                code,
			TotalResponsesCount: count,
		})
	}

	sort.Slice(stats.ResponseCodes, func(i, j int) bool {
		return stats.ResponseCodes[i].TotalResponsesCount > stats.ResponseCodes[j].TotalResponsesCount
	})
}

func (a *Analyzer) calculateRequestsPerDate(stats *domain.Statistics, entries []domain.LogEntry) {
	dateCounts := make(map[string]int)
	total := len(entries)

	for _, entry := range entries {
		date := entry.TimeLocal.Format("2006-01-02")
		dateCounts[date]++
	}

	for date, count := range dateCounts {
		t, _ := time.Parse("2006-01-02", date)
		weekday := getRussianWeekday(t.Weekday())

		percentage := roundFloat(float64(count)/float64(total)*100, 2)

		stats.RequestsPerDate = append(stats.RequestsPerDate, domain.DateStats{
			Date:                    date,
			Weekday:                 weekday,
			TotalRequestsCount:      count,
			TotalRequestsPercentage: percentage,
		})
	}

	sort.Slice(stats.RequestsPerDate, func(i, j int) bool {
		return stats.RequestsPerDate[i].Date < stats.RequestsPerDate[j].Date
	})
}

func (a *Analyzer) calculateUniqueProtocols(stats *domain.Statistics, entries []domain.LogEntry) {
	protocols := make(map[string]bool)

	for _, entry := range entries {
		if entry.Protocol != "" {
			protocols[entry.Protocol] = true
		}
	}

	for protocol := range protocols {
		stats.UniqueProtocols = append(stats.UniqueProtocols, protocol)
	}

	sort.Strings(stats.UniqueProtocols)
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getRussianWeekday(weekday time.Weekday) string {
	weekdays := map[time.Weekday]string{
		time.Monday:    "Monday",
		time.Tuesday:   "Tuesday",
		time.Wednesday: "Wednesday",
		time.Thursday:  "Thursday",
		time.Friday:    "Friday",
		time.Saturday:  "Saturday",
		time.Sunday:    "Sunday",
	}
	return weekdays[weekday]
}
