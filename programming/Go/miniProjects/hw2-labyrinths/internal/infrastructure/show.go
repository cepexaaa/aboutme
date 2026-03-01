package infrastructure

import (
	"io"
	"strings"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

var (
	// different unicode walls
	walls = [16]rune{' ', '╷', '╶', '┌', '╴', '┐', '─', '┬', '╵', '│', '└', '├', '┘', '┤', '┴', '┼'}
	// start, finish and path in maze
	solveWay = []rune{'▲', '★', '◦', '¢', '%'}
)

type getRune func(maze [][]domain.Cell, x, y int) rune

// show result maze in {writer}
// can paint using unicode symbols
func ShowResult(maze *domain.Labyrinth, writer io.Writer, isUnicode bool, isSolve bool) {
	var result strings.Builder
	step := 1
	var gr getRune
	if isUnicode {
		gr = unicodeWall
		if !isSolve {
			step = 2
		}
	} else {
		gr = cellToRune
	}
	for y := 0; y < len(maze.Map); y += step {
		for x := 0; x < len(maze.Map[0]); x += step {
			result.WriteRune(gr(maze.Map, x, y))
		}
		result.WriteRune('\n')
	}
	writer.Write([]byte(result.String()))
}

// convert abstract value of cell to certain ASCII symbol
func cellToRune(maze [][]domain.Cell, x, y int) rune {
	c := maze[y][x]
	switch c {
	case domain.Wall:
		return '#'
	case domain.Start:
		return 'O'
	case domain.End:
		return 'X'
	case domain.Empty:
		return ' '
	case domain.Way:
		return '.'
	case domain.Sand:
		return '%'
	case domain.Coin:
		return '$'
	default:
		return '?'
	}
}

// convert abstract value of cell to certain Unicode symbol
func unicodeWall(maze [][]domain.Cell, x, y int) rune {
	switch maze[y][x] {
	case domain.Start:
		return solveWay[0]
	case domain.End:
		return solveWay[1]
	case domain.Way:
		return solveWay[2]
	case domain.Coin:
		return solveWay[3]
	case domain.Sand:
		return solveWay[4]
	case domain.Empty:
		return walls[0]
	default:
		indOfWall := 0
		if y > 0 && maze[y-1][x] == domain.Wall {
			indOfWall += 8
		}
		if y < len(maze)-1 && maze[y+1][x] == domain.Wall {
			indOfWall += 1
		}
		if x > 0 && maze[y][x-1] == domain.Wall {
			indOfWall += 4
		}
		if x < len(maze[0])-1 && maze[y][x+1] == domain.Wall {
			indOfWall += 2
		}
		return walls[indOfWall]
	}
}
