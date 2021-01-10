package object

import "image/color"

const (
	defaultRGB               = 1
	defaultMaterialAmbient   = 0.1
	defaultMaterialDiffuse   = 0.9
	defaultMaterialSpecular  = 0.9
	defaultMaterialShininess = 200.0
)

type Material struct {
	Color                                 color.RGBA
	Ambient, Diffuse, Specular, Shininess float64
}

func DefaultMaterial() Material {
	return Material{
		Color:     NewColor(defaultRGB, defaultRGB, defaultRGB),
		Ambient:   defaultMaterialAmbient,
		Diffuse:   defaultMaterialDiffuse,
		Specular:  defaultMaterialSpecular,
		Shininess: defaultMaterialShininess,
	}
}

func NewColor(red, green, blue uint8) color.RGBA {
	return color.RGBA{
		R: red,
		G: green,
		B: blue,
	}
}
