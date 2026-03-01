package domain

import (
	"testing"
)

var defaultSymmetry = NewSymmetryTransformer(1)

func TestNewImageBuffer(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"small", 10, 10},
		{"medium", 100, 100},
		{"rectangular", 200, 100},
		{"large", 1920, 1080},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := NewImageBuffer(tt.width, tt.height, false)

			if buffer.Width != tt.width {
				t.Errorf("Width = %v, want %v", buffer.Width, tt.width)
			}

			if buffer.Height != tt.height {
				t.Errorf("Height = %v, want %v", buffer.Height, tt.height)
			}

			expectedSize := tt.width * tt.height
			if len(buffer.Data) != expectedSize {
				t.Errorf("Data length = %v, want %v", len(buffer.Data), expectedSize)
			}

			if len(buffer.Hits) != expectedSize {
				t.Errorf("Hits length = %v, want %v", len(buffer.Hits), expectedSize)
			}

			for i := 0; i < expectedSize; i++ {
				if buffer.Hits[i] != 0 {
					t.Errorf("Hits[%d] = %v, want 0", i, buffer.Hits[i])
				}
			}
		})
	}
}

func TestImageBuffer_AddPoint(t *testing.T) {
	buffer := NewImageBuffer(10, 10, false)
	color := Color{R: 0.5, G: 0.3, B: 0.2, A: 0.0}

	tests := []struct {
		name     string
		x, y     int
		color    Color
		expected struct {
			hits  int
			color Color
		}
	}{
		{
			name: "add single point",
			x:    5, y: 5,
			color: color,
			expected: struct {
				hits  int
				color Color
			}{hits: 1, color: color},
		},
		{
			name: "add same point twice",
			x:    5, y: 5,
			color: Color{R: 0.1, G: 0.2, B: 0.3, A: 1.0},
			expected: struct {
				hits  int
				color Color
			}{hits: 2, color: Color{R: 0.6, G: 0.5, B: 0.5, A: 0.0}},
		},
		{
			name: "add different point",
			x:    3, y: 7,
			color: Color{R: 0.8, G: 0.1, B: 0.1, A: 1.0},
			expected: struct {
				hits  int
				color Color
			}{hits: 1, color: Color{R: 0.8, G: 0.1, B: 0.1, A: 0.0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer.AddPoint(tt.x, tt.y, tt.color, defaultSymmetry)

			idx := tt.y*buffer.Width + tt.x
			if buffer.Hits[idx] != tt.expected.hits {
				t.Errorf("Hits[%d] = %v, want %v", idx, buffer.Hits[idx], tt.expected.hits)
			}

			if buffer.Data[idx].R != tt.expected.color.R ||
				buffer.Data[idx].G != tt.expected.color.G ||
				buffer.Data[idx].B != tt.expected.color.B ||
				buffer.Data[idx].A != tt.expected.color.A {
				t.Errorf("Color = %v, want %v", buffer.Data[idx], tt.expected.color)
			}
		})
	}
}

func TestImageBuffer_AddPointOutOfBounds(t *testing.T) {
	buffer := NewImageBuffer(5, 5, false)
	color := Color{R: 1, G: 1, B: 1, A: 1}

	outOfBounds := []struct {
		x, y int
		name string
	}{
		{-1, 0, "negative x"},
		{0, -1, "negative y"},
		{5, 0, "x equals width"},
		{0, 5, "y equals height"},
		{10, 10, "both out"},
	}

	for _, tt := range outOfBounds {
		t.Run(tt.name, func(t *testing.T) {

			hitsBefore := make([]int, len(buffer.Hits))
			copy(hitsBefore, buffer.Hits)

			buffer.AddPoint(tt.x, tt.y, color, defaultSymmetry)

			for i := range buffer.Hits {
				if buffer.Hits[i] != hitsBefore[i] {
					t.Errorf("Hits changed for out of bounds point at (%d, %d)", tt.x, tt.y)
					break
				}
			}
		})
	}
}

func TestImageBuffer_Normalize(t *testing.T) {
	tests := []struct {
		name      string
		logScaled bool
		setupFunc func(*ImageBuffer)
		checkFunc func(*ImageBuffer) bool
	}{
		{
			name:      "linear normalization empty buffer",
			logScaled: false,
			setupFunc: func(b *ImageBuffer) {
			},
			checkFunc: func(b *ImageBuffer) bool {
				for i := range b.Data {
					if b.Data[i].R != 0 || b.Data[i].G != 0 || b.Data[i].B != 0 {
						return false
					}
				}
				return true
			},
		},
		{
			name:      "linear normalization single hit",
			logScaled: false,
			setupFunc: func(b *ImageBuffer) {
				b.AddPoint(0, 0, Color{R: 0.5, G: 0.3, B: 0.2, A: 1.0}, defaultSymmetry)
			},
			checkFunc: func(b *ImageBuffer) bool {
				color := b.Data[0]
				return color.R == 0.5 && color.G == 0.3 && color.B == 0.2 && color.A == 1.0
			},
		},
		{
			name:      "linear normalization multiple hits",
			logScaled: false,
			setupFunc: func(b *ImageBuffer) {
				b.AddPoint(0, 0, Color{R: 0.6, G: 0.4, B: 0.2, A: 1.0}, defaultSymmetry)
				b.AddPoint(0, 0, Color{R: 0.2, G: 0.4, B: 0.6, A: 1.0}, defaultSymmetry)
			},
			checkFunc: func(b *ImageBuffer) bool {
				color := b.Data[0]
				return color.R == 0.4 && color.G == 0.4 && color.B == 0.4 && color.A == 1.0
			},
		},
		{
			name:      "log normalization",
			logScaled: true,
			setupFunc: func(b *ImageBuffer) {
				b.AddPoint(0, 0, Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0}, defaultSymmetry)
				b.AddPoint(1, 0, Color{R: 0.5, G: 0.5, B: 0.5, A: 1.0}, defaultSymmetry)
			},
			checkFunc: func(b *ImageBuffer) bool {
				color1 := b.Data[0]
				color2 := b.Data[1]
				return color1.R == 1.0 && color2.R < 1.0 && color1.R > color2.R
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := NewImageBuffer(5, 5, tt.logScaled)
			tt.setupFunc(buffer)

			buffer.Normalize()

			if !tt.checkFunc(buffer) {
				t.Errorf("Normalization failed for test: %s", tt.name)
			}
		})
	}
}
