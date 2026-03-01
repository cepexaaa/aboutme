package domain

type Cell uint8

const (
	Empty Cell = iota
	Wall
	Way
	Start
	End
	Sand
	Coin
)

func (c Cell) Cost() int {
	switch c {
	case Wall:
		return 100
	case Empty, Start, End, Way:
		return 1
	case Sand:
		return 3
	case Coin:
		return -1
	default:
		return 1
	}
}
