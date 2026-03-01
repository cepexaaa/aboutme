package generatorslabyrinth

import (
	"math/rand"
	"time"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

type direction int

type point domain.Pair[uint]

const (
	North direction = iota
	South
	East
	West
)

// use algorithm of Prim to create new maze
func Prim(w uint, h uint) [][]domain.Cell {
	w += (1 - w%2)
	h += (1 - h%2)

	maze := prepareField(w, h)

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	startX := uint(1 + 2*rnd.Intn(int(w/2)))
	startY := uint(1 + 2*rnd.Intn(int(h/2)))

	maze[startY][startX] = domain.Empty

	toCheck := []point{}
	if startY >= 2 {
		toCheck = append(toCheck, point{startX, startY - 2})
	}
	if startY+2 < h {
		toCheck = append(toCheck, point{startX, startY + 2})
	}
	if startX >= 2 {
		toCheck = append(toCheck, point{startX - 2, startY})
	}
	if startX+2 < w {
		toCheck = append(toCheck, point{startX + 2, startY})
	}

	for len(toCheck) > 0 {
		index := rnd.Intn(len(toCheck))
		cell := toCheck[index]
		x := cell.First
		y := cell.Second

		toCheck = append(toCheck[:index], toCheck[index+1:]...)

		if maze[y][x] == domain.Empty {
			continue
		}

		maze[y][x] = domain.Empty

		directions := []direction{North, South, East, West}
		rnd.Shuffle(len(directions), func(i, j int) {
			directions[i], directions[j] = directions[j], directions[i]
		})

		found := false
		for _, dir := range directions {
			switch dir {
			case North:
				if y >= 2 && maze[y-2][x] == domain.Empty {
					maze[y-1][x] = domain.Empty
					found = true
				}
			case South:
				if y+2 < h && maze[y+2][x] == domain.Empty {
					maze[y+1][x] = domain.Empty
					found = true
				}
			case East:
				if x >= 2 && maze[y][x-2] == domain.Empty {
					maze[y][x-1] = domain.Empty
					found = true
				}
			case West:
				if x+2 < w && maze[y][x+2] == domain.Empty {
					maze[y][x+1] = domain.Empty
					found = true
				}
			}
			if found {
				break
			}
		}

		if y >= 2 && maze[y-2][x] == domain.Wall {
			toCheck = append(toCheck, point{x, y - 2})
		}
		if y+2 < h && maze[y+2][x] == domain.Wall {
			toCheck = append(toCheck, point{x, y + 2})
		}
		if x >= 2 && maze[y][x-2] == domain.Wall {
			toCheck = append(toCheck, point{x - 2, y})
		}
		if x+2 < w && maze[y][x+2] == domain.Wall {
			toCheck = append(toCheck, point{x + 2, y})
		}
	}

	return maze
}
