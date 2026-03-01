package application

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"

	"fractalflame/internal/domain"
	"fractalflame/internal/infrastructure"
)

type weightedFunc struct {
	fn     domain.TransformationFunction
	weight float64
}

type Generator struct {
	logger *infrastructure.Logger
}

func NewGenerator(logger *infrastructure.Logger) *Generator {
	return &Generator{logger: logger}
}

func (g *Generator) Generate(config *domain.Config) (*image.RGBA, error) {
	rnd := rand.New(rand.NewSource(int64(config.Seed)))
	symmetry := domain.NewSymmetryTransformer(config.SymmetryLevel)

	buffer := domain.NewImageBuffer(config.Size.Width, config.Size.Height, config.GammaCorrection)
	transforms := g.createAffineTransforms(config, rnd)

	functions := g.createFunctions(config)
	normalizedFuncs := g.normalizeFunctionWeights(functions, config.Functions)

	progress := infrastructure.NewProgress(config.IterationCount, g.logger)

	if config.Threads == 1 {
		g.generateSingleThread(config, buffer, transforms, normalizedFuncs, progress, rnd, symmetry)
	} else {
		g.generateMultiThread(config, buffer, transforms, normalizedFuncs, progress, symmetry)
	}

	buffer.Normalize()
	return g.convertToImage(buffer, config.Gamma), nil
}

func (g *Generator) createAffineTransforms(config *domain.Config, random *rand.Rand) []domain.AffineTransform {
	var transforms []domain.AffineTransform

	for i, params := range config.AffineParams {
		if i == 0 && params.A == 0 && params.B == 0 {
			for j := 0; j < 3; j++ {
				transforms = append(transforms, domain.NewRandomAffine(config.Seed+float64(j), random))
			}
			break
		}

		transforms = append(transforms, domain.AffineTransform{
			Params: params,
			Color: domain.Color{
				R: 0.2 + 0.6*random.Float64(),
				G: 0.2 + 0.6*random.Float64(),
				B: 0.2 + 0.6*random.Float64(),
				A: 1.0,
			},
		})
	}

	return transforms
}

func (g *Generator) createFunctions(config *domain.Config) []weightedFunc {
	var functions []weightedFunc
	for _, fc := range config.Functions {
		f, err := domain.GetFunctionByName(fc.Name)
		if err != nil {
			g.logger.Warn(fmt.Sprintf("Skipping unknown function: %s", fc.Name))
			continue
		}

		functions = append(functions, weightedFunc{fn: f, weight: fc.Weight})
	}

	return functions
}

func (g *Generator) normalizeFunctionWeights(functions []weightedFunc, configFuncs []domain.FunctionConfig) []weightedFunc {
	totalWeight := 0.0
	for _, f := range configFuncs {
		totalWeight += f.Weight
	}
	normalized := make([]weightedFunc, len(functions))

	for i, f := range functions {
		normalized[i] = weightedFunc{fn: f.fn, weight: f.weight / totalWeight}
	}
	return normalized
}

func (g *Generator) generateSingleThread(
	config *domain.Config,
	buffer *domain.ImageBuffer,
	transforms []domain.AffineTransform,
	functions []weightedFunc,
	progress *infrastructure.Progress,
	random *rand.Rand,
	symetry *domain.SymmetryTransformer,
) {
	point := domain.Point{X: 0.1, Y: 0.1}
	for i := 0; i < config.IterationCount; i++ {
		transform := transforms[random.Intn(len(transforms))]
		point = transform.Apply(point)
		funcIdx := g.selectFunctionIndex(functions)
		if funcIdx >= 0 {
			point = functions[funcIdx].fn.Apply(point)
		}
		x, y := g.scaleToImage(point, config.Size.Width, config.Size.Height)
		buffer.AddPoint(x, y, transform.Color, symetry)
		progress.Update()
	}
}

func (g *Generator) generateMultiThread(
	config *domain.Config,
	buffer *domain.ImageBuffer,
	transforms []domain.AffineTransform,
	functions []weightedFunc,
	progress *infrastructure.Progress,
	symetry *domain.SymmetryTransformer,
) {
	var wg sync.WaitGroup
	iterationsPerThread := config.IterationCount / config.Threads

	var completed int64

	for t := 0; t < config.Threads; t++ {
		wg.Add(1)

		go func(threadID int) {
			defer wg.Done()

			rng := rand.New(rand.NewSource(int64(config.Seed) + int64(threadID)))

			point := domain.Point{
				X: 0.1 + 0.1*float64(threadID),
				Y: 0.1 + 0.1*float64(threadID),
			}

			localBuffer := domain.NewImageBuffer(config.Size.Width, config.Size.Height, config.GammaCorrection)

			for i := 0; i < iterationsPerThread; i++ {

				transform := transforms[rng.Intn(len(transforms))]

				point = transform.Apply(point)

				funcIdx := g.selectFunctionIndexWithRNG(functions, rng)
				if funcIdx >= 0 {
					point = functions[funcIdx].fn.Apply(point)
				}

				x, y := g.scaleToImage(point, config.Size.Width, config.Size.Height)

				localBuffer.AddPoint(x, y, transform.Color, symetry)

				if i%1000 == 0 {
					atomic.AddInt64(&completed, 1000)
					progress.Update(atomic.LoadInt64(&completed))
				}
			}

			g.mergeBuffers(buffer, localBuffer)

		}(t)
	}

	wg.Wait()
	progress.Finish()
}

func (g *Generator) scaleToImage(point domain.Point, width, height int) (int, int) {
	scale := 3.0
	offset := 0.5

	normalizedX := (point.X/scale + offset)
	normalizedY := (point.Y/scale + offset)

	normalizedX = g.clamp(normalizedX, 0.0, 1.0)
	normalizedY = g.clamp(normalizedY, 0.0, 1.0)

	pixelX := int(normalizedX * float64(width-1))
	pixelY := int(normalizedY * float64(height-1))

	return pixelX, pixelY
}

func (g *Generator) clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func (g *Generator) selectFunctionIndex(functions []weightedFunc) int {
	r := rand.Float64()
	cumulative := 0.0

	for i, f := range functions {
		cumulative += f.weight
		if r <= cumulative {
			return i
		}
	}

	return -1
}

func (g *Generator) selectFunctionIndexWithRNG(functions []weightedFunc, rng *rand.Rand) int {
	r := rng.Float64()
	cumulative := 0.0

	for i, f := range functions {
		cumulative += f.weight
		if r <= cumulative {
			return i
		}
	}

	return -1
}

func (g *Generator) mergeBuffers(mainBuffer, localBuffer *domain.ImageBuffer) {
	for i := 0; i < len(mainBuffer.Data); i++ {
		if localBuffer.Hits[i] > 0 {
			mainBuffer.Hits[i] += localBuffer.Hits[i]
			mainBuffer.Data[i].R += localBuffer.Data[i].R
			mainBuffer.Data[i].G += localBuffer.Data[i].G
			mainBuffer.Data[i].B += localBuffer.Data[i].B
		}
	}
}

func (ge *Generator) convertToImage(buffer *domain.ImageBuffer, gamma float64) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, buffer.Width, buffer.Height))

	for y := 0; y < buffer.Height; y++ {
		for x := 0; x < buffer.Width; x++ {
			idx := y*buffer.Width + x
			col := buffer.Data[idx]

			r := col.R
			g := col.G
			b := col.B

			if buffer.LogScaled && gamma > 0 {
				r = math.Pow(r, 1.0/gamma)
				g = math.Pow(g, 1.0/gamma)
				b = math.Pow(b, 1.0/gamma)
			}

			r = ge.clamp(r, 0.0, 1.0)
			g = ge.clamp(g, 0.0, 1.0)
			b = ge.clamp(b, 0.0, 1.0)

			img.Set(x, y, color.RGBA{
				R: uint8(r * 255),
				G: uint8(g * 255),
				B: uint8(b * 255),
				A: 255,
			})
		}
	}

	return img
}

func SaveImage(img *image.RGBA, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}
