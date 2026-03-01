package domain

type Point struct {
	X float64
	Y float64
}

type Color struct {
	R float64
	G float64
	B float64
	A float64
}

type AffineParams struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
	C float64 `json:"c"`
	D float64 `json:"d"`
	E float64 `json:"e"`
	F float64 `json:"f"`
}

type FunctionConfig struct {
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
}

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}
