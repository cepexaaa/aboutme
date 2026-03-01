package domain

type Labyrinth struct {
	Algorithm string
	Height    uint
	Widht     uint
	Map       [][]Cell
}

type Pair[T any] struct {
	First  T
	Second T
}
