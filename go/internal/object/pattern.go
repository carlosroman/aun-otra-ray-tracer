package object

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type StripePattern struct {
	A, B RGB
}

func NewStripePattern(a, b RGB) StripePattern {
	return StripePattern{
		A: a,
		B: b,
	}
}

func (s StripePattern) At(point ray.Vector) RGB {
	if math.Remainder(math.Floor(point.GetX()), 2) == 0 {
		return s.A
	}
	return s.B
}
