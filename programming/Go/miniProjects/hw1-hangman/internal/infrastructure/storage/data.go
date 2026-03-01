package storage

import (
	_ "embed"
	"encoding/json"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw1-hangman/internal/domain"
)

//go:embed data.json
var embeddedData []byte

func GetData() (*domain.Words, error) {
	var words domain.Words
	err := json.Unmarshal(embeddedData, &words)
	return &words, err
}
