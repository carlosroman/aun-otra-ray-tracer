package object

const (
	defaultRGB                = float64(1)
	defaultMaterialAmbient    = 0.1
	defaultMaterialDiffuse    = 0.9
	defaultMaterialSpecular   = 0.9
	defaultMaterialShininess  = 200.0
	defaultMaterialReflective = 0.0
)

type Material struct {
	Color                                             RGB
	Pattern                                           Pattern
	Ambient, Diffuse, Specular, Shininess, Reflective float64
}

func NewMaterial(color RGB,
	ambient, diffuse, specular, shininess, reflective float64) Material {
	return Material{
		Pattern:    Pattern{},
		Color:      color,
		Ambient:    ambient,
		Diffuse:    diffuse,
		Specular:   specular,
		Shininess:  shininess,
		Reflective: reflective,
	}
}

func DefaultMaterial() Material {
	return NewMaterial(
		NewColor(defaultRGB, defaultRGB, defaultRGB),
		defaultMaterialAmbient,
		defaultMaterialDiffuse,
		defaultMaterialSpecular,
		defaultMaterialShininess,
		defaultMaterialReflective,
	)
}
