package object

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type Pattern struct {
	Transform ray.Matrix
	At        func(point ray.Vector) RGB
}

func (p Pattern) AtObj(obj Object, worldPoint ray.Vector) RGB {
	objInv, err := obj.Transform().Inverse()
	if err != nil {
		return RGB{}
	}
	patternInv, err := p.Transform.Inverse()
	if err != nil {
		return RGB{}
	}
	objPoint := objInv.MultiplyByVector(worldPoint)
	patternPoint := patternInv.MultiplyByVector(objPoint)
	return p.At(patternPoint)
}

type StripePattern struct {
	Pattern
	A, B RGB
}

func EmptyStripePattern() StripePattern {
	return StripePattern{}
}

func NewStripePattern(a, b RGB) StripePattern {
	p := StripePattern{
		A: a,
		B: b,
	}
	p.Transform = ray.DefaultIdentityMatrix()
	p.At = func(point ray.Vector) RGB {
		if math.Remainder(math.Floor(point.GetX()), 2) == 0 {
			return p.A
		}
		return p.B
	}
	return p
}

func (s StripePattern) IsEmpty() bool {
	if s.A == emptyStripePattern.A &&
		s.B == emptyStripePattern.B &&
		len(s.Transform) < 1 {
		return true
	}
	return false
}

var (
	emptyStripePattern = StripePattern{}
)

type TestPattern struct {
	Pattern
}

func NewTestPattern() (p TestPattern) {
	p = TestPattern{}
	p.Transform = ray.DefaultIdentityMatrix()
	p.At = func(point ray.Vector) RGB {
		return RGB{
			R: point.GetX(),
			G: point.GetY(),
			B: point.GetZ(),
		}
	}
	return p
}
