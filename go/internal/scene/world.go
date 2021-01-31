package scene

import (
	"math"
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
	ColorAt(r ray.Ray, remaining int) object.RGB
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

func (w *world) ColorAt(r ray.Ray, remaining int) (color object.RGB) {
	intersections := Intersect(w, r)
	hit := Hit(intersections)
	if hit == NoHit {
		return color
	}
	return ShadeHit(
		w,
		PrepareComputations(hit, r, intersections...),
		remaining)
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
	Obj object.Object // Intersected object
}

type Intersections []Intersection

func Intersect(w World, r ray.Ray) (intersections Intersections) {
	for i := range w.Objects() {
		hits := object.Intersect(w.Objects()[i], r)
		tmpHits := make(Intersections, len(hits))
		for hit := range hits {
			tmpHits[hit] = Intersection{
				T:   hits[hit],
				Obj: w.Objects()[i],
			}
		}
		intersections = append(intersections, tmpHits...)
	}
	sort.SliceStable(intersections, func(i, j int) bool {
		return intersections[i].T < intersections[j].T
	})
	return intersections
}

func ShadeHit(w World, comps Computation, remaining int) object.RGB {
	reflected := ReflectedColor(w, comps, remaining)
	refracted := RefractedColor(w, comps, remaining)

	surface := object.Lighting(
		comps.obj.Material(),
		comps.obj,
		w.Light(),
		comps.overPoint, comps.eyev, comps.normalv,
		w.IsShadowed(comps.overPoint))

	if comps.obj.Material().Reflective > 0 &&
		comps.obj.Material().Transparency > 0 {
		reflectance := Schlick(comps)
		return surface.
			Add(reflected.MultiplyBy(reflectance)).
			Add(refracted.MultiplyBy(1 - reflectance))
	}
	return surface.
		Add(reflected).
		Add(refracted)
}

func ReflectedColor(w World, comps Computation, remaining int) object.RGB {
	if remaining <= 0 {
		return object.Black
	}
	if comps.obj.Material().Reflective == 0 {
		return object.Black
	}
	reflectRay := ray.NewRayAt(comps.overPoint, comps.reflectv)
	remaining--
	color := w.ColorAt(reflectRay, remaining)
	return color.MultiplyBy(comps.obj.Material().Reflective)
}

func RefractedColor(w World, comps Computation, remaining int) object.RGB {
	if remaining == 0 {
		return object.Black
	}
	if comps.obj.Material().Transparency == 0 {
		return object.Black
	}

	nRatio := comps.n1 / comps.n2
	cosI := ray.Dot(comps.eyev, comps.normalv)
	sin2t := math.Pow(nRatio, 2) * (1 - math.Pow(cosI, 2))
	if sin2t > 1 {
		return object.Black
	}

	cosT := math.Sqrt(1.0 - sin2t)
	direction := comps.normalv.
		Multiply(nRatio*cosI - cosT).
		Subtract(comps.eyev.Multiply(nRatio))

	refractRay := ray.NewRayAt(comps.underPoint, direction)
	remaining--
	return w.ColorAt(refractRay, remaining).
		MultiplyBy(comps.obj.Material().Transparency)
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
