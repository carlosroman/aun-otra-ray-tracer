package object

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

var (
	defaultPlaneLocalNormal = ray.NewVec(0, 1, 0)
)

type plane struct {
	t ray.Matrix
	m Material
}

func (p plane) Transform() ray.Matrix {
	return p.t
}

func (p *plane) SetTransform(t ray.Matrix) {
	p.t = t
}

func (p plane) Material() Material {
	return p.m
}

func (p *plane) SetMaterial(m Material) {
	p.m = m
}

func (p plane) LocalIntersect(ray ray.Ray) []float64 {
	if math.Abs(ray.Direction().GetY()) < epsilon {
		return nil
	}

	t := -ray.Origin().GetY() / ray.Direction().GetY()
	return []float64{t}
}

func (p plane) LocalNormalAt(_ ray.Vector) ray.Vector {
	return defaultPlaneLocalNormal
}

func NewPlane() Object {
	return &plane{
		t: ray.DefaultIdentityMatrix(),
		m: DefaultMaterial(),
	}
}
