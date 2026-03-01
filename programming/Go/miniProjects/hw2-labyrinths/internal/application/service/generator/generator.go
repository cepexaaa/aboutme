package generator

import (
	"math/rand"
	"time"

	generatorslabyrinth "gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/application/service/generator/generators_labyrinth"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

type Generator struct {
	Maze *domain.Labyrinth
}

func NewGenerator(algo string, w uint, h uint) *Generator {
	return &Generator{Maze: &domain.Labyrinth{Algorithm: algo, Widht: w, Height: h}}
}

func (g *Generator) Generate() {
	switch g.Maze.Algorithm {
	case "dfs":
		g.Maze.Map = generatorslabyrinth.DFS(g.Maze.Widht, g.Maze.Height)
	case "prim":
		g.Maze.Map = generatorslabyrinth.Prim(g.Maze.Widht, g.Maze.Height)
	default:
		g.Maze.Map = generatorslabyrinth.Kruskal(g.Maze.Widht, g.Maze.Height)
	}
	addSurfaces(g.Maze.Map)
}

// Adding surfaces (sand or coins) in maze
func addSurfaces(maze [][]domain.Cell) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for y := range maze {
		for x := range maze[y] {
			if maze[y][x] == domain.Empty {
				switch rnd.Intn(20) {
				case 0, 1: // 10% - send
					maze[y][x] = domain.Sand
				case 2: // 5% - coin
					maze[y][x] = domain.Coin
					// 85% Empty
				}
			}
		}
	}
}
