package ray_test

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

func Test_ray_RGBA(t *testing.T) {
	newRay := ray.NewRay(1, 2, 3, 2, 101, 202)
	r, g, b, a := newRay.RGBA()
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, newRay)

	// Only till NewPPMOutput other types of images.
	actualR, actualG, actualB, actualA := img.At(0, 0).RGBA()
	assert.Equal(t, r, actualR)
	assert.Equal(t, g, actualG)
	assert.Equal(t, b, actualB)
	assert.Equal(t, a, actualA)
}

func Test_ray_set(t *testing.T) {
	newRay := ray.NewRay(1, 2, 3, 0, 0, 0)
	newRay.SetR(2)
	newRay.SetG(101)
	newRay.SetB(202)
	r, g, b, a := newRay.RGBA()
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, newRay)

	// Only till NewPPMOutput other types of images.
	actualR, actualG, actualB, actualA := img.At(0, 0).RGBA()
	assert.Equal(t, r, actualR)
	assert.Equal(t, g, actualG)
	assert.Equal(t, b, actualB)
	assert.Equal(t, a, actualA)
}

func Test_ray_get(t *testing.T) {
	newRay := ray.NewRay(1, 2, 3, 0, 0, 0)
	assert.Equal(t, float64(1), newRay.GetX())
	assert.Equal(t, float64(2), newRay.GetY())
	assert.Equal(t, float64(3), newRay.GetZ())

}
