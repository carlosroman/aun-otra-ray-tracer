package object

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type sphere struct {
	obj
	c ray.Vector
	r float64
}

func NormalAt(obj Object, worldPoint ray.Vector) ray.Vector {
	localPoint := obj.WorldToObject(worldPoint)
	localNormal := obj.LocalNormalAt(localPoint)
	return obj.NormalToWorld(localNormal)
}

func (s sphere) LocalNormalAt(worldPoint ray.Vector) ray.Vector {
	return worldPoint.Subtract(s.c)
}

func Intersect(obj Object, rr ray.Ray) Intersections {
	return obj.
		LocalIntersect(
			rr.Transform(obj.TransformInverse()))
}

func (s *sphere) LocalIntersect(r ray.Ray) Intersections {
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
	return []Intersection{
		{T: t1, Obj: s},
		{T: t2, Obj: s},
	}
}

func NewSphere(center ray.Vector, radius float64) Object {
	s := &sphere{
		c: center,
		r: radius,
	}
	_ = s.SetTransform(ray.DefaultIdentityMatrix())
	s.SetMaterial(DefaultMaterial())
	return s
}

func NewGlassSphere(center ray.Vector, radius float64) Object {
	material := DefaultMaterial()
	material.Transparency = 1.0
	material.RefractiveIndex = 1.5
	s := NewSphere(center, radius)
	s.SetMaterial(material)
	return s
}

func DefaultSphere() Object {
	return NewSphere(ray.NewPoint(0, 0, 0), 1)
}

func DefaultGlassSphere() Object {
	return NewGlassSphere(ray.NewPoint(0, 0, 0), 1)
}
