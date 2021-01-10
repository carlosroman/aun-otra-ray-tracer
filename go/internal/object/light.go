package object

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type PointLight struct {
	Position  ray.Vector
	Intensity RGB
}

func NewPointLight(position ray.Vector, intensity RGB) PointLight {
	return PointLight{
		Position:  position,
		Intensity: intensity,
	}
}

func Lighting(
	material Material,
	light PointLight,
	point, eyev, normalv ray.Vector) (result RGB) {

	// combine the surface color with the light's color/intensity
	// effective color
	ec := material.Color.Multiply(light.Intensity)

	// find the direction to the light source
	lightv := light.Position.Subtract(point).Normalize()

	// compute the ambient contribution
	ambient := ec.MultiplyBy(material.Ambient)

	ldn := ray.Dot(lightv, normalv)

	diffuse := RGB{R: 0, G: 0, B: 0}
	specular := RGB{R: 0, G: 0, B: 0}
	if ldn >= 0 {

		// compute the diffuse contribution
		diffuse = ec.MultiplyBy(material.Diffuse).MultiplyBy(ldn)

		reflectv := lightv.Negate().Reflect(normalv)
		rde := ray.Dot(reflectv, eyev)

		if rde > 0 {
			factor := math.Pow(rde, material.Shininess)
			specular = light.Intensity.MultiplyBy(material.Specular).MultiplyBy(factor)
		}
	}
	return ambient.Plus(diffuse).Plus(specular)
}

type RGB struct {
	R, G, B float64
}

func (r RGB) Multiply(c RGB) (result RGB) {
	return RGB{
		R: r.R * c.R,
		G: r.G * c.G,
		B: r.B * c.B,
	}
}

func (r RGB) MultiplyBy(scalar float64) (result RGB) {
	return RGB{
		R: r.R * scalar,
		G: r.G * scalar,
		B: r.B * scalar,
	}
}

func (r RGB) Plus(c RGB) (result RGB) {
	return RGB{
		R: r.R + c.R,
		G: r.G + c.G,
		B: r.B + c.B,
	}
}
