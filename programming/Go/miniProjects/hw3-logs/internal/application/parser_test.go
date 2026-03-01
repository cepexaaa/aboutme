package application

import (
	"bufio"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_ParseLine_ValidLogEntry(t *testing.T) {
	p := NewParser()

	line := `93.180.71.3 - - [17/May/2015:08:05:32 +0000] "GET /downloads/product_1 HTTP/1.1" 304 0 "-" "Debian APT-HTTP/1.3 (0.8.16~exp12ubuntu10.21)"`

	entry, err := p.ParseLine(line)
	require.NoError(t, err)

	expectedTime, _ := time.Parse("02/Jan/2006:15:04:05 -0700", "17/May/2015:08:05:32 +0000")

	assert.Equal(t, "93.180.71.3", entry.RemoteAddr)
	assert.Equal(t, "-", entry.RemoteUser)
	assert.Equal(t, expectedTime, entry.TimeLocal)
	assert.Equal(t, "GET /downloads/product_1 HTTP/1.1", entry.Request)
	assert.Equal(t, 304, entry.Status)
	assert.Equal(t, int64(0), entry.BodyBytesSent)
	assert.Equal(t, "-", entry.HTTPReferer)
	assert.Equal(t, "Debian APT-HTTP/1.3 (0.8.16~exp12ubuntu10.21)", entry.HTTPUserAgent)
	assert.Equal(t, "/downloads/product_1", entry.Resource)
	assert.Equal(t, "HTTP/1.1", entry.Protocol)
}

func TestParser_ParseLine_WithSingleDigitDay(t *testing.T) {
	p := NewParser()

	line := `188.138.60.101 - - [1/May/2015:08:05:22 +0000] "GET /downloads/product_2 HTTP/1.1" 404 337 "-" "Debian APT-HTTP/1.3 (0.9.7.9)"`

	entry, err := p.ParseLine(line)
	require.NoError(t, err)

	assert.Equal(t, "188.138.60.101", entry.RemoteAddr)
	assert.Equal(t, 404, entry.Status)
	assert.Equal(t, int64(337), entry.BodyBytesSent)
	assert.Equal(t, "/downloads/product_2", entry.Resource)
}

func TestParser_ParseLine_InvalidFormat(t *testing.T) {
	p := NewParser()

	invalidLines := []string{
		"invalid log line",
		"",
		"93.180.71.3 - - [invalid] \"GET\" 200 0 \"\" \"\"",
		"93.180.71.3 - - [17/May/2015:08:05:32 +0000] \"GET\" invalid 0 \"\" \"\"",
	}

	for _, line := range invalidLines {
		_, err := p.ParseLine(line)
		assert.Error(t, err, "Expected error for line: %s", line)
	}
}

func TestParser_ParseLine_DifferentProtocols(t *testing.T) {
	p := NewParser()

	testCases := []struct {
		line             string
		expectedResource string
		expectedProtocol string
	}{
		{
			`93.180.71.3 - - [17/May/2015:08:05:32 +0000] "GET /downloads/product_1 HTTP/1.1" 304 0 "-" "Debian APT-HTTP/1.3"`,
			"/downloads/product_1",
			"HTTP/1.1",
		},
		{
			`80.91.33.133 - - [17/May/2015:09:05:15 +0000] "GET /downloads/product_1 HTTP/2.1" 304 0 "-" "Debian APT-HTTP/1.3"`,
			"/downloads/product_1",
			"HTTP/2.1",
		},
		{
			`91.121.161.213 - - [17/May/2015:09:05:58 +0000] "GET /downloads/product_2 grpc" 404 346 "-" "Debian APT-HTTP/1.3"`,
			"/downloads/product_2",
			"grpc",
		},
	}

	for _, tc := range testCases {
		entry, err := p.ParseLine(tc.line)
		require.NoError(t, err)
		assert.Equal(t, tc.expectedResource, entry.Resource)
		assert.Equal(t, tc.expectedProtocol, entry.Protocol)
	}
}

func TestParser_ParseFile(t *testing.T) {
	p := NewParser()

	logContent := `93.180.71.3 - - [17/May/2015:08:05:32 +0000] "GET /downloads/product_1 HTTP/1.1" 304 0 "-" "Debian APT-HTTP/1.3"
80.91.33.133 - - [17/May/2015:08:05:23 +0000] "GET /downloads/product_1 HTTP/1.1" 304 0 "-" "Debian APT-HTTP/1.3"
217.168.17.5 - - [17/May/2015:08:05:34 +0000] "GET /downloads/product_1 HTTP/1.1" 200 490 "-" "Debian APT-HTTP/1.3"`

	scanner := bufio.NewScanner(strings.NewReader(logContent))
	entries, err := p.ParseFile(scanner)

	require.NoError(t, err)
	assert.Len(t, entries, 3)

	assert.Equal(t, 200, entries[2].Status)
	assert.Equal(t, int64(490), entries[2].BodyBytesSent)
}
