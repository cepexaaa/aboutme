package infrastructure

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLogger_Info(t *testing.T) {
	logger := NewLogger()
	var buf bytes.Buffer

	logger.SetOutput(&buf)

	logger.Info("Test info message")

	output := buf.String()

	if !strings.Contains(output, "[INFO]") {
		t.Errorf("Output doesn't contain [INFO]: %s", output)
	}

	if !strings.Contains(output, "Test info message") {
		t.Errorf("Output doesn't contain message: %s", output)
	}

	if !strings.Contains(output, time.Now().Format("2006-01-02")) {
		t.Errorf("Output doesn't contain date: %s", output)
	}
}

func TestLogger_Warn(t *testing.T) {
	logger := NewLogger()
	var buf bytes.Buffer

	logger.SetOutput(&buf)
	logger.Warn("Test warning message")

	output := buf.String()

	if !strings.Contains(output, "[WARN]") {
		t.Errorf("Output doesn't contain [WARN]: %s", output)
	}

	if !strings.Contains(output, "Test warning message") {
		t.Errorf("Output doesn't contain message: %s", output)
	}
}

func TestLogger_Error(t *testing.T) {
	logger := NewLogger()
	var buf bytes.Buffer

	logger.SetOutput(&buf)
	logger.Error("Test error message")

	output := buf.String()

	if !strings.Contains(output, "[ERROR]") {
		t.Errorf("Output doesn't contain [ERROR]: %s", output)
	}

	if !strings.Contains(output, "Test error message") {
		t.Errorf("Output doesn't contain message: %s", output)
	}
}

func TestLogger_Debug(t *testing.T) {
	logger := NewLogger()
	var buf bytes.Buffer

	logger.SetOutput(&buf)
	logger.Debug("Test debug message")

	output := buf.String()

	if !strings.Contains(output, "[DEBUG]") {
		t.Errorf("Output doesn't contain [DEBUG]: %s", output)
	}

	if !strings.Contains(output, "Test debug message") {
		t.Errorf("Output doesn't contain message: %s", output)
	}
}

func TestLogger_ConcurrentLogging(t *testing.T) {
	logger := NewLogger()
	var buf bytes.Buffer
	var mu sync.Mutex

	safeBuf := &safeBuffer{buf: &buf, mu: &mu}
	logger.SetOutput(safeBuf)

	var wg sync.WaitGroup
	numGoroutines := 10

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			logger.Info("Message from goroutine")
		}(i)
	}

	wg.Wait()

	output := buf.String()
	lines := strings.Count(output, "\n")

	if lines != numGoroutines {
		t.Errorf("Expected %d log lines, got %d", numGoroutines, lines)
	}

	count := strings.Count(output, "Message from goroutine")
	if count != numGoroutines {
		t.Errorf("Expected %d messages, got %d", numGoroutines, count)
	}
}

func TestNewLogger(t *testing.T) {
	logger := NewLogger()

	if logger == nil {
		t.Fatal("NewLogger() returned nil")
	}

	var buf bytes.Buffer
	logger.SetOutput(&buf)
	logger.Info("Test")

	if buf.String() == "" {
		t.Error("Logger doesn't produce output")
	}
}

type safeBuffer struct {
	buf *bytes.Buffer
	mu  *sync.Mutex
}

func (sb *safeBuffer) Write(p []byte) (n int, err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.Write(p)
}

func TestLogger_MultipleMessages(t *testing.T) {
	logger := NewLogger()
	var buf bytes.Buffer

	logger.SetOutput(&buf)

	messages := []struct {
		level string
		fn    func(string)
		msg   string
	}{
		{"INFO", logger.Info, "First message"},
		{"WARN", logger.Warn, "Second message"},
		{"ERROR", logger.Error, "Third message"},
		{"DEBUG", logger.Debug, "Fourth message"},
	}

	for _, m := range messages {
		m.fn(m.msg)
	}

	output := buf.String()
	lines := strings.Count(output, "\n")

	if lines != len(messages) {
		t.Errorf("Expected %d lines, got %d", len(messages), lines)
	}

	for _, m := range messages {
		if !strings.Contains(output, m.level) {
			t.Errorf("Output doesn't contain level %s", m.level)
		}
		if !strings.Contains(output, m.msg) {
			t.Errorf("Output doesn't contain message: %s", m.msg)
		}
	}
}

func TestLogger_EmptyMessage(t *testing.T) {
	logger := NewLogger()
	var buf bytes.Buffer

	logger.SetOutput(&buf)
	logger.Info("")

	output := buf.String()

	if !strings.Contains(output, "[INFO]") {
		t.Errorf("Output doesn't contain [INFO] for empty message: %s", output)
	}
}

func TestLogger_SpecialCharacters(t *testing.T) {
	logger := NewLogger()
	var buf bytes.Buffer

	logger.SetOutput(&buf)

	testCases := []string{
		"Message with spaces",
		"Message\nwith\nnewlines",
		"Message\twith\ttabs",
		"Message with Unicode: 🚀",
		"Message with quotes: \"test\"",
		"Message with backslashes: \\test\\",
	}

	for _, msg := range testCases {
		buf.Reset()
		logger.Info(msg)

		output := buf.String()
		if !strings.Contains(output, msg) {
			t.Errorf("Output doesn't contain original message: %s", output)
		}
	}
}
