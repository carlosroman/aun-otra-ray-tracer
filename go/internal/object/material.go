package object

const (
	defaultRGB               = float64(1)
	defaultMaterialAmbient   = 0.1
	defaultMaterialDiffuse   = 0.9
	defaultMaterialSpecular  = 0.9
	defaultMaterialShininess = 200.0
)

type Material struct {
	Color                                 RGB
	Pattern                               StripePattern
	Ambient, Diffuse, Specular, Shininess float64
}

func NewMaterial(color RGB,
	ambient, diffuse, specular, shininess float64) Material {
	return Material{
		Pattern:   EmptyStripePattern(),
		Color:     color,
		Ambient:   ambient,
		Diffuse:   diffuse,
		Specular:  specular,
		Shininess: shininess,
	}
}

func DefaultMaterial() Material {
	return NewMaterial(
		NewColor(defaultRGB, defaultRGB, defaultRGB),
		defaultMaterialAmbient,
		defaultMaterialDiffuse,
		defaultMaterialSpecular,
		defaultMaterialShininess,
	)
}
