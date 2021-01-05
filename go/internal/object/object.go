package object

import (
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type Object interface {
	Intersect(ray ray.Ray) []float64
	NormalAt(worldPoint ray.Vector) ray.Vector
	Transform() ray.Matrix
	SetTransform(t ray.Matrix)
}
