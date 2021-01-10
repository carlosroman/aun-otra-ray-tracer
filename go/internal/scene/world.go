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

func Hit(intersections Intersections) Intersection {
	for i := range intersections {
		if 0 <= intersections[i].T {
			return intersections[i]
		}
	}
	return NoHit
}

var NoHit = Intersection{}
