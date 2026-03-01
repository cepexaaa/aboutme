package infrastructure

import (
	"bufio"
	"fmt"
	"os"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

// read maze from file
func ReadMaze(fileName string) [][]domain.Cell {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error durind opening the file:", err)
		return nil
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	firstRune, _, err := reader.ReadRune()
	if err != nil {
		return nil
	}

	isUnicode := firstRune == 65533 || firstRune == 9484

	err = reader.UnreadRune()
	if err != nil {
		fmt.Println("Error unreading rune:", err)
		return nil
	}

	scanner := bufio.NewScanner(reader)
	var result [][]domain.Cell
	if isUnicode {
		// it is compressed field. That why we need to expand it
		position := 0
		for scanner.Scan() {
			line := []rune(scanner.Text())
			position = 0
			row := make([]domain.Cell, len([]rune(line))*2-1)
			row2 := make([]domain.Cell, len([]rune(line))*2-1)
			for _, r := range line {
				// each cell control itself + right/left and up/down from it
				ind := getIndex(r)
				row[position] = domain.Wall
				if ind%2 == 1 {
					row2[position] = domain.Wall
				}
				if ind/2%2 == 1 && len(row) > position+1 {
					row[position+1] = domain.Wall
				}
				position += 2
			}
			result = append(result, row)
			result = append(result, row2)
		}
	} else {
		for scanner.Scan() {
			line := scanner.Text()
			row := make([]domain.Cell, 0, len([]rune(line)))
			for _, r := range line {
				row = append(row, runeToCell(r))
			}
			result = append(result, row)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error during reading the file:", err)
	}
	for _, line := range result {
		if len(line) != len(result[0]) {
			fmt.Println("uncorrent maze format")
			return nil
		}
	}
	return result
}

// convert symbols from ASCII
func runeToCell(r rune) domain.Cell {
	switch r {
	case '#':
		return domain.Wall
	case 'O':
		return domain.Start
	case 'X':
		return domain.End
	case ' ':
		return domain.Empty
	case '.':
		return domain.Way
	case '%':
		return domain.Sand
	case '$':
		return domain.Coin
	default:
		return domain.Wall
	}
}

func getIndex(c rune) int {
	for i := range walls {
		if walls[i] == c {
			return i
		}
	}
	return 0
}
