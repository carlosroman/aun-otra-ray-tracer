package scene_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/ray"
	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/scene"
)

func TestDefaultWorld(t *testing.T) {

	light := object.NewPointLight(ray.NewPoint(-10, 10, -10), object.NewColor(1, 1, 1))

	s1 := object.NewSphere(ray.ZeroPoint, 1)
	m1 := object.DefaultMaterial()
	m1.Color = object.NewColor(0.8, 1.0, 0.6)
	m1.Diffuse = 0.7
	m1.Specular = 0.2
	s1.SetMaterial(m1)

	s2 := object.NewSphere(ray.ZeroPoint, 1)
	require.NoError(t, s2.SetTransform(ray.Scaling(0.5, 0.5, 0.5)))

	w, err := scene.DefaultWorld()
	require.NoError(t, err)

	assert.Len(t, w.Objects(), 2)
	assert.Contains(t, w.Objects(), s1)
	assert.Contains(t, w.Objects(), s2)
	assert.Equal(t, light, w.Light())
}
