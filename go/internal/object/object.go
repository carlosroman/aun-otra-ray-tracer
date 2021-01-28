package object

import (
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type Object interface {
	LocalIntersect(ray ray.Ray) []float64
	LocalNormalAt(worldPoint ray.Vector) ray.Vector

	Transform() ray.Matrix
	SetTransform(t ray.Matrix)

	Material() (m Material)
	SetMaterial(m Material)
}
