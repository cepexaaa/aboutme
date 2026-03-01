package solverslabyrinths

import (
	"container/heap"
	"math"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

type dNode struct {
	position domain.Pair[uint]
	parent   *dNode
	distance int // distance from start
	index    int // index in heap
}

type priorityQueueDijkstra []*dNode

func (pq priorityQueueDijkstra) Len() int { return len(pq) }

func (pq priorityQueueDijkstra) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}

func (pq priorityQueueDijkstra) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueueDijkstra) Push(x interface{}) {
	n := len(*pq)
	node := x.(*dNode)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *priorityQueueDijkstra) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*pq = old[0 : n-1]
	return node
}

// use Dijkstra for find path
// See more: https://habr.com/ru/companies/otus/articles/748470/
func Dijkstra(maze *domain.Labyrinth, start, end domain.Pair[uint]) {
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

	distances := make(map[domain.Pair[uint]]int)
	visited := make(map[domain.Pair[uint]]bool)
	parent := make(map[domain.Pair[uint]]domain.Pair[uint])

	for y := range result {
		for x := range result[y] {
			if isValid(x, y, result) {
				pos := domain.Pair[uint]{First: uint(x), Second: uint(y)}
				distances[pos] = math.MaxInt32
			}
		}
	}

	startPos := start
	distances[startPos] = 0
	parent[startPos] = startPos

	pq := make(priorityQueueDijkstra, 0)
	heap.Push(&pq, &dNode{
		position: startPos,
		distance: 0,
	})

	directions := []domain.Pair[int]{
		{First: 0, Second: -1},
		{First: 1, Second: 0},
		{First: 0, Second: 1},
		{First: -1, Second: 0},
	}

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*dNode)
		currentX := int(current.position.First)
		currentY := int(current.position.Second)

		if currentX == endX && currentY == endY {
			reconstructPathDijkstra(result, parent, current.position)
			result[startY][startX] = domain.Start
			result[endY][endX] = domain.End
			maze.Map = result
			return
		}

		if visited[current.position] {
			continue
		}
		visited[current.position] = true

		for _, dir := range directions {
			neighborX := currentX + dir.First
			neighborY := currentY + dir.Second
			neighborPos := domain.Pair[uint]{First: uint(neighborX), Second: uint(neighborY)}

			if !isValid(neighborX, neighborY, result) {
				continue
			}

			if visited[neighborPos] {
				continue
			}

			newDistance := distances[current.position] + 1

			if newDistance < distances[neighborPos] {
				distances[neighborPos] = newDistance
				parent[neighborPos] = current.position

				heap.Push(&pq, &dNode{
					position: neighborPos,
					distance: newDistance,
				})
			}
		}
	}
}

// leaving the final route
func reconstructPathDijkstra(maze [][]domain.Cell, parent map[domain.Pair[uint]]domain.Pair[uint], end domain.Pair[uint]) {
	current := end

	for {
		prev, exists := parent[current]
		if !exists || current == prev {
			break
		}

		x := int(current.First)
		y := int(current.Second)
		if maze[y][x] != domain.Start && maze[y][x] != domain.End {
			maze[y][x] = domain.Way
		}
		current = prev
	}
}
