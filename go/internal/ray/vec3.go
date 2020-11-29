package ray

import "math"

type vec3 struct {
	x, y, z float64
}

var zero = NewVec(0, 0, 0)

func NewVec(x, y, z float64) vec3 {
	return vec3{
		x: x,
		y: y,
		z: z,
	}
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

func (v vec3) Add(vec vec3) vec3 {
	return NewVec(
		v.x+vec.x,
		v.y+vec.y,
		v.z+vec.z)
}

func (v vec3) Subtract(vec vec3) vec3 {
	return NewVec(
		v.x-vec.x,
		v.y-vec.y,
		v.z-vec.z)
}

func (v vec3) Multiply(by float64) vec3 {
	return NewVec(
		v.x*by,
		v.y*by,
		v.z*by)
}

func (v vec3) Divide(by float64) vec3 {
	return v.Multiply(1 / by)
}

func (v vec3) Negate() vec3 {
	return zero.Subtract(v)
}

func (v vec3) Magnitude() float64 {
	return math.Sqrt(
		(v.x * v.x) + (v.y * v.y) + (v.z * v.z))
}

func (v vec3) Normalize() vec3 {
	magnitude := v.Magnitude()
	return v.Divide(magnitude)
}
