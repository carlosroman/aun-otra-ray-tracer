package object

const (
	defaultRGB                     = float64(1)
	defaultMaterialAmbient         = 0.1
	defaultMaterialDiffuse         = 0.9
	defaultMaterialSpecular        = 0.9
	defaultMaterialShininess       = 200.0
	defaultMaterialReflective      = 0.0
	defaultMaterialTransparency    = 0.0
	defaultMaterialRefractiveIndex = 1.0
)

type Material struct {
	Color                                                                            RGB
	Pattern                                                                          Pattern
	Ambient, Diffuse, Specular, Shininess, Reflective, Transparency, RefractiveIndex float64
}

func NewMaterial(color RGB,
	ambient, diffuse, specular, shininess, reflective, transparency, refractiveIndex float64) Material {
	return Material{
		Pattern:         Pattern{},
		Color:           color,
		Ambient:         ambient,
		Diffuse:         diffuse,
		Specular:        specular,
		Shininess:       shininess,
		Reflective:      reflective,
		Transparency:    transparency,
		RefractiveIndex: refractiveIndex,
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
		defaultMaterialTransparency,
		defaultMaterialRefractiveIndex,
	)
}
