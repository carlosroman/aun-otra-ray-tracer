package object_test

import (
	"math"
	"testing"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
	"github.com/stretchr/testify/assert"
)

func TestNewPointLight(t *testing.T) {
	expectedIntensity := object.NewColor(1, 1, 1)
	expectedPosition := ray.NewPoint(0, 0, 0)
	actual := object.NewPointLight(expectedPosition, expectedIntensity)

	assert.Equal(t, expectedIntensity, actual.Intensity)
	assert.Equal(t, expectedPosition, actual.Position)
}

func TestLighting(t *testing.T) {

	testCases := []struct {
		name          string
		eyev, normalv ray.Vector
		light         object.PointLight
		inShadows     bool
		expectedColor object.RGB
	}{
		{
			name:          "Lighting with the eye between the light and the surface",
			eyev:          ray.NewVec(0, 0, -1),
			normalv:       ray.NewVec(0, 0, -1),
			light:         object.NewPointLight(ray.NewPoint(0, 0, -10), object.NewColor(1, 1, 1)),
			expectedColor: object.RGB{R: 1.9, B: 1.9, G: 1.9},
		},
		{
			name:          "Lighting with the eye between light and surface, eye offset 45°",
			eyev:          ray.NewVec(0, math.Sqrt(2)/2, math.Sqrt(2)/2),
			normalv:       ray.NewVec(0, 0, -1),
			light:         object.NewPointLight(ray.NewPoint(0, 0, -10), object.NewColor(1, 1, 1)),
			expectedColor: object.RGB{R: 1.0, B: 1.0, G: 1.0},
		},
		{
			name:          "Lighting with eye opposite surface, light offset 45°",
			eyev:          ray.NewVec(0, 0, -1),
			normalv:       ray.NewVec(0, 0, -1),
			light:         object.NewPointLight(ray.NewPoint(0, 10, -10), object.NewColor(1, 1, 1)),
			expectedColor: object.RGB{R: 0.7364, B: 0.7364, G: 0.7364},
		},
		{
			name:          "Lighting with eye in the path of the reflection vector",
			eyev:          ray.NewVec(0, -math.Sqrt(2)/2, -math.Sqrt(2)/2),
			normalv:       ray.NewVec(0, 0, -1),
			light:         object.NewPointLight(ray.NewPoint(0, 10, -10), object.NewColor(1, 1, 1)),
			expectedColor: object.RGB{R: 1.6364, B: 1.6364, G: 1.6364},
		},
		{
			name:          "Lighting with the light behind the surface",
			eyev:          ray.NewVec(0, 0, -1),
			normalv:       ray.NewVec(0, 0, -1),
			light:         object.NewPointLight(ray.NewPoint(0, 0, 10), object.NewColor(1, 1, 1)),
			expectedColor: object.RGB{R: 0.1, B: 0.1, G: 0.1},
		},
		{
			name:          "Lighting with the surface in shadow",
			eyev:          ray.NewVec(0, 0, -1),
			normalv:       ray.NewVec(0, 0, -1),
			inShadows:     true,
			light:         object.NewPointLight(ray.NewPoint(0, 0, -10), object.NewColor(1, 1, 1)),
			expectedColor: object.RGB{R: 0.1, B: 0.1, G: 0.1},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			m := object.DefaultMaterial()
			obj := object.NewSphere(ray.ZeroPoint, 1)
			position := ray.NewPoint(0, 0, 0)

			// When
			actual := object.Lighting(m, obj, tt.light, position, tt.eyev, tt.normalv, tt.inShadows)

			// Then
			assertColorEqual(t, tt.expectedColor, actual)
		})
	}
}

func TestLighting_WithPattern(t *testing.T) {
	m := object.DefaultMaterial()
	m.Ambient = 1
	m.Diffuse = 0
	m.Specular = 0
	m.Pattern = object.NewStripePattern(object.NewColor(1, 1, 1), object.NewColor(0, 0, 0))

	eyev := ray.NewVec(0, 0, -1)
	normalv := ray.NewVec(0, 0, -1)
	light := object.NewPointLight(ray.NewPoint(0, 0, -10), object.NewColor(1, 1, 1))

	obj := object.NewSphere(ray.ZeroPoint, 1)
	c1 := object.Lighting(m, obj, light, ray.NewPoint(0.9, 0, 0), eyev, normalv, false)
	c2 := object.Lighting(m, obj, light, ray.NewPoint(1.1, 0, 0), eyev, normalv, false)

	assertColorEqual(t, object.NewColor(1, 1, 1), c1)
	assertColorEqual(t, object.NewColor(0, 0, 0), c2)
}

func assertColorEqual(t *testing.T, expected object.RGB, actual object.RGB) {
	assert.InDelta(t, expected.R, actual.R, 0.00001, "R")
	assert.InDelta(t, expected.G, actual.G, 0.00001, "G")
	assert.InDelta(t, expected.B, actual.B, 0.00001, "B")
}
