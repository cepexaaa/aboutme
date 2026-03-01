package main

import (
	"fmt"
	"os"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw1-hangman/internal/application"
)

func main() {
	switch len(os.Args) {
	case 3:
		application.NoInteractiveGame(os.Args)
	case 1:
		application.InteractiveGame()
	default:
		fmt.Println("You can start game with 0 or 2 words")
	}
}
