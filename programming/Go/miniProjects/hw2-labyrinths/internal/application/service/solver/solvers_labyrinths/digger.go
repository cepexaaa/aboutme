package solverslabyrinths

import (
	"fmt"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

// go towards final point
func Digger(maze *domain.Labyrinth, start, end domain.Pair[uint]) {
	fmt.Println("Сила есть - ума не надо")

	myCoords := domain.Pair[uint]{First: start.First, Second: start.Second}

	for myCoords != end {
		goTowards(maze.Map, myCoords, end.First, &(myCoords.First))
		goTowards(maze.Map, myCoords, end.Second, &(myCoords.Second))

	}
	maze.Map[start.Second][start.First] = domain.Start
	maze.Map[end.Second][end.First] = domain.End
}

// find certain stap in correct destination
func goTowards(field [][]domain.Cell, myCoords domain.Pair[uint], aim uint, me *uint) {
	l := int(aim) - int(*me)
	*me += uint(l / normalized(l))
	field[myCoords.Second][myCoords.First] = domain.Way

}

func normalized(x int) int {
	switch {
	case x < 0:
		return -x
	case x == 0:
		return 1
	default:
		return x
	}
}
