package object

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type Object interface {
	Intersect(ray ray.Ray) []float64
	NormalAt(worldPoint ray.Vector) ray.Vector
	Transform() ray.Matrix
	SetTransform(t ray.Matrix)
}

type sphere struct {
	c ray.Vector
	r float64
	t ray.Matrix
}

func (s sphere) NormalAt(worldPoint ray.Vector) ray.Vector {
	inv, err := s.t.Inverse()
	if err != nil {
		return nil
	}
	objPt := inv.MultiplyByVector(worldPoint)
	objN := objPt.Subtract(s.c)
	return inv.
		Transpose().
		MultiplyByVector(objN).
		SetW(0).
		Normalize()
}

func (s sphere) Intersect(rr ray.Ray) []float64 {
	inv, err := s.t.Inverse()
	if err != nil {
		return nil
	}
	r := rr.Transform(inv)
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

func NewSphere(center ray.Vector, radius float64) Object {
	return &sphere{
		c: center,
		r: radius,
		t: ray.DefaultIdentityMatrix(),
	}
}
