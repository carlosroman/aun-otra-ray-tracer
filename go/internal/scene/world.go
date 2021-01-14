package scene

import (
	"sort"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type World interface {
	Objects() []object.Object
	AddObject(obj object.Object)
	AddObjects(objs ...object.Object)
	Light() object.PointLight
	AddLight(light object.PointLight)
	ColorAt(r ray.Ray) object.RGB
	IsShadowed(point ray.Vector) bool
}

type world struct {
	objs  []object.Object
	light object.PointLight
}

func (w world) Objects() []object.Object {
	return w.objs
}

func (w *world) AddObject(obj object.Object) {
	w.objs = append(w.objs, obj)
}

func (w *world) AddObjects(objs ...object.Object) {
	for _, obj := range objs {
		w.objs = append(w.objs, obj)
	}
}

func (w world) Light() object.PointLight {
	return w.light
}

func (w *world) AddLight(light object.PointLight) {
	w.light = light
}

func (w *world) ColorAt(r ray.Ray) (color object.RGB) {
	hit := Hit(Intersect(w, r))
	if hit == NoHit {
		return color
	}
	return ShadeHit(
		w,
		PrepareComputations(hit, r))
}

func (w *world) IsShadowed(point ray.Vector) bool {
	v := w.light.Position.Subtract(point)
	distance := v.Magnitude()
	direction := v.Normalize()
	r := ray.NewRayAt(point, direction)
	intersections := Intersect(w, r)
	h := Hit(intersections)
	if h != NoHit && h.T < distance {
		return true
	}
	return false
}

func NewWorld() World {
	return &world{
		objs: nil,
	}
}

type Intersection struct {
	T   float64       // Intersection value
	Obj object.Object // Interesected object
}

type Intersections []Intersection

func Intersect(w World, r ray.Ray) (intersections Intersections) {
	for i := range w.Objects() {
		hits := w.Objects()[i].Intersect(r)
		for hit := range hits {
			intersections = append(intersections, Intersection{
				T:   hits[hit],
				Obj: w.Objects()[i],
			})
		}
	}
	sort.SliceStable(intersections, func(i, j int) bool {
		return intersections[i].T < intersections[j].T
	})
	return intersections
}

func ShadeHit(w World, comps Computation) object.RGB {
	return object.Lighting(
		comps.obj.Material(),
		w.Light(),
		comps.overPoint, comps.eyev, comps.normalv,
		w.IsShadowed(comps.overPoint))
}

func Hit(intersections Intersections) Intersection {
	for i := range intersections {
		if 0 <= intersections[i].T {
			return intersections[i]
		}
	}
	return NoHit
}

var NoHit = Intersection{}
