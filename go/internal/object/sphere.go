package object

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type sphere struct {
	c ray.Vector
	r float64
	t ray.Matrix
	m Material
}

func NormalAt(obj Object, worldPoint ray.Vector) ray.Vector {
	inv, err := obj.Transform().Inverse()
	if err != nil {
		return nil
	}
	localPoint := inv.MultiplyByVector(worldPoint)
	localNormal := obj.LocalNormalAt(localPoint)
	return inv.
		Transpose().
		MultiplyByVector(localNormal).
		SetW(0).
		Normalize()
}

func (s sphere) LocalNormalAt(worldPoint ray.Vector) ray.Vector {
	return worldPoint.Subtract(s.c)
}

func Intersect(obj Object, rr ray.Ray) []float64 {
	inv, err := obj.Transform().Inverse()
	if err != nil {
		return nil
	}
	return obj.LocalIntersect(rr.Transform(inv))
}

func (s *sphere) LocalIntersect(r ray.Ray) []float64 {
	sphereToRay := r.Origin().Subtract(s.c)
	a := ray.Dot(r.Direction(), r.Direction())
	b := 2 * ray.Dot(r.Direction(), sphereToRay)
	c := ray.Dot(sphereToRay, sphereToRay) - (s.r * s.r)
	discriminant := (b * b) - (4 * a * c)
	if discriminant < 0 {
		return nil
	}

	f := 2 * a
	nB := -b
	sqrtDisc := math.Sqrt(discriminant)
	t1 := (nB - sqrtDisc) / f
	t2 := (nB + sqrtDisc) / f
	return []float64{t1, t2}
}

func (s sphere) Transform() ray.Matrix {
	return s.t
}

func (s *sphere) SetTransform(t ray.Matrix) {
	s.t = t
}

func (s sphere) Material() Material {
	return s.m
}

func (s *sphere) SetMaterial(m Material) {
	s.m = m
}

func NewSphere(center ray.Vector, radius float64) Object {
	return &sphere{
		c: center,
		r: radius,
		t: ray.DefaultIdentityMatrix(),
		m: DefaultMaterial(),
	}
}

func NewGlassSphere(center ray.Vector, radius float64) Object {
	material := DefaultMaterial()
	material.Transparency = 1.0
	material.RefractiveIndex = 1.5
	return &sphere{
		c: center,
		r: radius,
		t: ray.DefaultIdentityMatrix(),
		m: material,
	}
}

func DefaultSphere() Object {
	return NewSphere(ray.NewPoint(0, 0, 0), 1)
}

func DefaultGlassSphere() Object {
	return NewGlassSphere(ray.NewPoint(0, 0, 0), 1)
}
