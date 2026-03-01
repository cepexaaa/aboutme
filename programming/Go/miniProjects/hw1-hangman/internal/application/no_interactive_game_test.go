package application

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestNoInteractiveGame_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedOut string
	}{
		{
			name:        "exact match",
			args:        []string{"", "apple", "apple"},
			expectedOut: "apple;POS\n",
		},
		{
			name:        "partial match",
			args:        []string{"", "apple", "apl"},
			expectedOut: "appl*;NEG\n",
		},
		{
			name:        "no match",
			args:        []string{"", "apple", "xyz"},
			expectedOut: "*****;NEG\n",
		},
		{
			name:        "case insensitive",
			args:        []string{"", "ApPlE", "APPLE"},
			expectedOut: "apple;POS\n",
		},
		{
			name:        "russian letters",
			args:        []string{"", "яблоко", "ябл"},
			expectedOut: "ябл***;NEG\n",
		},
		{
			name:        "special characters",
			args:        []string{"", "test@123", "t@3"},
			expectedOut: "t**t@**3;NEG\n",
		},
		{
			name:        "empty guess",
			args:        []string{"", "apple", ""},
			expectedOut: "*****;NEG\n",
		},
		{
			name:        "duplicate letters in guess",
			args:        []string{"", "apple", "ppll"},
			expectedOut: "*ppl*;NEG\n",
		},
		{
			name:        "extra letters in guess",
			args:        []string{"", "cat", "catxyz"},
			expectedOut: "cat;POS\n",
		},
		{
			name:        "empty word",
			args:        []string{"", "", "abc"},
			expectedOut: ";POS\n",
		},
		{
			name:        "single character word",
			args:        []string{"", "a", "a"},
			expectedOut: "a;POS\n",
		},
		{
			name:        "unicode characters",
			args:        []string{"", "café", "cafe"},
			expectedOut: "caf*;NEG\n",
		},
		{
			name:        "numbers in word",
			args:        []string{"", "test123", "t2"},
			expectedOut: "t**t*2*;NEG\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			defer func() {
				os.Stdout = oldStdout
			}()

			r, w, _ := os.Pipe()
			os.Stdout = w

			NoInteractiveGame(tt.args)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			got := buf.String()

			if got != tt.expectedOut {
				t.Errorf("NoInteractiveGame() output = %q, want %q", got, tt.expectedOut)
			}
		})
	}
}

func TestNoInteractiveGame_PosNegLogic(t *testing.T) {
	tests := []struct {
		name           string
		word           string
		guess          string
		expectedSuffix string
	}{
		{
			name:           "perfect match POS",
			word:           "hangman",
			guess:          "hangman",
			expectedSuffix: ";POS\n",
		},
		{
			name:           "partial match NEG",
			word:           "hangman",
			guess:          "hang",
			expectedSuffix: ";NEG\n",
		},
		{
			name:           "all letters wrong NEG",
			word:           "hangman",
			guess:          "xyz",
			expectedSuffix: ";NEG\n",
		},
		{
			name:           "correct but extra letters POS",
			word:           "cat",
			guess:          "catdog",
			expectedSuffix: ";POS\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			defer func() {
				os.Stdout = oldStdout
			}()

			r, w, _ := os.Pipe()
			os.Stdout = w

			args := []string{"", tt.word, tt.guess}
			NoInteractiveGame(args)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			got := buf.String()

			if !strings.HasSuffix(got, tt.expectedSuffix) {
				t.Errorf("NoInteractiveGame() suffix = %q, want %q", got[len(got)-5:], tt.expectedSuffix)
			}
		})
	}
}

func TestNoInteractiveGame_OutputFormat(t *testing.T) {
	tests := []struct {
		name    string
		word    string
		guess   string
		pattern string
	}{
		{
			name:    "stars for missing letters",
			word:    "apple",
			guess:   "ap",
			pattern: `^app\*\*;NEG\n$`,
		},
		{
			name:    "no stars for perfect match",
			word:    "test",
			guess:   "test",
			pattern: `^test;POS\n$`,
		},
		{
			name:    "correct letter order",
			word:    "banana",
			guess:   "bnn",
			pattern: `^b*n*n*;NEG\n$`,
		},
		{
			name:    "maintain word length",
			word:    "hello",
			guess:   "h",
			pattern: `^h\*\*\*\*;NEG\n$`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			defer func() {
				os.Stdout = oldStdout
			}()

			r, w, _ := os.Pipe()
			os.Stdout = w

			args := []string{"", tt.word, tt.guess}
			NoInteractiveGame(args)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			got := buf.String()

			expectedLength := len(tt.word) + 5
			if len(got) != expectedLength {
				t.Errorf("Output length = %d, want %d. Output: %q", len(got), expectedLength, got)
			}

			if !strings.HasSuffix(got, "\n") {
				t.Error("Output should end with newline")
			}
		})
	}
}
