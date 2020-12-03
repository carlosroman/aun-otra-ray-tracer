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

type HitPoint struct {
	T   float64
	Obj object.Object
}

func Intersect(w World, r ray.Ray) (intersects []HitPoint) {
	for i := range w.Objects() {
		hits := w.Objects()[i].Intersect(r)
		for hit := range hits {
			intersects = append(intersects, HitPoint{
				T:   hits[hit],
				Obj: w.Objects()[i],
			})
		}
	}
	sort.SliceStable(intersects, func(i, j int) bool {
		return intersects[i].T < intersects[j].T
	})
	return intersects
}

func Hit(intersects []HitPoint) HitPoint {
	for i := range intersects {
		if 0 <= intersects[i].T {
			return intersects[i]
		}
	}
	return NoHit
}

var NoHit = HitPoint{}
