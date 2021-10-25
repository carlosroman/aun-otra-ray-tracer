package scene

import (
	"image"
	"image/color"
	"math"

	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/object"
)

type Canvas [][]object.RGB

func NewCanvas(width, height int) (c Canvas) {
	c = make(Canvas, width)
	for i := range c {
		c[i] = make([]object.RGB, height)
	}
	return c
}

func (c Canvas) Get(x, y int) object.RGB {
	return c[x][y]
}

func (c Canvas) Set(x, y int, red, green, blue float64) {
	c[x][y] = object.NewColor(red, green, blue)
}

func (c Canvas) SetColor(x, y int, col object.RGB) {
	c[x][y] = col
}

func (c Canvas) GenerateImg() (img *image.RGBA) {
	const maxC = 255
	const min = 0.0
	const max = 1
	const alpha = uint8(0)
	img = image.NewRGBA(
		image.Rect(0, 0, len(c), len(c[0])),
	)
	for x := 0; x < len(c); x++ {
		for y := 0; y < len(c[x]); y++ {
			img.Set(x, y, color.RGBA{
				R: uint8(math.Round(clamp(c[x][y].R, min, max) * maxC)),
				G: uint8(math.Round(clamp(c[x][y].G, min, max) * maxC)),
				B: uint8(math.Round(clamp(c[x][y].B, min, max) * maxC)),
				A: alpha,
			})
		}
	}
	return img
}

func clamp(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
