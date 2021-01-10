package scene

import (
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

func DefaultWorld() World {
	light := DefaultWorldLight()

	s1 := object.NewSphere(ray.ZeroPoint, 1)
	m1 := object.DefaultMaterial()
	m1.Color = object.NewColor(0.8, 1.0, 0.6)
	m1.Diffuse = 0.7
	m1.Specular = 0.2
	s1.SetMaterial(m1)

	s2 := object.NewSphere(ray.ZeroPoint, 1)
	s2.SetTransform(ray.Scaling(0.5, 0.5, 0.5))

	w := NewWorld()
	w.AddObjects(s1, s2)
	w.AddLight(light)
	return w
}

func DefaultWorldLight() object.PointLight {
	return object.NewPointLight(ray.NewPoint(-10, 10, -10), object.NewColor(1, 1, 1))
}

func TwoSphereWorld() World {
	s1 := object.NewSphere(ray.NewPoint(0, 0, -1), 0.5)
	s2 := object.NewSphere(ray.NewPoint(0, -100.5, -1), 100)
	w := NewWorld()
	w.AddObjects(s1, s2)
	return w
}
