package rawmode

import (
	"bytes"
	"fmt"
	"os"
	"syscall"
	"testing"
)

func TestTermiosStructure(t *testing.T) {
	tios := &termios{
		Iflag:  0,
		Oflag:  0,
		Cflag:  0,
		Lflag:  syscall.ICANON | syscall.ECHO,
		Ispeed: 9600,
		Ospeed: 9600,
	}

	if tios.Lflag&syscall.ICANON == 0 {
		t.Error("ICANON flag should be set initially")
	}
	if tios.Lflag&syscall.ECHO == 0 {
		t.Error("ECHO flag should be set initially")
	}
}

func TestClearScreen(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	r, w, _ := os.Pipe()
	os.Stdout = w

	ClearScreen()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expected := "\033[2J\033[H"
	if output != expected {
		t.Errorf("ClearScreen() output = %q, expected %q", output, expected)
	}
}

func TestMoveCursor(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	r, w, _ := os.Pipe()
	os.Stdout = w

	MoveCursor(10, 5)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expected := "\033[5;10H"
	if output != expected {
		t.Errorf("MoveCursor() output = %q, expected %q", output, expected)
	}
}

func TestHideShowCursor(t *testing.T) {
	tests := []struct {
		name     string
		function func()
		expected string
	}{
		{
			name:     "hide cursor",
			function: hideCursor,
			expected: "\033[?25l",
		},
		{
			name:     "show cursor",
			function: showCursor,
			expected: "\033[?25h",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			defer func() { os.Stdout = oldStdout }()

			r, w, _ := os.Pipe()
			os.Stdout = w

			tt.function()

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			if output != tt.expected {
				t.Errorf("%s output = %q, expected %q", tt.name, output, tt.expected)
			}
		})
	}
}

func TestRawMode_StopRawMode(t *testing.T) {
	tests := []struct {
		name            string
		initialTerminal *termios
		expectRestore   bool
	}{
		{
			name:            "with terminal settings",
			initialTerminal: &termios{Lflag: syscall.ICANON},
			expectRestore:   true,
		},
		{
			name:            "with nil terminal",
			initialTerminal: nil,
			expectRestore:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			defer func() { os.Stdout = oldStdout }()

			r, w, _ := os.Pipe()
			os.Stdout = w

			rm := &RawMode{terminal: tt.initialTerminal}
			rm.StopRawMode()

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			expectedOutput := "\x1b[?1049l\x1b[?25h"
			if output != expectedOutput {
				t.Errorf("StopRawMode() output = %q, expected %q", output, expectedOutput)
			}
		})
	}
}

func TestRawMode_Integration(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	r, w, _ := os.Pipe()
	os.Stdout = w

	rm := &RawMode{}

	rm.RunRawMode()

	rm.StopRawMode()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expectedOutput := "\x1b[?1049h\x1b[?25l\x1b[?1049l\x1b[?25h"
	if output != expectedOutput {
		t.Errorf("Full cycle output = %q, expected %q", output, expectedOutput)
	}
}

func TestRawMode_ErrorScenarios(t *testing.T) {
	t.Run("stop without run", func(t *testing.T) {
		rm := &RawMode{}
		rm.StopRawMode()

		if rm.terminal != nil {
			t.Error("StopRawMode should handle nil terminal gracefully")
		}
	})

	t.Run("multiple stop calls", func(t *testing.T) {
		rm := &RawMode{}
		rm.RunRawMode()
		rm.StopRawMode()
		rm.StopRawMode()
		rm.StopRawMode()
	})
}

func TestANSISequences(t *testing.T) {
	testCases := []struct {
		name     string
		call     func()
		expected string
	}{
		{
			name:     "ClearScreen sequence",
			call:     ClearScreen,
			expected: "\033[2J\033[H",
		},
		{
			name:     "MoveCursor sequence",
			call:     func() { MoveCursor(15, 20) },
			expected: "\033[20;15H",
		},
		{
			name:     "hideCursor sequence",
			call:     hideCursor,
			expected: "\033[?25l",
		},
		{
			name:     "showCursor sequence",
			call:     showCursor,
			expected: "\033[?25h",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oldStdout := os.Stdout
			defer func() { os.Stdout = oldStdout }()

			r, w, _ := os.Pipe()
			os.Stdout = w

			tc.call()

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)

			if buf.String() != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, buf.String())
			}
		})
	}
}

func captureOutput(f func()) string {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)

	return buf.String()
}

func TestCaptureOutputHelper(t *testing.T) {
	expected := "test output"

	output := captureOutput(func() {
		fmt.Print(expected)
	})

	if output != expected {
		t.Errorf("captureOutput() = %q, expected %q", output, expected)
	}
}
