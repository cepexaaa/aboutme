package solverslabyrinths

import (
	"container/heap"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

type node struct {
	position domain.Pair[uint]
	parent   *node
	gCost    int // cost from start
	hCost    int // heuristic cost to aim
	fCost    int // both cost (gCost + hCost)
}

type priorityQueue []*node

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].fCost < pq[j].fCost
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*node))
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// use Astar for find path
// See more: https://habr.com/ru/companies/otus/articles/748470/
func Astar(maze *domain.Labyrinth, start, end domain.Pair[uint]) {
	// make a copy
	result := make([][]domain.Cell, len(maze.Map))
	for i := range maze.Map {
		result[i] = make([]domain.Cell, len(maze.Map[i]))
		copy(result[i], maze.Map[i])
	}

	startX, startY := int(start.First), int(start.Second)
	endX, endY := int(end.First), int(end.Second)

	if !isValid(startX, startY, result) || !isValid(endX, endY, result) {
		return
	}

	openList := make(priorityQueue, 0)
	closedList := make(map[domain.Pair[uint]]bool)

	startNode := &node{
		position: start,
		parent:   nil,
		gCost:    0,
		hCost:    heuristic(startX, startY, endX, endY),
		fCost:    0 + heuristic(startX, startY, endX, endY),
	}

	heap.Push(&openList, startNode)

	directions := []domain.Pair[int]{
		{First: 0, Second: -1},
		{First: 1, Second: 0},
		{First: 0, Second: 1},
		{First: -1, Second: 0},
	}

	for openList.Len() > 0 {
		current := heap.Pop(&openList).(*node)
		currentX := int(current.position.First)
		currentY := int(current.position.Second)

		// We have been gone to the finish
		if currentX == endX && currentY == endY {
			reconstructPath(result, current)
			result[startY][startX] = domain.Start
			result[endY][endX] = domain.End
			maze.Map = result
			return
		}

		closedList[current.position] = true

		for _, dir := range directions {
			neighborX := currentX + dir.First
			neighborY := currentY + dir.Second

			if !isValid(neighborX, neighborY, result) {
				continue
			}

			neighborPos := domain.Pair[uint]{First: uint(neighborX), Second: uint(neighborY)}

			if closedList[neighborPos] {
				continue
			}

			newGCost := current.gCost + 1

			var neighborNode *node
			for _, node := range openList {
				if node.position.First == uint(neighborX) && node.position.Second == uint(neighborY) {
					neighborNode = node
					break
				}
			}

			if neighborNode == nil {
				neighborNode = &node{
					position: neighborPos,
					parent:   current,
					gCost:    newGCost,
					hCost:    heuristic(neighborX, neighborY, endX, endY),
					fCost:    newGCost + heuristic(neighborX, neighborY, endX, endY),
				}
				heap.Push(&openList, neighborNode)
			} else if newGCost < neighborNode.gCost {
				neighborNode.parent = current
				neighborNode.gCost = newGCost
				neighborNode.fCost = newGCost + neighborNode.hCost
				heap.Init(&openList)
			}
		}
	}
}

// leaving the final route
func reconstructPath(maze [][]domain.Cell, endNode *node) {
	current := endNode
	for current != nil {
		x := int(current.position.First)
		y := int(current.position.Second)

		if maze[y][x] != domain.Start && maze[y][x] != domain.End {
			maze[y][x] = domain.Way
		}
		current = current.parent
	}
}

func heuristic(x1, y1, x2, y2 int) int {
	return abs(x1-x2) + abs(y1-y2)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// check that is valid cell to be in
func isValid(x, y int, maze [][]domain.Cell) bool {
	if y < 0 || y >= len(maze) || x < 0 || x >= len(maze[0]) {
		return false
	}
	return maze[y][x] == domain.Empty || maze[y][x] == domain.Start || maze[y][x] == domain.End || maze[y][x] == domain.Sand || maze[y][x] == domain.Coin
}
