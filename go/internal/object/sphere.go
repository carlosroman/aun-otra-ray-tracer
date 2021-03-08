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

func NewSphere(center ray.Vector, radius float64, opts ...Option) Object {
	s := &sphere{
		c: center,
		r: radius,
	}
	_ = s.SetTransform(ray.DefaultIdentityMatrix())
	s.SetMaterial(DefaultMaterial())

	for i := range opts {
		opts[i].Apply(s)
	}
	return s
}

func NewGlassSphere(center ray.Vector, radius float64, opts ...Option) (s Object) {
	material := DefaultMaterial()
	material.Transparency = 1.0
	material.RefractiveIndex = 1.5
	opts = append(opts, WithMaterial(material))
	s = NewSphere(center, radius,
		opts...,
	)
	return s
}

func DefaultSphere() Object {
	return NewSphere(ray.NewPoint(0, 0, 0), 1)
}

func DefaultGlassSphere() Object {
	return NewGlassSphere(ray.NewPoint(0, 0, 0), 1)
}

func NewMetalSphere(center ray.Vector, radius float64, opts ...Option) (s Object) {
	material := DefaultMaterial()
	material.Diffuse = 0.6
	material.Reflective = 0.1
	material.Specular = 0.4
	material.Shininess = 10
	material.Color = NewColor(0.9, 1, 0.9)
	opts = append(opts, WithMaterial(material))
	s = NewSphere(center, radius,
		opts...,
	)
	return s
}

func DefaultMetalSphere() Object {
	return NewMetalSphere(ray.NewPoint(0, 0, 0), 1)
}
