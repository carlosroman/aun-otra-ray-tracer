package object

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type Pattern struct {
	Transform        ray.Matrix
	TransformInverse ray.Matrix
	At               func(point ray.Vector) RGB
	IsNotEmpty       bool
}

func (p Pattern) AtObj(obj Object, worldPoint ray.Vector) RGB {
	objPoint := obj.WorldToObject(worldPoint)
	patternPoint := p.TransformInverse.MultiplyByVector(objPoint)
	return p.At(patternPoint)
}

func (p *Pattern) SetTransform(transform ray.Matrix) error {
	p.Transform = transform
	inverse, err := transform.Inverse()
	p.TransformInverse = inverse
	return err
}

type StripePattern struct {
	Pattern
	A, B RGB
}

func EmptyPattern() Pattern {
	return Pattern{
		Transform:        ray.DefaultIdentityMatrix(),
		TransformInverse: ray.DefaultIdentityMatrixInverse(),
		At: func(_ ray.Vector) RGB {
			return Black
		},
	}
}

func NewStripePattern(a, b RGB) (p Pattern) {
	p = NewTestPattern()
	p.At = func(point ray.Vector) RGB {
		if math.Remainder(math.Floor(point.GetX()), 2) == 0 {
			return a
		}
		return b
	}
	return p
}

func NewTestPattern() (p Pattern) {
	p = Pattern{
		IsNotEmpty: true,
	}
	_ = p.SetTransform(ray.DefaultIdentityMatrix())
	p.At = func(point ray.Vector) RGB {
		return RGB{
			R: point.GetX(),
			G: point.GetY(),
			B: point.GetZ(),
		}
	}
	return p
}

func NewGradientPattern(a, b RGB) (p Pattern) {
	p = NewTestPattern()
	d := b.Subtract(a)
	p.At = func(point ray.Vector) RGB {
		fraction := point.GetX() - math.Floor(point.GetX())
		return a.Add(d.MultiplyBy(fraction))
	}
	return p
}

func NewRingPattern(a, b RGB) (p Pattern) {
	p = NewTestPattern()
	p.At = func(point ray.Vector) RGB {
		valX := math.Pow(point.GetX(), 2)
		valZ := math.Pow(point.GetZ(), 2)
		val := math.Floor(math.Sqrt(valX + valZ))
		if math.Mod(val, 2) == 0 {
			return a
		}
		return b
	}
	return p
}

func NewCheckerPattern(a, b RGB) (p Pattern) {
	p = NewTestPattern()
	p.At = func(point ray.Vector) RGB {
		val := math.Floor(point.GetX()) +
			math.Floor(point.GetY()) +
			math.Floor(point.GetZ())
		if math.Mod(val, 2) == 0 {
			return a
		}
		return b
	}
	return p
}
