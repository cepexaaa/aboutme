package solverslabyrinths

import (
	"container/heap"
	"fmt"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

type NodeWithCost struct {
	position domain.Pair[uint]
	parent   *NodeWithCost
	cost     int
	coins    int
	index    int
}

type CostPriorityQueue []*NodeWithCost

func (pq CostPriorityQueue) Len() int { return len(pq) }

func (pq CostPriorityQueue) Less(i, j int) bool {
	if pq[i].cost == pq[j].cost {
		return pq[i].coins > pq[j].coins // It will be better than it was if we have more coins
	}
	return pq[i].cost < pq[j].cost
}

func (pq CostPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *CostPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*NodeWithCost)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *CostPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*pq = old[0 : n-1]
	return node
}

// find optimal path based on sand and coins
func Optimal(maze *domain.Labyrinth, start, end domain.Pair[uint]) {
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

	openList := make(CostPriorityQueue, 0)
	visited := make(map[string]bool)

	startNode := &NodeWithCost{
		position: start,
		parent:   nil,
		cost:     0,
		coins:    0,
	}

	heap.Push(&openList, startNode)

	directions := []domain.Pair[int]{
		{First: 0, Second: -1},
		{First: 1, Second: 0},
		{First: 0, Second: 1},
		{First: -1, Second: 0},
	}

	for openList.Len() > 0 {
		current := heap.Pop(&openList).(*NodeWithCost)
		currentX := int(current.position.First)
		currentY := int(current.position.Second)

		visitedKey := fmt.Sprintf("%d,%d", currentX, currentY)

		if currentX == endX && currentY == endY {
			reconstructOptimalPath(result, current)
			result[startY][startX] = domain.Start
			result[endY][endX] = domain.End
			maze.Map = result
			return
		}

		if visited[visitedKey] {
			continue
		}
		visited[visitedKey] = true

		for _, dir := range directions {
			neighborX := currentX + dir.First
			neighborY := currentY + dir.Second

			if !isValid(neighborX, neighborY, result) {
				continue
			}

			cell := result[neighborY][neighborX]
			stepCost := cell.Cost()

			newCost := current.cost + stepCost

			newCoins := current.coins
			if cell == domain.Coin {
				newCoins++
			}

			neighborPos := domain.Pair[uint]{First: uint(neighborX), Second: uint(neighborY)}
			neighborNode := &NodeWithCost{
				position: neighborPos,
				parent:   current,
				cost:     newCost,
				coins:    newCoins,
			}

			found := false
			for i, node := range openList {
				if node.position == neighborPos {
					found = true
					if newCost < node.cost || (newCost == node.cost && newCoins > node.coins) {
						openList[i] = neighborNode
						heap.Init(&openList)
					}
					break
				}
			}

			if !found {
				heap.Push(&openList, neighborNode)
			}
		}
	}
	result[startY][startX] = domain.Start
	result[endY][endX] = domain.End
	maze.Map = result
}

func reconstructOptimalPath(maze [][]domain.Cell, endNode *NodeWithCost) {
	current := endNode

	for current != nil && current.parent != nil {
		x := int(current.position.First)
		y := int(current.position.Second)

		if maze[y][x] != domain.Start && maze[y][x] != domain.End &&
			maze[y][x] != domain.Coin && maze[y][x] != domain.Sand {
			maze[y][x] = domain.Way
		}

		current = current.parent
	}
}
