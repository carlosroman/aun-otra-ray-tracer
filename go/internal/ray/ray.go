package ray

import (
	"image/color"
)

type ray struct {
	v     Vector
	color color.RGBA
}

func (r ray) GetX() float64 {
	return r.v.GetX()
}

func (r ray) GetY() float64 {
	return r.v.GetY()
}

func (r ray) GetZ() float64 {
	return r.v.GetZ()
}

func (r *ray) SetR(redVal uint8) {
	r.color.R = redVal
}

func (r *ray) SetG(greenVal uint8) {
	r.color.G = greenVal
}
func (r *ray) SetB(blueVal uint8) {
	r.color.B = blueVal
}

func (r ray) RGBA() (red, green, blue, alpha uint32) {
	return r.color.RGBA()
}

type Ray interface {
	color.Color
	Vector
	SetR(redVal uint8)
	SetG(greenVal uint8)
	SetB(blueVal uint8)
	GetX() float64
}

func NewRay(x, y, z float64, r, g, b uint8) Ray {
	return &ray{
		v: NewVec(x, y, z),
		color: color.RGBA{
			R: r,
			G: g,
			B: b,
			A: 0xff,
		},
	}
}
