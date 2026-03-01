package domain

import (
	"math"
)

type SymmetryTransformer struct {
	level int
	angle float64
}

func NewSymmetryTransformer(level int) *SymmetryTransformer {
	if level < 1 {
		level = 1
	}

	return &SymmetryTransformer{
		level: level,
		angle: (2 * math.Pi) / float64(level),
	}
}

func (st *SymmetryTransformer) Apply(point Point) []Point {
	if st.level == 1 {
		return []Point{point}
	}

	points := make([]Point, st.level)
	points[0] = point

	for i := 1; i < st.level; i++ {
		angle := st.angle * float64(i)
		points[i] = st.rotatePoint(point, angle)
	}

	return points
}

func (st *SymmetryTransformer) rotatePoint(point Point, angle float64) Point {
	sin := math.Sin(angle)
	cos := math.Cos(angle)

	return Point{
		X: point.X*cos - point.Y*sin,
		Y: point.X*sin + point.Y*cos,
	}
}

func (st *SymmetryTransformer) ApplyWithColor(point Point, color Color) []struct {
	Point Point
	Color Color
} {
	if st.level == 1 {
		return []struct {
			Point Point
			Color Color
		}{
			{Point: point, Color: color},
		}
	}

	result := make([]struct {
		Point Point
		Color Color
	}, st.level)

	result[0] = struct {
		Point Point
		Color Color
	}{Point: point, Color: color}

	for i := 1; i < st.level; i++ {
		angle := st.angle * float64(i)
		rotatedPoint := st.rotatePoint(point, angle)
		result[i] = struct {
			Point Point
			Color Color
		}{Point: rotatedPoint, Color: color}
	}

	return result
}

func (st *SymmetryTransformer) GetLevel() int {
	return st.level
}
