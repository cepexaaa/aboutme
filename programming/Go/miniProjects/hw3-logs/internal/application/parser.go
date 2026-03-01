package application

import (
	"bufio"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain"
)

var (
	logPattern     = regexp.MustCompile(`^(\S+) - (\S+) \[([^\]]+)\] "([^"]*)" (\d+) (\d+) "([^"]*)" "([^"]*)"$`)
	requestPattern = regexp.MustCompile(`^(\S+)\s+(\S+)\s+(\S+)$`)
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseLine(line string) (*domain.LogEntry, error) {
	matches := logPattern.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("invalid format of the log string")
	}

	var timeLocal time.Time

	timeLocal, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[3])
	if err != nil {
		timeLocal, err = time.Parse("2/Jan/2006:15:04:05 -0700", matches[3])
		if err != nil {
			return nil, fmt.Errorf("time parsing error: %v", err)
		}
	}

	status, err := strconv.Atoi(matches[5])
	if err != nil {
		return nil, fmt.Errorf("status parsing error: %v", err)
	}

	bodyBytesSent, err := strconv.ParseInt(matches[6], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("response size parsing error: %v", err)
	}

	resource, protocol := parseRequest(matches[4])

	return &domain.LogEntry{
		RemoteAddr:    matches[1],
		RemoteUser:    matches[2],
		TimeLocal:     timeLocal,
		Request:       matches[4],
		Status:        status,
		BodyBytesSent: bodyBytesSent,
		HTTPReferer:   matches[7],
		HTTPUserAgent: matches[8],
		Resource:      resource,
		Protocol:      protocol,
	}, nil
}

func parseRequest(request string) (string, string) {
	matches := requestPattern.FindStringSubmatch(request)
	if matches == nil {
		return "", ""
	}

	return matches[2], matches[3]
}

func (p *Parser) ParseFile(scanner *bufio.Scanner) ([]domain.LogEntry, error) {
	var entries []domain.LogEntry
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		entry, err := p.ParseLine(line)
		if err != nil {
			slog.Warn("Skipping an invalid line",
				"line", lineNumber,
				"error", err,
				"content", line[:min(50, len(line))]+"...")
			continue
		}

		entries = append(entries, *entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("file reading error: %v", err)
	}

	slog.Info("Successfully parsed records", "count", len(entries))
	return entries, nil
}
