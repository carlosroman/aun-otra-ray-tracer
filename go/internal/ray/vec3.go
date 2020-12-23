package ray

import (
	"math"
)

type vec3 struct {
	x, y, z, w float64
}

var zero = NewVec(0, 0, 0)

func newTuple(x, y, z, w float64) Vector {
	return vec3{
		x: x,
		y: y,
		z: z,
		w: w,
	}
}

func NewVec(x, y, z float64) Vector {
	return newTuple(x, y, z, 0)
}

func NewPoint(x, y, z float64) Vector {
	return newTuple(x, y, z, 1)
}

func (v vec3) GetX() float64 {
	return v.x
}

func (v vec3) GetY() float64 {
	return v.y
}

func (v vec3) GetZ() float64 {
	return v.z
}

func (v vec3) GetW() float64 {
	return v.w
}

func (v vec3) Add(vec Vector) Vector {
	return NewVec(
		v.x+vec.GetX(),
		v.y+vec.GetY(),
		v.z+vec.GetZ())
}

func (v vec3) Subtract(vec Vector) Vector {
	return NewVec(
		v.x-vec.GetX(),
		v.y-vec.GetY(),
		v.z-vec.GetZ())
}

func (v vec3) Multiply(by float64) Vector {
	return NewVec(
		v.x*by,
		v.y*by,
		v.z*by)
}

func (v vec3) Divide(by float64) Vector {
	return v.Multiply(1 / by)
}

func (v vec3) Negate() Vector {
	return zero.Subtract(v)
}

func (v vec3) Magnitude() float64 {
	return math.Sqrt(
		(v.x * v.x) + (v.y * v.y) + (v.z * v.z))
}

func (v vec3) Normalize() Vector {
	magnitude := v.Magnitude()
	return v.Divide(magnitude)
}
