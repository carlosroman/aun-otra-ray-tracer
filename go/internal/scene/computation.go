package scene

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

const (
	epsilon = 0.00000001
)

type Computation struct {
	t          float64 //Intersect
	obj        object.Object
	point      ray.Vector
	overPoint  ray.Vector
	eyev       ray.Vector
	normalv    ray.Vector
	reflectv   ray.Vector
	underPoint ray.Vector
	inside     bool
	n1         float64
	n2         float64
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

func (c Computation) Reflectv() ray.Vector {
	return c.reflectv
}

func (c Computation) Inside() bool {
	return c.inside
}

func (c Computation) OverPoint() ray.Vector {
	return c.overPoint
}

func (c Computation) UnderPoint() ray.Vector {
	return c.underPoint
}

func (c Computation) N1() float64 {
	return c.n1
}

func (c Computation) N2() float64 {
	return c.n2
}

func contains(objs []object.Object, obj object.Object) (found bool, idx int) {
	for i := range objs {
		if objs[i] == obj {
			return true, i
		}
	}
	return false, 0
}

func PrepareComputations(i Intersection, r ray.Ray, xs ...Intersection) (comps Computation) {

	comps.t = i.T
	comps.obj = i.Obj
	comps.point = r.PointAt(comps.t)
	comps.eyev = r.Direction().Negate()
	comps.normalv = object.NormalAt(i.Obj, comps.point)
	comps.reflectv = r.Direction().Reflect(comps.normalv)

	if ray.Dot(comps.normalv, comps.eyev) < 0 {
		comps.inside = true
		comps.normalv = comps.normalv.Negate()
	}

	normalvMultiplyByEpsilon := comps.normalv.Multiply(epsilon)
	comps.overPoint = comps.point.Add(normalvMultiplyByEpsilon)
	comps.underPoint = comps.point.Subtract(normalvMultiplyByEpsilon)

	var containers []object.Object

	comps.n1 = 1.0
	comps.n2 = 1.0
	for idx := range xs {
		if i == xs[idx] {
			if len(containers) > 0 {
				comps.n1 = containers[len(containers)-1].Material().RefractiveIndex
			}
		}

		if found, at := contains(containers, xs[idx].Obj); found {
			containers = append(containers[:at], containers[at+1:]...)
		} else {
			containers = append(containers, xs[idx].Obj)
		}

		if i == xs[idx] {
			if len(containers) > 0 {
				comps.n2 = containers[len(containers)-1].Material().RefractiveIndex
			}
			break
		}
	}
	return comps
}

func Schlick(comps Computation) float64 {
	cos := ray.Dot(comps.eyev, comps.normalv)
	if comps.n1 > comps.n2 {
		n := comps.n1 / comps.n2
		sin2t := math.Pow(n, 2) * (1.0 - math.Pow(cos, 2))
		if sin2t > 1.0 {
			return 1.0
		}

		cos = math.Sqrt(1.0 - sin2t)
	}
	r0 := math.Pow((comps.n1-comps.n2)/(comps.n1+comps.n2), 2)
	return r0 + (1-r0)*math.Pow(1-cos, 5)
}
