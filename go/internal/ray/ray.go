package ray

import (
	"image/color"
)

type ray struct {
	o     Vector
	d     Vector
	color color.RGBA
}

func (r ray) Origin() Vector {
	return r.o
}

func (r ray) Direction() Vector {
	return r.d
}

func (r ray) PointAt(parameter float64) Vector {
	return r.o.Add(r.d.Multiply(parameter))
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
	Origin() Vector
	Direction() Vector
	PointAt(parameter float64) Vector
	SetR(redVal uint8)
	SetG(greenVal uint8)
	SetB(blueVal uint8)
}

func NewRayAt(point, vector Vector) Ray {
	return &ray{
		o: point,
		d: vector,
	}
}

func NewRay(x, y, z float64, r, g, b uint8) Ray {
	return &ray{
		o: NewVec(x, y, z),
		color: color.RGBA{
			R: r,
			G: g,
			B: b,
			A: 0xff,
		},
	}
}
