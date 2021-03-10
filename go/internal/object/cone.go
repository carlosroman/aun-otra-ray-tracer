package object

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type cone struct {
	obj
	minimum, maximum float64
	closed           bool
}

func (c cone) LocalIntersect(r ray.Ray) (xs Intersections) {
	a := math.Pow(r.Direction().GetX(), 2) -
		math.Pow(r.Direction().GetY(), 2) +
		math.Pow(r.Direction().GetZ(), 2)

	b := (2 * r.Origin().GetX() * r.Direction().GetX()) -
		(2 * r.Origin().GetY() * r.Direction().GetY()) +
		(2 * r.Origin().GetZ() * r.Direction().GetZ())

	cc := math.Pow(r.Origin().GetX(), 2) -
		math.Pow(r.Origin().GetY(), 2) + math.Pow(r.Origin().GetZ(), 2)

	absA := math.Abs(a)

	if absA <= epsilon && math.Abs(b) > epsilon {
		t := -cc / (2 * b)
		xs = append(xs, Intersection{
			T:   t,
			Obj: &c,
		})
	}

	if absA <= epsilon {
		return c.intersectCaps(r, xs)
	}

	disc := math.Pow(b, 2) - 4*a*cc
	if disc < 0 {
		return xs
	}

	sqrtDisc := math.Sqrt(disc)
	a2 := 2 * a
	t0 := (-b - sqrtDisc) / a2
	t1 := (-b + sqrtDisc) / a2

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

func (c cone) intersectCaps(r ray.Ray, xs Intersections) Intersections {
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

func (c cone) LocalNormalAt(worldPoint ray.Vector, hit Intersection) ray.Vector {

	dist := math.Pow(worldPoint.GetX(), 2) + math.Pow(worldPoint.GetZ(), 2)

	if dist < 1 && worldPoint.GetY() >= (c.maximum-epsilon) {
		return ray.NewVec(0, 1, 0)
	}

	if dist < 1 && worldPoint.GetY() <= (c.minimum+epsilon) {
		return ray.NewVec(0, -1, 0)
	}

	y := math.Sqrt(dist)
	if worldPoint.GetY() > 0 {
		y = -y
	}

	return ray.NewVec(worldPoint.GetX(), y, worldPoint.GetZ())
}

func DefaultCone() Object {
	return NewCone(math.Inf(-1), math.Inf(1), false)
}

func NewCone(minimum, maximum float64, closed bool) Object {
	c := cone{
		minimum: minimum,
		maximum: maximum,
		closed:  closed,
	}
	_ = c.SetTransform(ray.DefaultIdentityMatrix())
	c.m = DefaultMaterial()
	return &c
}
