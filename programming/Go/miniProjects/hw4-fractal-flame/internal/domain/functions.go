package domain

import (
	"errors"
	"math"
)

type TransformationFunction interface {
	Apply(p Point) Point
}

type Swirl struct{}

func (s Swirl) Apply(p Point) Point {
	r2 := p.X*p.X + p.Y*p.Y
	sinr2 := math.Sin(r2)
	cosr2 := math.Cos(r2)

	return Point{
		X: p.X*sinr2 - p.Y*cosr2,
		Y: p.X*cosr2 + p.Y*sinr2,
	}
}

type Horseshoe struct{}

func (h Horseshoe) Apply(p Point) Point {
	r := math.Sqrt(p.X*p.X + p.Y*p.Y)
	if r == 0 {
		return Point{X: 0, Y: 0}
	}

	return Point{
		X: (1.0 / r) * (p.X - p.Y) * (p.X + p.Y),
		Y: (2.0 / r) * p.X * p.Y,
	}
}

type Disk struct{}

func (d Disk) Apply(p Point) Point {
	r := math.Sqrt(p.X*p.X + p.Y*p.Y)
	if r == 0 {
		return Point{X: 0, Y: 0}
	}

	theta := math.Atan2(p.Y, p.X) / math.Pi
	return Point{
		X: theta * math.Sin(math.Pi*r),
		Y: theta * math.Cos(math.Pi*r),
	}
}

type Heart struct{}

func (h Heart) Apply(p Point) Point {
	r := math.Sqrt(p.X*p.X + p.Y*p.Y)
	theta := math.Atan2(p.Y, p.X)

	return Point{
		X: r * math.Sin(theta*r),
		Y: -r * math.Cos(theta*r),
	}
}

type Linear struct{}

func (h Linear) Apply(p Point) Point {
	return p
}

func GetFunctionByName(name string) (TransformationFunction, error) {
	switch name {
	case "swirl":
		return Swirl{}, nil
	case "horseshoe":
		return Horseshoe{}, nil
	case "disk":
		return Disk{}, nil
	case "heart":
		return Heart{}, nil
	case "linear":
		return Linear{}, nil
	default:
		return nil, errors.New("unknown function: " + name)
	}
}
