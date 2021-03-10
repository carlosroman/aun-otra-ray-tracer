package object

import (
	"fmt"
	"math"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type triangle struct {
	obj
	p1, p2, p3 ray.Vector
	e1, e2     ray.Vector
	normal     ray.Vector
}

func NewTriangle(p1, p2, p3 ray.Vector) Object {
	t := triangle{
		p1: p1,
		p2: p2,
		p3: p3,
		e1: p2.Subtract(p1),
		e2: p3.Subtract(p1),
	}
	t.normal = ray.Cross(t.e2, t.e1).Normalize()

	_ = t.SetTransform(ray.DefaultIdentityMatrix())
	t.SetMaterial(DefaultMaterial())
	return &t
}

func (t triangle) String() string {
	return fmt.Sprintf("e1: %v, e2: %v, normal: %v", t.e1, t.e2, t.normal)
}

func (t triangle) LocalNormalAt(_ ray.Vector, _ Intersection) ray.Vector {
	return t.normal
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
	}}
}
