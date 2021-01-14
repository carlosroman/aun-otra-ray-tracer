package scene

import (
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

const (
	epsilon = 0.00000001
)

type Computation struct {
	t         float64 //Intersect
	obj       object.Object
	point     ray.Vector
	overPoint ray.Vector
	eyev      ray.Vector
	normalv   ray.Vector
	inside    bool
}

func (c Computation) Intersect() float64 {
	return c.t
}

func (c Computation) Object() object.Object {
	return c.obj
}

func (c Computation) Point() ray.Vector {
	return c.point
}

func (c Computation) Eyev() ray.Vector {
	return c.eyev
}

func (c Computation) Normalv() ray.Vector {
	return c.normalv
}

func (c Computation) Inside() bool {
	return c.inside
}

func (c Computation) OverPoint() ray.Vector {
	return c.overPoint
}

func PrepareComputations(i Intersection, r ray.Ray) (comps Computation) {
	comps.t = i.T
	comps.obj = i.Obj
	comps.point = r.PointAt(comps.t)
	comps.eyev = r.Direction().Negate()
	comps.normalv = i.Obj.NormalAt(comps.point)

	if ray.Dot(comps.normalv, comps.eyev) < 0 {
		comps.inside = true
		comps.normalv = comps.normalv.Negate()
	}
	comps.overPoint = comps.point.Add(comps.normalv.Multiply(epsilon))
	return comps
}
