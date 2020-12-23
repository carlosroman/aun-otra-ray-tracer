package scene

import (
	"sort"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type World interface {
	Objects() []object.Object
	AddObject(obj object.Object)
}

type world struct {
	objs []object.Object
}

func (w world) Objects() []object.Object {
	return w.objs
}

func (w *world) AddObject(obj object.Object) {
	w.objs = append(w.objs, obj)
}

func NewWorld() World {
	return &world{
		objs: nil,
	}
}

type Intersection struct {
	T   float64
	Obj object.Object
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
