package domain

import (
	"time"
)

type LogEntry struct {
	RemoteAddr    string
	RemoteUser    string
	TimeLocal     time.Time
	Request       string
	Status        int
	BodyBytesSent int64
	HTTPReferer   string
	HTTPUserAgent string
	Resource      string
	Protocol      string
}

type Statistics struct {
	Files               []string            `json:"files"`
	TotalRequestsCount  int                 `json:"totalRequestsCount"`
	ResponseSizeInBytes ResponseSizeStats   `json:"responseSizeInBytes"`
	Resources           []ResourceStats     `json:"resources"`
	ResponseCodes       []ResponseCodeStats `json:"responseCodes"`
	RequestsPerDate     []DateStats         `json:"requestsPerDate,omitempty"`
	UniqueProtocols     []string            `json:"uniqueProtocols,omitempty"`
}

type ResponseSizeStats struct {
	Average float64 `json:"average"`
	Max     float64 `json:"max"`
	P95     float64 `json:"p95"`
}

type ResourceStats struct {
	Resource           string `json:"resource"`
	TotalRequestsCount int    `json:"totalRequestsCount"`
}

type ResponseCodeStats struct {
	Code                int `json:"code"`
	TotalResponsesCount int `json:"totalResponsesCount"`
}

type DateStats struct {
	Date                    string  `json:"date"`
	Weekday                 string  `json:"weekday"`
	TotalRequestsCount      int     `json:"totalRequestsCount"`
	TotalRequestsPercentage float64 `json:"totalRequestsPercentage"`
}
