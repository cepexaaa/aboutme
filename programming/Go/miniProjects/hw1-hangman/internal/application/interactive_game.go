package application

import (
	"fmt"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw1-hangman/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw1-hangman/internal/infrastructure/images"
	rawmode "gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw1-hangman/internal/infrastructure/raw_mode"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw1-hangman/internal/infrastructure/storage"
)

func InteractiveGame() {
	var rawMode rawmode.RawMode
	rawMode.RunRawMode()
	defer rawMode.StopRawMode()

	setUpGame()
}

func setUpGame() {
	game := domain.GameSettings{Quit: false}
	images.Prepearing(&game)
	if game.Quit {
		return
	}
	words, err := storage.GetData()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	game.FilterWords(words)

	images.Proccess(&game)
}
