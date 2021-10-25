package object

import (
	"fmt"
	"math"

	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/ray"
)

type triangle struct {
	obj
	p1, p2, p3 ray.Vector
	n1, n2, n3 ray.Vector
	e1, e2     ray.Vector
}

func WithNormals(n1, n2, n3 ray.Vector) Option {
	return OptionFunc(func(o Object) {
		if t, ok := o.(*triangle); ok {
			t.n1 = n1
			t.n2 = n2
			t.n3 = n3
		}
	})
}

func NewSmoothTriangle(p1, p2, p3, n1, n2, n3 ray.Vector, opts ...Option) Object {
	opts = append(opts, WithNormals(n1, n2, n3))
	return NewTriangle(p1, p2, p3, opts...)
}

func NewTriangle(p1, p2, p3 ray.Vector, opts ...Option) Object {
	t := triangle{
		p1: p1,
		p2: p2,
		p3: p3,
		e1: p2.Subtract(p1),
		e2: p3.Subtract(p1),
	}
	normal := ray.Cross(t.e2, t.e1).Normalize()
	t.n1 = normal
	t.n2 = normal
	t.n3 = normal

	for i := range opts {
		opts[i].Apply(&t)
	}

	_ = t.SetTransform(ray.DefaultIdentityMatrix())
	t.SetMaterial(DefaultMaterial())
	return &t
}

func (t triangle) String() string {
	return fmt.Sprintf("e1: %v, e2: %v, n1: %v, n2: %v, n3: %v", t.e1, t.e2, t.n1, t.n2, t.n3)
}

func (t triangle) LocalNormalAt(_ ray.Vector, hit Intersection) ray.Vector {
	a := t.n2.Multiply(hit.U)
	b := t.n3.Multiply(hit.V)
	c := t.n1.Multiply(1 - hit.U - hit.V)
	return a.Add(b).Add(c)
}

func (t triangle) LocalIntersect(r ray.Ray) Intersections {

	dirCrossE2 := ray.Cross(r.Direction(), t.e2)
	det := ray.Dot(t.e1, dirCrossE2)
	if math.Abs(det) < epsilon {
		return nil
	}

	f := 1.0 / det
	p1ToOrig := r.Origin().Subtract(t.p1)
	u := f * ray.Dot(p1ToOrig, dirCrossE2)
	if u < 0 || u > 1 {
		return nil
	}

	origCrossE1 := ray.Cross(p1ToOrig, t.e1)
	v := f * ray.Dot(r.Direction(), origCrossE1)
	if v < 0 || (u+v) > 1 {
		return nil
	}

	return Intersections{Intersection{
		T:   f * ray.Dot(t.e2, origCrossE1),
		Obj: &t,
		U:   u,
		V:   v,
	}}
}
