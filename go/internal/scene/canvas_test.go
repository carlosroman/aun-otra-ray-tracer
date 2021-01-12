package scene_test

import (
	"image/color"
	"testing"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/scene"
	"github.com/stretchr/testify/assert"
)

func TestNewCanvas(t *testing.T) {
	// Given
	c := scene.NewCanvas(5, 3)

	// And
	c1 := object.NewColor(1.5, 0, 0)
	c2 := object.NewColor(0, 0.5, 0)
	c3 := object.NewColor(-0.5, 0, 1)

	// When
	c.SetColor(0, 0, c1)
	c.SetColor(2, 1, c2)
	c.Set(4, 2, c3.R, c3.G, c3.B)

	// Then
	img := c.GenerateImg()

	assert.Equal(t, color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 0,
	}, img.At(0, 0))

	assert.Equal(t, color.RGBA{
		R: 0,
		G: 128,
		B: 0,
		A: 0,
	}, img.At(2, 1))

	assert.Equal(t, color.RGBA{
		R: 0,
		G: 0,
		B: 255,
		A: 0,
	}, img.At(4, 2))
}
