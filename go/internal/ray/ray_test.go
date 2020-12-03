package ray_test

import (
	"fmt"
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
	assert.Equal(t, float64(1), newRay.Origin().GetX())
	assert.Equal(t, float64(2), newRay.Origin().GetY())
	assert.Equal(t, float64(3), newRay.Origin().GetZ())
}

func TestRay_Direction(t *testing.T) {
	newRayAt := ray.NewRayAt(ray.NewVec(2, 3, 4), ray.NewVec(1, 0, 0))
	assertVec(t, ray.NewVec(1, 0, 0), newRayAt.Direction())
}

func TestRay_Origin(t *testing.T) {
	newRayAt := ray.NewRayAt(ray.NewVec(2, 3, 4), ray.NewVec(1, 0, 0))
	assertVec(t, ray.NewVec(2, 3, 4), newRayAt.Origin())

}

func TestRay_PointAt(t *testing.T) {
	newRayAt := ray.NewRayAt(ray.NewVec(2, 3, 4), ray.NewVec(1, 0, 0))

	testCases := []struct {
		tick     float64
		expected ray.Vector
	}{
		{
			tick:     0,
			expected: ray.NewVec(2, 3, 4),
		},
		{
			tick:     1,
			expected: ray.NewVec(3, 3, 4),
		},
		{
			tick:     -1,
			expected: ray.NewVec(1, 3, 4),
		},
		{
			tick:     2.5,
			expected: ray.NewVec(4.5, 3, 4),
		},
	}
	for _, tt := range testCases {
		name := fmt.Sprintf("tick_%v", tt.tick)
		t.Run(name, func(t *testing.T) {
			position := newRayAt.PointAt(tt.tick)
			assertVec(t, tt.expected, position)
		})
	}
}