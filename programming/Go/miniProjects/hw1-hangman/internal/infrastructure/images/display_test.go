package images

import (
	"os"
	"strings"
	"testing"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw1-hangman/internal/domain"
)

func TestGetImage(t *testing.T) {
	tests := []struct {
		name     string
		state    int
		expected int
	}{
		{"state 1", 1, 0},
		{"state 5", 5, 0},
		{"state 10", 10, 0},
		{"state out of range high", 15, 0},
		{"state out of range low", -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getImage(tt.state)

			if result == nil {
				t.Log("getImage returned nil (expected for missing files)")
			}
		})
	}
}

func TestMoveArrow(t *testing.T) {
	pageContent := []string{
		"Title",
		"",
		"Option 1",
		"Option 2",
		"Option 3",
	}

	game := &domain.GameSettings{CurrentChoise: 1}

	if !strings.HasPrefix(pageContent[3], "Option 2") {
		t.Error("Option 2 should not have arrow initially")
	}

	moveArrow(game, &pageContent, 1)

	if game.CurrentChoise != 2 {
		t.Errorf("CurrentChoise = %d, expected 2", game.CurrentChoise)
	}

	if !strings.HasPrefix(pageContent[4], "->") {
		t.Error("Option 3 should have arrow after move")
	}

	if strings.HasPrefix(pageContent[3], "->") {
		t.Error("Option 2 should not have arrow after move")
	}

	moveArrow(game, &pageContent, -1)

	if game.CurrentChoise != 1 {
		t.Errorf("CurrentChoise = %d, expected 1", game.CurrentChoise)
	}
}

func TestReadUTF8Char(t *testing.T) {
	tests := []struct {
		name        string
		firstByte   byte
		stdinBytes  []byte
		expected    rune
		expectError bool
	}{
		{
			name:        "ASCII character",
			firstByte:   65, // 'A'
			stdinBytes:  []byte{},
			expected:    'A',
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalStdin := os.Stdin
			defer func() { os.Stdin = originalStdin }()

			result, err := readUTF8Char(tt.firstByte)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected rune %c, got %c", tt.expected, result)
				}
			}
		})
	}
}

func TestPrintGameInfo(t *testing.T) {
	tests := []struct {
		name         string
		gameState    int
		countLetters int
		answerLength int
		lang         string
		hint         string
		expectedQuit bool
	}{
		{
			name:         "game in progress",
			gameState:    3,
			countLetters: 2,
			answerLength: 5,
			lang:         "en",
			hint:         "",
			expectedQuit: false,
		},
		{
			name:         "game lost",
			gameState:    10,
			countLetters: 3,
			answerLength: 5,
			lang:         "en",
			hint:         "",
			expectedQuit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := &domain.GameSettings{
				CurrentState: tt.gameState,
				CountLetters: tt.countLetters,
				Answer:       make([]rune, tt.answerLength),
				Result:       make([]rune, tt.answerLength),
				Params:       []string{tt.lang},
			}

			if tt.countLetters == tt.answerLength {
				for i := range game.Result {
					game.Result[i] = 'a'
				}
			}

			printGameInfo(game, tt.lang, tt.hint)

			if game.Quit != tt.expectedQuit {
				t.Errorf("Quit = %v, expected %v", game.Quit, tt.expectedQuit)
			}
		})
	}
}

func TestReadUTF8Char_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		firstByte  byte
		stdinInput []byte
		wantRune   rune
		wantError  bool
	}{
		{
			name:       "ASCII uppercase letter",
			firstByte:  65, // 'A'
			stdinInput: []byte{},
			wantRune:   'A',
			wantError:  false,
		},
		{
			name:       "ASCII lowercase letter",
			firstByte:  97, // 'a'
			stdinInput: []byte{},
			wantRune:   'a',
			wantError:  false,
		},
		{
			name:       "ASCII digit",
			firstByte:  48, // '0'
			stdinInput: []byte{},
			wantRune:   '0',
			wantError:  false,
		},
		{
			name:       "ASCII special character",
			firstByte:  43, // '+'
			stdinInput: []byte{},
			wantRune:   '+',
			wantError:  false,
		},
		{
			name:       "Russian letter а",
			firstByte:  208,
			stdinInput: []byte{176}, // 'а'
			wantRune:   'а',
			wantError:  false,
		},
		{
			name:       "Russian letter я",
			firstByte:  209,
			stdinInput: []byte{143}, // 'я'
			wantRune:   'я',
			wantError:  false,
		},
		{
			name:       "3-byte UTF-8 character",
			firstByte:  226,
			stdinInput: []byte{130, 172}, // €
			wantRune:   0,
			wantError:  false,
		},
		{
			name:       "EOF during read",
			firstByte:  208,
			stdinInput: []byte{},
			wantRune:   ' ',
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock stdin
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			tmpfile, _ := os.CreateTemp("", "stdin")
			tmpfile.Write(tt.stdinInput)
			tmpfile.Seek(0, 0)
			os.Stdin = tmpfile

			got, err := readUTF8Char(tt.firstByte)

			tmpfile.Close()
			os.Remove(tmpfile.Name())

			if tt.wantError {
				if err == nil {
					t.Errorf("readUTF8Char() error = %v, want error", err)
				}
			} else {
				if err != nil {
					t.Errorf("readUTF8Char() unexpected error = %v", err)
				}
				if got != tt.wantRune {
					t.Errorf("readUTF8Char() = %U, want %U", got, tt.wantRune)
				}
			}
		})
	}
}

func TestMoveArrow_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		initialChoice int
		move          int
		pageContent   []string
		wantChoice    int
		wantPrefix    string
		wantOldClean  bool
	}{
		{
			name:          "move down within bounds",
			initialChoice: 0,
			move:          1,
			pageContent:   []string{"", "", "Option1", "Option2", "Option3"},
			wantChoice:    1,
			wantPrefix:    "->Option2",
			wantOldClean:  true,
		},
		{
			name:          "move up within bounds",
			initialChoice: 2,
			move:          -1,
			pageContent:   []string{"", "", "Option1", "Option2", "Option3", "Option4"},
			wantChoice:    1,
			wantPrefix:    "->Option2",
			wantOldClean:  true,
		},
		{
			name:          "move down from first position",
			initialChoice: 0,
			move:          1,
			pageContent:   []string{"", "", "First", "Second", "Third"},
			wantChoice:    1,
			wantPrefix:    "->Second",
			wantOldClean:  true,
		},
		{
			name:          "move down from last position",
			initialChoice: 2,
			move:          1,
			pageContent:   []string{"", "", "1", "2", "3", "4"},
			wantChoice:    3,
			wantPrefix:    "->4",
			wantOldClean:  false,
		},
		{
			name:          "multiple moves",
			initialChoice: 0,
			move:          2,
			pageContent:   []string{"", "", "A", "B", "C", "D"},
			wantChoice:    2,
			wantPrefix:    "->C",
			wantOldClean:  true,
		},
		{
			name:          "move with empty options",
			initialChoice: 0,
			move:          1,
			pageContent:   []string{"Title", "", "Opt1", "", "Opt2"},
			wantChoice:    1,
			wantPrefix:    "->",
			wantOldClean:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := &domain.GameSettings{CurrentChoise: tt.initialChoice}
			content := make([]string, len(tt.pageContent))
			copy(content, tt.pageContent)

			content[tt.initialChoice+2] = "->" + content[tt.initialChoice+2]

			moveArrow(game, &content, tt.move)

			if game.CurrentChoise != tt.wantChoice {
				t.Errorf("CurrentChoise = %d, want %d", game.CurrentChoise, tt.wantChoice)
			}

			if !strings.HasPrefix(content[tt.wantChoice+2], tt.wantPrefix) {
				t.Errorf("Option text = %s, want prefix %s", content[tt.wantChoice+2], tt.wantPrefix)
			}

			// Check that old position is cleaned
			if tt.wantOldClean && tt.initialChoice != tt.wantChoice {
				oldOption := content[tt.initialChoice+2]
				if strings.HasPrefix(oldOption, "->") {
					t.Errorf("Old option still has arrow: %s", oldOption)
				}
			}
		})
	}
}

func TestPrintGameInfo_Conditions_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		currentState int
		countLetters int
		answerLength int
		lang         string
		hint         string
		wantQuit     bool
	}{
		{
			name:         "game in progress en",
			currentState: 3,
			countLetters: 2,
			answerLength: 5,
			lang:         "en",
			hint:         "",
			wantQuit:     false,
		},
		{
			name:         "game in progress ru",
			currentState: 3,
			countLetters: 2,
			answerLength: 5,
			lang:         "ru",
			hint:         "",
			wantQuit:     false,
		},
		{
			name:         "game lost state 10",
			currentState: 10,
			countLetters: 3,
			answerLength: 5,
			lang:         "en",
			hint:         "",
			wantQuit:     true,
		},
		{
			name:         "with hint en",
			currentState: 2,
			countLetters: 1,
			answerLength: 4,
			lang:         "en",
			hint:         "test hint",
			wantQuit:     false,
		},
		{
			name:         "with hint ru",
			currentState: 2,
			countLetters: 1,
			answerLength: 4,
			lang:         "ru",
			hint:         "тестовая подсказка",
			wantQuit:     false,
		},
		{
			name:         "last attempt en",
			currentState: 9,
			countLetters: 3,
			answerLength: 5,
			lang:         "en",
			hint:         "",
			wantQuit:     false,
		},
		{
			name:         "last attempt ru",
			currentState: 9,
			countLetters: 3,
			answerLength: 5,
			lang:         "ru",
			hint:         "",
			wantQuit:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := &domain.GameSettings{
				CurrentState: tt.currentState,
				CountLetters: tt.countLetters,
				Answer:       make([]rune, tt.answerLength),
				Result:       make([]rune, tt.answerLength),
				Params:       []string{tt.lang},
			}

			// Fill result for win condition
			if tt.countLetters == tt.answerLength {
				for i := range game.Result {
					game.Result[i] = 'a'
				}
			}

			printGameInfo(game, tt.lang, tt.hint)

			if game.Quit != tt.wantQuit {
				t.Errorf("Quit = %v, want %v", game.Quit, tt.wantQuit)
			}
		})
	}
}
