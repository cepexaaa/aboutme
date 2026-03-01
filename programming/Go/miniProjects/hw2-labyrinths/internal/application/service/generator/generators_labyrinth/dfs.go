package generatorslabyrinth

import (
	"math/rand"
	"time"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

type steps domain.Pair[int]

// create maze using dfs
func DFS(w uint, h uint) [][]domain.Cell {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	w += (1 - w%2)
	h += (1 - h%2)

	maze := prepareField(w, h)

	startX := 1 + 2*rnd.Intn(int(w)/2)
	startY := 1 + 2*rnd.Intn(int(h)/2)

	dfs(maze, startX, startY, rnd)

	return maze
}

// use dfs to create maze
func dfs(maze [][]domain.Cell, x, y int, rnd *rand.Rand) {
	maze[y][x] = domain.Empty

	directions := []steps{
		{0, -2}, // up
		{2, 0},  // right
		{0, 2},  // down
		{-2, 0}, // left
	}

	rnd.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	for _, dir := range directions {
		newX := x + dir.First
		newY := y + dir.Second

		if newX > 0 && newX < len(maze[0])-1 && newY > 0 && newY < len(maze)-1 {
			if maze[newY][newX] == domain.Wall {
				maze[y+dir.Second/2][x+dir.First/2] = domain.Empty
				dfs(maze, newX, newY, rnd)
			}
		}
	}
}

// prepare field for new maze
// in the begining it fills Walls
func prepareField(w uint, h uint) [][]domain.Cell {
	maze := make([][]domain.Cell, h)
	for i := range maze {
		maze[i] = make([]domain.Cell, w)
		for j := range maze[i] {
			maze[i][j] = domain.Wall
		}
	}
	return maze
}
