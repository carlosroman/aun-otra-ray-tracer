package object

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type cube struct {
	t    ray.Matrix
	tInv ray.Matrix
	m    Material
}

func (c cube) Transform() ray.Matrix {
	return c.t
}

func (c *cube) SetTransform(t ray.Matrix) error {
	c.t = t
	inverse, err := t.Inverse()
	c.tInv = inverse
	return err
}

func (c cube) TransformInverse() ray.Matrix {
	return c.tInv
}

func (c cube) Material() Material {
	return c.m
}

func (c *cube) SetMaterial(m Material) {
	c.m = m
}

func (c cube) LocalIntersect(ray ray.Ray) []float64 {
	xtmin, xtmax := checkAxis(ray.Origin().GetX(), ray.Direction().GetX())
	ytmin, ytmax := checkAxis(ray.Origin().GetY(), ray.Direction().GetY())
	ztmin, ztmax := checkAxis(ray.Origin().GetZ(), ray.Direction().GetZ())

	tmax := math.Min(math.Min(xtmax, ytmax), ztmax)
	tmin := math.Max(math.Max(xtmin, ytmin), ztmin)

	if tmin > tmax {
		return nil
	}
	return []float64{tmin, tmax}
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

func (c cube) LocalNormalAt(point ray.Vector) ray.Vector {
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
	return &cube{
		t:    ray.DefaultIdentityMatrix(),
		tInv: ray.DefaultIdentityMatrixInverse(),
		m:    DefaultMaterial(),
	}
}
