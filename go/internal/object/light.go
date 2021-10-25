package object

import (
	"math"

	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/ray"
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
	obj Object,
	light PointLight,
	position, eyev, normalv ray.Vector,
	inShadows bool) RGB {

	color := material.Color
	if material.Pattern.IsNotEmpty {
		color = material.Pattern.AtObj(obj, position)
	}
	// combine the surface color with the light's color/intensity
	// effective color
	ec := color.Multiply(light.Intensity)

	// find the direction to the light source
	lightv := light.Position.Subtract(position).Normalize()

	// compute the ambient contribution
	ambient := ec.MultiplyBy(material.Ambient)

	if inShadows {
		return ambient
	}

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
