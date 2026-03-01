package solverslabyrinths

import (
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

type point struct {
	X, Y   int
	Parent *point
}

// use bfs to find path
func BFS(maze *domain.Labyrinth, start, end domain.Pair[uint]) {
	// make a copy
	result := make([][]domain.Cell, len(maze.Map))
	for i := range maze.Map {
		result[i] = make([]domain.Cell, len(maze.Map[i]))
		copy(result[i], maze.Map[i])
	}

	startX, startY := int(start.First), int(start.Second)
	endX, endY := int(end.First), int(end.Second)

	if !isValid(startX, startY, result) || !isValid(endX, endY, result) {
		maze.Map = result
		return
	}

	queue := []*point{}
	visited := make(map[point]bool)

	directions := []point{
		{0, -1, nil},
		{1, 0, nil},
		{0, 1, nil},
		{-1, 0, nil},
	}

	startPoint := &point{startX, startY, nil}
	queue = append(queue, startPoint)
	visited[point{startX, startY, nil}] = true

	var endPoint *point

	// BFS
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.X == endX && current.Y == endY {
			endPoint = current
			break
		}

		for _, dir := range directions {
			neighborX := current.X + dir.X
			neighborY := current.Y + dir.Y

			if !isValid(neighborX, neighborY, result) {
				continue
			}

			neighborPoint := point{neighborX, neighborY, nil}
			if !visited[neighborPoint] {
				visited[neighborPoint] = true
				queue = append(queue, &point{neighborX, neighborY, current})
			}
		}
	}

	if endPoint != nil {
		reconstructPathBFS(result, endPoint)
	}

	result[startY][startX] = domain.Start
	result[endY][endX] = domain.End

	maze.Map = result
}

func reconstructPathBFS(maze [][]domain.Cell, endPoint *point) {
	current := endPoint.Parent

	for current != nil && current.Parent != nil {
		x, y := current.X, current.Y
		if maze[y][x] != domain.Start && maze[y][x] != domain.End {
			maze[y][x] = domain.Way
		}
		current = current.Parent
	}
}
