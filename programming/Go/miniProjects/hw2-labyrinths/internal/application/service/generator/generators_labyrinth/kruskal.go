package generatorslabyrinth

import (
	"math/rand"
	"time"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

type Edge struct {
	X1, Y1, X2, Y2 int
}

type DisjointSet struct {
	parent map[int]int
	rank   map[int]int
}

func newDisjointSet() *DisjointSet {
	return &DisjointSet{
		parent: make(map[int]int),
		rank:   make(map[int]int),
	}
}

func (ds *DisjointSet) Find(x int) int {
	if ds.parent[x] != x {
		ds.parent[x] = ds.Find(ds.parent[x])
	}
	return ds.parent[x]
}

func (ds *DisjointSet) Union(x, y int) {
	rootX := ds.Find(x)
	rootY := ds.Find(y)

	if rootX != rootY {
		if ds.rank[rootX] < ds.rank[rootY] {
			ds.parent[rootX] = rootY
		} else if ds.rank[rootX] > ds.rank[rootY] {
			ds.parent[rootY] = rootX
		} else {
			ds.parent[rootY] = rootX
			ds.rank[rootX]++
		}
	}
}

func Kruskal(w uint, h uint) [][]domain.Cell {
	w += (1 - w%2)
	h += (1 - h%2)

	width := int(w)
	height := int(h)

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	maze := prepareField(w, h)

	ds := newDisjointSet()

	edges := make([]Edge, 0)

	for y := 1; y < height-1; y += 2 {
		for x := 1; x < width-1; x += 2 {
			cellID := y*width + x
			ds.parent[cellID] = cellID
			ds.rank[cellID] = 0
			maze[y][x] = domain.Empty

			if x+2 < width-1 {
				edges = append(edges, Edge{x, y, x + 2, y})
			}
			if y+2 < height-1 {
				edges = append(edges, Edge{x, y, x, y + 2})
			}
		}
	}

	rnd.Shuffle(len(edges), func(i, j int) {
		edges[i], edges[j] = edges[j], edges[i]
	})

	for _, edge := range edges {
		cell1ID := edge.Y1*width + edge.X1
		cell2ID := edge.Y2*width + edge.X2

		if ds.Find(cell1ID) != ds.Find(cell2ID) {
			ds.Union(cell1ID, cell2ID)

			wallX := (edge.X1 + edge.X2) / 2
			wallY := (edge.Y1 + edge.Y2) / 2
			maze[wallY][wallX] = domain.Empty
		}
	}

	return maze
}
