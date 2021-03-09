package object_test

import (
	"math"
	"testing"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
	"github.com/stretchr/testify/assert"
)

func TestWithMaterial(t *testing.T) {
	obj := object.NewTestShape()
	expectedMaterial := object.DefaultMaterial()
	expectedMaterial.Diffuse = 1.11
	opt := object.WithMaterial(expectedMaterial)
	opt.Apply(obj)
	assert.Equal(t, expectedMaterial, obj.Material())
}

func TestWithTransform(t *testing.T) {
	obj := object.NewTestShape()
	expectedTransformation := ray.Rotation(ray.Y, math.Pi/2)
	opt := object.WithTransform(expectedTransformation)
	opt.Apply(obj)
	assert.Equal(t, expectedTransformation, obj.Transform())
}
