package object

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

var (
	defaultPlaneLocalNormal = ray.NewVec(0, 1, 0)
)

type plane struct {
	obj
}

func (p plane) LocalIntersect(ray ray.Ray) Intersections {
	if math.Abs(ray.Direction().GetY()) < epsilon {
		return nil
	}

	t := -ray.Origin().GetY() / ray.Direction().GetY()
	return []Intersection{{
		T:   t,
		Obj: &p,
	}}
}

func (p plane) LocalNormalAt(worldPoint ray.Vector, hit Intersection) ray.Vector {
	return defaultPlaneLocalNormal
}

func NewPlane() Object {
	p := plane{}
	p.SetMaterial(DefaultMaterial())
	_ = p.SetTransform(ray.DefaultIdentityMatrix())
	return &p
}
