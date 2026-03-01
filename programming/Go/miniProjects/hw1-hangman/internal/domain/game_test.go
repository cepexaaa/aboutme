package domain

import (
	"strings"
	"testing"
)

func TestGameSettings_GetHint(t *testing.T) {
	tests := []struct {
		name     string
		hints    []string
		expected string
	}{
		{
			name:     "no hints",
			hints:    []string{},
			expected: "No any hints",
		},
		{
			name:     "single hint",
			hints:    []string{"Это животное"},
			expected: "Это животное",
		},
		{
			name:     "multiple hints",
			hints:    []string{"Подсказка 1", "Подсказка 2", "Подсказка 3"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GameSettings{Hints: tt.hints}

			result := gs.GetHint()

			if len(tt.hints) == 0 {
				if result != tt.expected {
					t.Errorf("GetHint() = %v, expected %v", result, tt.expected)
				}
			} else {
				found := false
				for _, hint := range tt.hints {
					if result == hint {
						found = true
						break
					}
				}
				if !found && result != "No any hints" {
					t.Errorf("GetHint() = %v, expected one of %v", result, tt.hints)
				}
			}
		})
	}
}

func TestGameSettings_ProcessLetter(t *testing.T) {
	tests := []struct {
		name           string
		initialAnswer  string
		initialResult  []rune
		initialState   int
		initialCount   int
		inputRune      rune
		expectedResult string
		expectedState  int
		expectedCount  int
	}{
		{
			name:           "correct guess new letter",
			initialAnswer:  "apple",
			initialResult:  []rune{'_', '_', '_', '_', '_'},
			initialState:   0,
			initialCount:   0,
			inputRune:      'a',
			expectedResult: "a____",
			expectedState:  0,
			expectedCount:  1,
		},
		{
			name:           "correct guess existing letter",
			initialAnswer:  "apple",
			initialResult:  []rune{'a', '_', '_', '_', '_'},
			initialState:   0,
			initialCount:   1,
			inputRune:      'a',
			expectedResult: "a____",
			expectedState:  0,
			expectedCount:  1,
		},
		{
			name:           "incorrect guess",
			initialAnswer:  "apple",
			initialResult:  []rune{'_', '_', '_', '_', '_'},
			initialState:   0,
			initialCount:   0,
			inputRune:      'z',
			expectedResult: "_____",
			expectedState:  1,
			expectedCount:  0,
		},
		{
			name:           "correct guess multiple letters",
			initialAnswer:  "banana",
			initialResult:  []rune{'_', '_', '_', '_', '_', '_'},
			initialState:   0,
			initialCount:   0,
			inputRune:      'a',
			expectedResult: "_a_a_a",
			expectedState:  0,
			expectedCount:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GameSettings{
				Answer:       []rune(tt.initialAnswer),
				Result:       tt.initialResult,
				CurrentState: tt.initialState,
				CountLetters: tt.initialCount,
			}

			gs.ProccessLetter(tt.inputRune)

			resultStr := string(gs.Result)
			if resultStr != tt.expectedResult {
				t.Errorf("Result = %v, expected %v", resultStr, tt.expectedResult)
			}

			if gs.CurrentState != tt.expectedState {
				t.Errorf("CurrentState = %v, expected %v", gs.CurrentState, tt.expectedState)
			}

			if gs.CountLetters != tt.expectedCount {
				t.Errorf("CountLetters = %v, expected %v", gs.CountLetters, tt.expectedCount)
			}
		})
	}
}

func TestGameSettings_setResult(t *testing.T) {
	tests := []struct {
		name          string
		answer        string
		expectedLen   int
		expectedChars string
	}{
		{
			name:          "empty word",
			answer:        "",
			expectedLen:   0,
			expectedChars: "",
		},
		{
			name:          "short word",
			answer:        "cat",
			expectedLen:   3,
			expectedChars: "___",
		},
		{
			name:          "long word",
			answer:        "elephant",
			expectedLen:   8,
			expectedChars: "________",
		},
		{
			name:          "word with spaces",
			answer:        "hello world",
			expectedLen:   11,
			expectedChars: "___________",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GameSettings{
				Answer: []rune(tt.answer),
			}

			gs.setResult()

			if len(gs.Result) != tt.expectedLen {
				t.Errorf("Result length = %v, expected %v", len(gs.Result), tt.expectedLen)
			}

			resultStr := string(gs.Result)
			if resultStr != tt.expectedChars {
				t.Errorf("Result = %v, expected %v", resultStr, tt.expectedChars)
			}

			for _, char := range gs.Result {
				if char != '_' {
					t.Errorf("Found non-underscore character: %v", char)
				}
			}
		})
	}
}

func TestGameSettings_FilterWords(t *testing.T) {
	mockWords := &Words{
		En: ComplexityWords{
			Easy: []Word{
				{Name: "Apple", Hints: []string{"Fruit", "Red"}},
				{Name: "Cat", Hints: []string{"Animal", "Pet"}},
			},
			Hard: []Word{
				{Name: "Elephant", Hints: []string{"Big animal", "Gray"}},
				{Name: "Computer", Hints: []string{"Electronic device"}},
			},
		},
		Ru: ComplexityWords{
			Easy: []Word{
				{Name: "Яблоко", Hints: []string{"Фрукт", "Красное"}},
				{Name: "Кот", Hints: []string{"Животное", "Домашнее"}},
			},
			Hard: []Word{
				{Name: "Слон", Hints: []string{"Большое животное", "Серое"}},
				{Name: "Компьютер", Hints: []string{"Электронное устройство"}},
			},
		},
	}

	tests := []struct {
		name         string
		params       []string
		expectedLang string
		expectedCat  string
		expectWord   bool
		expectHints  bool
	}{
		{
			name:         "english easy",
			params:       []string{"en", "easy", "low"},
			expectedLang: "en",
			expectedCat:  "easy",
			expectWord:   true,
			expectHints:  true,
		},
		{
			name:         "russian hard",
			params:       []string{"ru", "hard", "high"},
			expectedLang: "ru",
			expectedCat:  "hard",
			expectWord:   true,
			expectHints:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GameSettings{
				Params: tt.params,
			}

			selectedWord := gs.FilterWords(mockWords)

			if tt.expectWord && selectedWord.Name == "" {
				t.Error("Expected word to be selected, but got empty")
			}

			if len(gs.Params) >= 2 {
				if tt.expectedLang != "" && gs.Params[0] != tt.expectedLang {
					t.Errorf("Language param = %v, expected %v", gs.Params[0], tt.expectedLang)
				}

				if tt.expectedCat != "" && gs.Params[1] != tt.expectedCat {
					t.Errorf("Category param = %v, expected %v", gs.Params[1], tt.expectedCat)
				}
			}

			answerStr := string(gs.Answer)
			if answerStr != strings.ToLower(selectedWord.Name) {
				t.Errorf("Answer = %v, expected lowercase of %v", answerStr, selectedWord.Name)
			}

			if tt.expectHints && len(gs.Hints) == 0 {
				t.Error("Expected hints to be set, but got empty")
			}

			if len(gs.Result) != len(gs.Answer) {
				t.Errorf("Result length = %v, expected %v", len(gs.Result), len(gs.Answer))
			}

			for _, char := range gs.Result {
				if char != '_' {
					t.Errorf("Result contains non-underscore character: %v", char)
				}
			}
		})
	}
}

func TestGameSettings_FilterWords_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		words       *Words
		params      []string
		shouldPanic bool
	}{
		{
			name: "empty words struct",
			words: &Words{
				En: ComplexityWords{},
				Ru: ComplexityWords{},
			},
			params:      []string{"en", "easy", "low"},
			shouldPanic: true,
		},
		{
			name: "empty easy category",
			words: &Words{
				En: ComplexityWords{
					Easy: []Word{},
					Hard: []Word{{Name: "Test"}},
				},
				Ru: ComplexityWords{
					Easy: []Word{},
					Hard: []Word{},
				},
			},
			params:      []string{"en", "easy", "low"},
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GameSettings{Params: tt.params}

			defer func() {
				if r := recover(); r != nil && !tt.shouldPanic {
					t.Errorf("Unexpected panic: %v", r)
				} else if r == nil && tt.shouldPanic {
					t.Error("Expected panic, but none occurred")
				}
			}()

			gs.FilterWords(tt.words)
		})
	}
}

func TestGameSettings_Integration(t *testing.T) {
	mockWords := &Words{
		En: ComplexityWords{
			Easy: []Word{
				{Name: "Test", Hints: []string{"Hint1", "Hint2"}},
			},
		},
		Ru: ComplexityWords{
			Easy: []Word{},
		},
	}

	gs := &GameSettings{
		Params: []string{"en", "easy", "low"},
	}

	selectedWord := gs.FilterWords(mockWords)
	if selectedWord.Name != "Test" {
		t.Errorf("Expected word 'Test', got '%v'", selectedWord.Name)
	}

	if string(gs.Answer) != "test" {
		t.Errorf("Expected answer 'test', got '%v'", string(gs.Answer))
	}

	if len(gs.Result) != 4 {
		t.Errorf("Expected result length 4, got %v", len(gs.Result))
	}

	gs.ProccessLetter('t')
	if string(gs.Result) != "t__t" {
		t.Errorf("Expected result 't__t', got '%v'", string(gs.Result))
	}
	if gs.CountLetters != 2 {
		t.Errorf("Expected count 2, got %v", gs.CountLetters)
	}

	initialState := gs.CurrentState
	gs.ProccessLetter('x')
	if gs.CurrentState != initialState+1 {
		t.Errorf("Expected state %v, got %v", initialState+1, gs.CurrentState)
	}

	hint := gs.GetHint()
	if hint != "Hint1" && hint != "Hint2" {
		t.Errorf("Expected hint to be 'Hint1' or 'Hint2', got '%v'", hint)
	}
}

func TestWord_Structure(t *testing.T) {
	word := Word{
		Name:  "Test",
		Hints: []string{"Hint1", "Hint2"},
	}

	if word.Name != "Test" {
		t.Errorf("Word.Name = %v, expected 'Test'", word.Name)
	}

	if len(word.Hints) != 2 {
		t.Errorf("Word.Hints length = %v, expected 2", len(word.Hints))
	}
}

func TestGameSettings_EdgeCases(t *testing.T) {
	t.Run("get hint modifies original slice", func(t *testing.T) {
		originalHints := []string{"Hint1", "Hint2", "Hint3"}
		gs := &GameSettings{Hints: originalHints}

		hint1 := gs.GetHint()
		hint2 := gs.GetHint()
		hint3 := gs.GetHint()

		if len(gs.Hints) != len(originalHints) {
			t.Error("GetHint should not modify original hints slice")
		}

		hints := []string{hint1, hint2, hint3}
		for _, hint := range hints {
			found := false
			for _, originalHint := range originalHints {
				if hint == originalHint {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Hint '%v' not found in original hints", hint)
			}
		}
	})
}
