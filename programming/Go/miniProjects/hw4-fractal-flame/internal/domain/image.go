package domain

type ImageBuffer struct {
	Width     int
	Height    int
	Data      []Color
	Hits      []int
	LogScaled bool
}

func NewImageBuffer(width, height int, logScaled bool) *ImageBuffer {
	return &ImageBuffer{
		Width:     width,
		Height:    height,
		Data:      make([]Color, width*height),
		Hits:      make([]int, width*height),
		LogScaled: logScaled,
	}
}

func (ib *ImageBuffer) AddPoint(x, y int, color Color, symmetry *SymmetryTransformer) {
	if symmetry == nil || symmetry.GetLevel() == 1 {
		ib.addSinglePoint(x, y, color)
		return
	}

	nx := (float64(x)/float64(ib.Width-1))*2 - 1
	ny := (float64(y)/float64(ib.Height-1))*2 - 1
	point := Point{X: nx, Y: ny}
	symPoints := symmetry.ApplyWithColor(point, color)
	for _, sp := range symPoints {
		sx := int((sp.Point.X + 1) * 0.5 * float64(ib.Width-1))
		sy := int((sp.Point.Y + 1) * 0.5 * float64(ib.Height-1))
		ib.addSinglePoint(sx, sy, sp.Color)
	}
}

func (ib *ImageBuffer) Normalize() {
	for i := 0; i < len(ib.Data); i++ {
		if ib.Hits[i] > 0 {
			ib.Data[i].R /= float64(ib.Hits[i])
			ib.Data[i].G /= float64(ib.Hits[i])
			ib.Data[i].B /= float64(ib.Hits[i])
			ib.Data[i].A = 1.0
		}
	}
}

func (ib *ImageBuffer) addSinglePoint(x, y int, color Color) {
	if x >= 0 && x < ib.Width && y >= 0 && y < ib.Height {
		idx := y*ib.Width + x
		ib.Hits[idx]++

		ib.Data[idx].R += color.R
		ib.Data[idx].G += color.G
		ib.Data[idx].B += color.B
	}
}
