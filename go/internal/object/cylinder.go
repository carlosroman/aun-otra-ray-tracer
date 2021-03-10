package object

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type cylinder struct {
	obj
	minimum, maximum float64
	closed           bool
}

func checkCap(r ray.Ray, t float64) bool {
	x := r.Origin().GetX() + t*r.Direction().GetX()
	z := r.Origin().GetZ() + t*r.Direction().GetZ()
	return (math.Pow(x, 2) + math.Pow(z, 2)) <= 1
}

func (c cylinder) LocalIntersect(r ray.Ray) (xs Intersections) {
	a := 2 * (math.Pow(r.Direction().GetX(), 2) + math.Pow(r.Direction().GetZ(), 2))
	if a < epsilon {
		return c.intersectCaps(r, xs)
	}

	b := 2 * (r.Origin().GetX()*r.Direction().GetX() +
		r.Origin().GetZ()*r.Direction().GetZ())

	cc := math.Pow(r.Origin().GetX(), 2) +
		math.Pow(r.Origin().GetZ(), 2) - 1

	disc := math.Pow(b, 2) - 2*a*cc
	if disc < 0 {
		return xs
	}

	sqrtDisc := math.Sqrt(disc) / a
	x := -b / a
	t0 := x - sqrtDisc
	t1 := x + sqrtDisc

	if t0 > t1 {
		t := t0
		t0 = t1
		t1 = t
	}

	y0 := r.Origin().GetY() + t0*r.Direction().GetY()
	if (c.minimum < y0) && (y0 < c.maximum) {
		xs = append(xs, Intersection{
			T:   t0,
			Obj: &c,
		})
	}
	y1 := r.Origin().GetY() + t1*r.Direction().GetY()

	if (c.minimum < y1) && (y1 < c.maximum) {
		xs = append(xs, Intersection{
			T:   t1,
			Obj: &c,
		})
	}

	return c.intersectCaps(r, xs)
}

func (c cylinder) intersectCaps(r ray.Ray, xs Intersections) Intersections {
	if !c.closed || math.Abs(r.Direction().GetY()) <= epsilon {
		return xs
	}

	t0 := (c.minimum - r.Origin().GetY()) / r.Direction().GetY()

	if checkCap(r, t0) {
		xs = append(xs, Intersection{
			T:   t0,
			Obj: &c,
		})
	}

	t1 := (c.maximum - r.Origin().GetY()) / r.Direction().GetY()

	if checkCap(r, t1) {
		xs = append(xs, Intersection{
			T:   t1,
			Obj: &c,
		})
	}

	return xs
}

func (c cylinder) LocalNormalAt(worldPoint ray.Vector, hit Intersection) ray.Vector {

	dist := math.Pow(worldPoint.GetX(), 2) + math.Pow(worldPoint.GetZ(), 2)

	if dist < 1 && worldPoint.GetY() >= (c.maximum-epsilon) {
		return ray.NewVec(0, 1, 0)
	}

	if dist < 1 && worldPoint.GetY() <= (c.minimum+epsilon) {
		return ray.NewVec(0, -1, 0)
	}

	return ray.NewVec(worldPoint.GetX(), 0, worldPoint.GetZ())
}

func DefaultCylinder() Object {
	return NewCylinder(math.Inf(-1), math.Inf(1), false)
}

func NewCylinder(minimum, maximum float64, closed bool) Object {
	c := cylinder{
		minimum: minimum,
		maximum: maximum,
		closed:  closed,
	}
	_ = c.SetTransform(ray.DefaultIdentityMatrix())
	c.m = DefaultMaterial()
	return &c
}
