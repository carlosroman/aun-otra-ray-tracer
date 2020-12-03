package object

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type Object interface {
	Intersect(ray ray.Ray) []float64
}

type sphere struct {
	c ray.Vector
	r float64
}

func (s sphere) Intersect(r ray.Ray) []float64 {
	sphereToRay := r.Origin().Subtract(s.c)
	a := ray.Dot(r.Direction(), r.Direction())
	b := 2 * ray.Dot(sphereToRay, r.Direction())
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

func NewSphere(center ray.Vector, radius float64) Object {
	return &sphere{
		c: center,
		r: radius,
	}
}
