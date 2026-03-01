package localisation

import (
	_ "embed"
	"encoding/json"
)

type messages struct {
	Russian message `json:"ru"`
	English message `json:"en"`
}

type message struct {
	GameInfo    string `json:"gameInfo"`
	UsingArrows string `json:"usingArrows"`
	Win         string `json:"win"`
	Lose        string `json:"lose"`
	Hint        string `json:"hint"`
}

//go:embed messages.json
var messagesData []byte

func (m message) GetGameInfo() string    { return m.GameInfo }
func (m message) GetUsingArrows() string { return m.UsingArrows }
func (m message) GetWin() string         { return m.Win }
func (m message) GetLose() string        { return m.Lose }
func (m message) GetHint() string        { return m.Hint }

func NewMessageProvider(language string) message {
	var messages messages
	json.Unmarshal(messagesData, &messages)

	var source message
	if language == "ru" {
		source = messages.Russian
	} else {
		source = messages.English
	}

	return source
}
