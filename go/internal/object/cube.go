package object

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/ray"
)

type cube struct {
	obj
}

func (c cube) LocalIntersect(r ray.Ray) Intersections {
	xtmin, xtmax := checkAxis(r.Origin().GetX(), r.Direction().GetX())
	ytmin, ytmax := checkAxis(r.Origin().GetY(), r.Direction().GetY())
	ztmin, ztmax := checkAxis(r.Origin().GetZ(), r.Direction().GetZ())

	tmax := math.Min(math.Min(xtmax, ytmax), ztmax)
	tmin := math.Max(math.Max(xtmin, ytmin), ztmin)

	if tmin > tmax {
		return nil
	}
	return []Intersection{{T: tmin, Obj: &c}, {T: tmax, Obj: &c}}
}

func checkAxis(origin, direction float64) (tmin, tmax float64) {
	tminNumerator := -1 - origin
	tmaxNumerator := 1 - origin

	tmin = tminNumerator / direction
	tmax = tmaxNumerator / direction

	if tmin > tmax {
		return tmax, tmin
	}
	return tmin, tmax
}

func (c cube) LocalNormalAt(point ray.Vector, _ Intersection) ray.Vector {
	absX := math.Abs(point.GetX())
	absY := math.Abs(point.GetY())
	maxc := math.Max(math.Max(absX, absY),
		math.Abs(point.GetZ()))

	if maxc == absX {
		return ray.NewVec(point.GetX(), 0, 0)
	}

	if maxc == absY {
		return ray.NewVec(0, point.GetY(), 0)
	}

	return ray.NewVec(0, 0, point.GetZ())
}

func NewCube() Object {
	c := cube{}
	_ = c.SetTransform(ray.DefaultIdentityMatrix())
	c.m = DefaultMaterial()
	return &c
}
