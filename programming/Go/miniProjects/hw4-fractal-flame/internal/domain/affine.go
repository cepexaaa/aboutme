package domain

import (
	"math"
	"math/rand"
)

type AffineTransform struct {
	Params AffineParams
	Color  Color
}

func (at *AffineTransform) Apply(p Point) Point {
	return Point{
		X: at.Params.A*p.X + at.Params.B*p.Y + at.Params.C,
		Y: at.Params.D*p.X + at.Params.E*p.Y + at.Params.F,
	}
}

func NewRandomAffine(seed float64, random *rand.Rand) AffineTransform {
	sin := math.Sin(seed)
	cos := math.Cos(seed)

	return AffineTransform{
		Params: AffineParams{
			A: 0.7*cos + (random.Float64()-0.5)*0.1,
			B: -0.7*sin + (random.Float64()-0.5)*0.1,
			C: 0.3*sin + (random.Float64()-0.5)*0.1,
			D: 0.7*sin + (random.Float64()-0.5)*0.1,
			E: 0.7*cos + (random.Float64()-0.5)*0.1,
			F: 0.3*cos + (random.Float64()-0.5)*0.1,
		},
		Color: Color{
			R: 0.5 + 0.5*sin,
			G: 0.5 + 0.5*cos,
			B: 0.5 + 0.5*math.Sin(seed+1),
			A: 1.0,
		},
	}
}
