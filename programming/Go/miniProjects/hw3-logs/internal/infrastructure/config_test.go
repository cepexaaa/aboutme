package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags_RequiredParameters(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd"}

	_, err := ParseFlags()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parameter --path required")
}

func TestValidateOutputExtension(t *testing.T) {
	testCases := []struct {
		format   OutputFormat
		output   string
		expected error
	}{
		{JSONFormat, "report.json", nil},
		{JSONFormat, "report.txt", assert.AnError},
		{MarkdownFormat, "report.md", nil},
		{MarkdownFormat, "report.json", assert.AnError},
		{AsciiDocFormat, "report.adoc", nil},
		{AsciiDocFormat, "report.ad", nil},
		{AsciiDocFormat, "report.txt", assert.AnError},
	}

	for _, tc := range testCases {
		err := validateOutputExtension(tc.format, tc.output)
		if tc.expected == nil {
			assert.NoError(t, err, "Format: %s, Output: %s", tc.format, tc.output)
		} else {
			assert.Error(t, err, "Format: %s, Output: %s", tc.format, tc.output)
		}
	}
}
