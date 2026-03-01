package images

import (
	"bufio"
	"embed"
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw1-hangman/internal/domain"
	rawmode "gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw1-hangman/internal/infrastructure/raw_mode"
)

//go:embed frames/title.txt
var introFile string

//go:embed victory/dance.txt
var victoryDanceFile string

//go:embed settings/*.txt
var settingsFiles embed.FS

//go:embed states/*
var statesFiles embed.FS

func Prepearing(game *domain.GameSettings) {
	rawmode.ClearScreen()
	intro()
	rawmode.ClearScreen()

	settings(game)
	if game.Quit {
		fmt.Println("You are exiting")
		time.Sleep(time.Second)
		return
	}
}

func PrintImage(image []string) {
	rawmode.ClearScreen()
	for i, line := range image {
		rawmode.MoveCursor(5, i+2)
		fmt.Print(line)
	}
}

func Proccess(game *domain.GameSettings) {
	game.CurrentState = 1
	lang := game.Params[0]
	printGameInfo(game, lang, "")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		game.Quit = true
	}()

	for !game.Quit {
		buf := make([]byte, 1)
		n, _ := os.Stdin.Read(buf)

		if n > 0 {
			char := buf[0]

			switch char {
			case 27: // ESC
				if !checkForMoreData() {
					game.Quit = true
				}
			case 43: // +
				printGameInfo(game, lang, game.GetHint())
			default:
				letter, err := readUTF8Char(char)
				if err != nil {
					fmt.Println("something was wrong with error in reading")
					fmt.Println(err.Error())
				}
				game.ProccessLetter(letter)
				printGameInfo(game, lang, "")
			}
		}
	}
}

func getImage(state int) []string {
	if !(0 < state && state <= 10) {
		state = 10
	}
	bytes, _ := statesFiles.ReadFile("states/" + strconv.Itoa(state))
	return strings.Split(string(bytes), "\n")
}

func intro() {
	lines := strings.Split(introFile, "\n")
	rawmode.ClearScreen()
	for i, line := range lines {
		rawmode.MoveCursor(5, i+2)
		for _, c := range line {
			fmt.Print(string(c))
			time.Sleep(time.Millisecond * 10)
		}
	}
	time.Sleep(time.Millisecond * 500)
}

func settings(game *domain.GameSettings) {
	files, err := settingsFiles.ReadDir("settings")
	if err != nil {
		fmt.Println(err.Error())
		time.Sleep(time.Second * 3)
		return
	}
	for _, file := range files {
		if !game.Quit {
			bytes, _ := settingsFiles.ReadFile("settings/" + file.Name())
			renderPage(game, strings.Split(string(bytes), "\n"))
		} else {
			rawmode.ClearScreen()
			break
		}
	}
	rawmode.ClearScreen()
}

func renderPage(game *domain.GameSettings, s []string) {
	if game.CurrentChoise+2 >= len(s) {
		fmt.Println("index out of bonds")
		return
	}
	s[game.CurrentChoise+2] = "->" + s[game.CurrentChoise+2]
	PrintImage(s)
	buf := make([]byte, 1)
	for !game.Quit {
		n, _ := os.Stdin.Read(buf)
		if n > 0 {
			char := buf[0]
			if char == 27 {
				seq := make([]byte, 2)
				n2, err2 := os.Stdin.Read(seq)
				if err2 != nil || n2 < 2 {
					game.Quit = true
					break
				}
				if seq[0] == 91 {
					switch seq[1] {
					case 65: //up
						if game.CurrentChoise > 0 {
							moveArrow(game, &s, -1)
						}
					case 66: //down
						if game.CurrentChoise < (len(s) - 6) {
							moveArrow(game, &s, 1)
						}
					}
				} else {
					game.Quit = true
				}

			} else if char == 13 || char == 10 { // Enter
				game.Params = append(game.Params, strings.TrimSpace(s[game.CurrentChoise+2][2:]))
				game.CurrentChoise = 0
				break
			} else {
				rawmode.MoveCursor(5, len(s)+2)
				fmt.Print("Use down arrow and up arrow to make a choice")
			}
		}
	}
}

func moveArrow(game *domain.GameSettings, s *[]string, move int) {
	(*s)[game.CurrentChoise+2] = (*s)[game.CurrentChoise+2][2:] + "  "
	game.CurrentChoise += move
	(*s)[game.CurrentChoise+2] = "->" + (*s)[game.CurrentChoise+2]
	PrintImage((*s))
}

func readUTF8Char(firstByte byte) (rune, error) {
	reader := bufio.NewReader(os.Stdin)
	switch {
	case firstByte&0x80 == 0x00:
		return rune(firstByte), nil
	case firstByte&0xF0 == 0xE0:
		reader.ReadByte()
		reader.ReadByte()
		return 0, nil
	case firstByte&0xF8 == 0xF0:
		reader.ReadByte()
		reader.ReadByte()
		reader.ReadByte()
		return 0, nil
	default:
		b, err := reader.ReadByte()
		if err != nil {
			return ' ', err
		}
		bb := []rune(strings.ToLower(string([]byte{firstByte, b})))[0]
		return bb, nil
	}
}

func printGameInfo(game *domain.GameSettings, lang string, hint string) {
	if game.CurrentState == 10 {
		printGameLose(game)
		return
	} else if game.CountLetters == len(game.Answer) {
		printVictory(game)
		return
	}

	image := getImage(game.CurrentState)
	PrintImage(image)

	rawmode.MoveCursor(5, len(image)+2)
	fmt.Println(string(game.Result))

	rawmode.MoveCursor(5, len(image)+4)
	fmt.Printf(game.Info.GetGameInfo(), 10-game.CurrentState)
	if hint != "" {
		fmt.Println(game.Info.GetHint() + hint)
	}
}

func printGameLose(game *domain.GameSettings) {
	rawmode.ClearScreen()
	fmt.Println(game.Info.GetLose() + string(game.Answer))
	game.Quit = true
	time.Sleep(time.Second * 2)
}

func printVictory(game *domain.GameSettings) {
	rawmode.ClearScreen()
	fmt.Println(game.Info.GetWin())
	time.Sleep(time.Second)
	victoryDance(game)
	game.Quit = true
}

func victoryDance(game *domain.GameSettings) {
	s := strings.Split(victoryDanceFile, "\n")
	frame := make([]string, 5)
	frame[4] = "  (Ctrl+C)"
	frame[0] = game.Info.GetWin()
	for i := 0; !game.Quit; i += 4 {
		copy(frame[1:4], s[i+1:i+4])
		PrintImage(frame)
		time.Sleep(time.Millisecond * 200)
		i %= 44
	}
}

func checkForMoreData() bool {
	ch := make(chan bool, 1)
	reader := bufio.NewReader(os.Stdin)

	go func() {
		reader.Peek(1)
		ch <- true
	}()

	select {
	case <-ch:
		buf := make([]byte, 1)
		reader.Read(buf)
		reader.Read(buf)
		return true
	case <-time.After(50 * time.Millisecond):
		return false
	}
}
