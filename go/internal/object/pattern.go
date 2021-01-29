package object

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type StripePattern struct {
	A, B      RGB
	Transform ray.Matrix
}

func EmptyStripePattern() StripePattern {
	return StripePattern{}
}

func NewStripePattern(a, b RGB) StripePattern {
	return StripePattern{
		A:         a,
		B:         b,
		Transform: ray.DefaultIdentityMatrix(),
	}
}

func (s StripePattern) At(point ray.Vector) RGB {
	if math.Remainder(math.Floor(point.GetX()), 2) == 0 {
		return s.A
	}
	return s.B
}

func (s StripePattern) AtObj(obj Object, worldPoint ray.Vector) RGB {
	objInv, err := obj.Transform().Inverse()
	if err != nil {
		return RGB{}
	}
	patternInv, err := s.Transform.Inverse()
	if err != nil {
		return RGB{}
	}
	objPoint := objInv.MultiplyByVector(worldPoint)
	patternPoint := patternInv.MultiplyByVector(objPoint)
	return s.At(patternPoint)
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
