package object_test

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
)

func TestDefaultMaterial(t *testing.T) {
	material := object.DefaultMaterial()

	assert.Equal(t, 0.1, material.Ambient)
	assert.Equal(t, 0.9, material.Diffuse)
	assert.Equal(t, 0.9, material.Specular)
	assert.Equal(t, 200.0, material.Shininess)
	assert.Equal(t, color.RGBA{R: 1, G: 1, B: 1, A: 0}, material.Color)
}
