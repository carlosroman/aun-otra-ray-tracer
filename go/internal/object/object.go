package object

import (
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

const (
	epsilon = 0.00000001
)

type Object interface {
	LocalIntersect(ray ray.Ray) []float64
	LocalNormalAt(worldPoint ray.Vector) ray.Vector

	Transform() ray.Matrix
	TransformInverse() ray.Matrix
	SetTransform(t ray.Matrix) error

	Material() (m Material)
	SetMaterial(m Material)
}

type BasicObject struct {
	Transform ray.Matrix
	Material  Material
}
