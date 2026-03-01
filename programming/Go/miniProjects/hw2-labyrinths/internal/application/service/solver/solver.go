package solver

import (
	"errors"
	"os"

	solverslabyrinths "gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/application/service/solver/solvers_labyrinths"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

type Solver struct {
	Maze  *domain.Labyrinth
	Start domain.Pair[uint]
	End   domain.Pair[uint]
	File  os.File
}

func NewSolver(start, end domain.Pair[uint]) *Solver {
	return &Solver{Start: start, End: end}
}

func (s *Solver) NewMaze(Algorithm string, Height uint, Width uint, field [][]domain.Cell) {
	s.Maze = &domain.Labyrinth{Algorithm: Algorithm, Height: Height, Widht: Width, Map: field}
}

// is valid parameters to solve labirint
func (s *Solver) IsLogicValid() error {
	if int(s.Start.Second) >= len(s.Maze.Map) || int(s.Start.First) >= len(s.Maze.Map[0]) {
		return errors.New("start point out of labyrinth. Change them")
	}
	if int(s.End.Second) >= len(s.Maze.Map) || int(s.End.First) >= len(s.Maze.Map[0]) {
		return errors.New("end point out of labyrinth. Change them")
	}
	if s.Maze.Map[s.Start.Second][s.Start.First] == domain.Wall {
		s.Maze.Map[s.Start.Second][s.Start.First] = domain.Start
		return errors.New("start on wall. Change start point")
	}
	if s.Maze.Map[s.End.Second][s.End.First] == domain.Wall {
		s.Maze.Map[s.End.Second][s.End.First] = domain.End
		return errors.New("end on wall. Change end point")
	}
	return nil
}

// use algorithm for solving
func (s *Solver) Solve() {
	switch s.Maze.Algorithm {
	case "dijkstra":
		cleanMap(s.Maze.Map)
		solverslabyrinths.Dijkstra(s.Maze, s.Start, s.End)
	case "astar":
		cleanMap(s.Maze.Map)
		solverslabyrinths.Astar(s.Maze, s.Start, s.End)
	case "bfs":
		cleanMap(s.Maze.Map)
		solverslabyrinths.BFS(s.Maze, s.Start, s.End)
	case "digger":
		cleanMap(s.Maze.Map)
		solverslabyrinths.Digger(s.Maze, s.Start, s.End)
	default:
		solverslabyrinths.Optimal(s.Maze, s.Start, s.End)
	}
}

// remove sand and coins from maze
func cleanMap(field [][]domain.Cell) {
	for y := range field {
		for x := range field[y] {
			if field[y][x] == domain.Coin || field[y][x] == domain.Sand {
				field[y][x] = domain.Empty
			}
		}
	}
}
