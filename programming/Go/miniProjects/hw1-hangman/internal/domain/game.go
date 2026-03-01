package domain

import (
	"math/rand"
	"strings"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw1-hangman/internal/infrastructure/localisation"
)

type GameSettings struct {
	Params        []string // lang, complexity
	Quit          bool
	CurrentChoise int
	CurrentState  int
	Result        []rune
	Answer        []rune
	Hints         []string
	CountLetters  int
	Info          gameMessages
}

type gameMessages interface {
	GetGameInfo() string
	GetUsingArrows() string
	GetWin() string
	GetLose() string
	GetHint() string
}

func (g *GameSettings) GetHint() string {
	if len(g.Hints) == 0 {
		return "No any hints"
	}
	r := rand.Intn(len(g.Hints))
	return g.Hints[r]
}

func (g *GameSettings) ProccessLetter(r rune) {
	guess := false
	for i, c := range g.Answer {
		if r == c {
			guess = true
			if g.Result[i] == '_' {
				g.CountLetters++
				g.Result[i] = c
			}
		}
	}
	if !guess {
		g.CurrentState++
	}
}

func (g *GameSettings) FilterWords(words *Words) Word {
	var filter1 ComplexityWords
	r := rand.Intn(2)
	if len(g.Params) < 2 {
		g.Params = []string{"random", "random"}
	}
	if g.Params[1] == "en" || r == 0 && g.Params[1] == "random" {
		g.Params[1] = "en"
		filter1 = words.En
	} else {
		g.Params[1] = "ru"
		filter1 = words.Ru
	}
	g.Info = localisation.NewMessageProvider(g.Params[1])
	var filter2 []Word
	r = rand.Intn(2)
	if g.Params[0] == "easy" || r == 0 && g.Params[0] == "random" {
		g.Params[0] = "easy"
		filter2 = filter1.Easy
	} else {
		g.Params[0] = "hard"
		filter2 = filter1.Hard
	}
	r = rand.Intn(len(filter2))
	g.Answer = []rune(strings.ToLower(filter2[r].Name))
	g.Hints = filter2[r].Hints
	g.setResult()
	return filter2[r]
}

func (g *GameSettings) setResult() {
	for i := 0; i < len(g.Answer); i++ {
		g.Result = append(g.Result, '_')
	}
}
