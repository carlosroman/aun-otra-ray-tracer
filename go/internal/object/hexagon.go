package object

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

func hexagonCorner() (corner Object) {
	corner = DefaultSphere()
	_ = corner.SetTransform(
		ray.Translation(0, 0, -1).
			Multiply(ray.Scaling(0.25, 0.25, 0.25)))
	return corner
}

func hexagonEdge() (edge Object) {
	edge = NewCylinder(0, 1, false)
	_ = edge.SetTransform(
		ray.Translation(0, 0, -1).
			Multiply(ray.Rotation(ray.Y, -math.Pi/6)).
			Multiply(ray.Rotation(ray.Z, -math.Pi/2)).
			Multiply(ray.Scaling(0.25, 1, 0.25)))
	return edge
}

func hexagonSide() (side Object) {
	s := NewGroup()
	s.AddChild(hexagonCorner(), hexagonEdge())
	return &s
}

func NewHexagon() (hex Object) {
	h := NewGroup()

	for i := float64(0); i < 6; i++ {
		s := hexagonSide()
		_ = s.SetTransform(
			ray.Rotation(ray.Y, i*(math.Pi/3)))
		h.AddChild(s)
	}
	return &h
}
