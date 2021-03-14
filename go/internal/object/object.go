package object

import (
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

const (
	epsilon = 1e-5
)

type Object interface {
	LocalIntersect(ray ray.Ray) Intersections
	LocalNormalAt(worldPoint ray.Vector, hit Intersection) ray.Vector

	Transform() ray.Matrix
	TransformInverse() ray.Matrix
	SetTransform(t ray.Matrix) error

	Material() (m Material)
	SetMaterial(m Material)

	Parent() Object
	SetParent(obj Object)

	WorldToObject(worldPoint ray.Vector) (point ray.Vector)
	NormalToWorld(normalVector ray.Vector) (vector ray.Vector)
}

func NormalAt(xs Intersection, worldPoint ray.Vector) ray.Vector {
	localPoint := xs.Obj.WorldToObject(worldPoint)
	localNormal := xs.Obj.LocalNormalAt(localPoint, xs)
	return xs.Obj.NormalToWorld(localNormal)
}

func Intersect(obj Object, rr ray.Ray) Intersections {
	return obj.
		LocalIntersect(
			rr.Transform(obj.TransformInverse()))
}

type BasicObject struct {
	Transform ray.Matrix
	Material  Material
}

type obj struct {
	t    ray.Matrix
	tInv ray.Matrix
	m    Material
	p    Object
}

func (o obj) LocalNormalAt(worldPoint ray.Vector, _ Intersection) ray.Vector {
	return worldPoint
}

func (o obj) LocalIntersect(r ray.Ray) Intersections {
	return nil
}

func (o obj) Parent() Object {
	return o.p
}

func (o *obj) SetParent(obj Object) {
	o.p = obj
}

func NewTestShape(opts ...Option) Object {
	s := obj{}
	_ = s.SetTransform(ray.DefaultIdentityMatrix())
	m := DefaultMaterial()
	m.Ambient = 1
	s.SetMaterial(m)
	for i := range opts {
		opts[i].Apply(&s)
	}
	return &s
}

func (o obj) Transform() ray.Matrix {
	return o.t
}

func (o *obj) SetTransform(t ray.Matrix) error {
	o.t = t
	inverse, err := t.Inverse()
	o.tInv = inverse
	return err
}

func (o obj) TransformInverse() ray.Matrix {
	return o.tInv
}

func (o obj) Material() Material {
	if o.p != nil {
		return o.p.Material()
	}
	return o.m
}

func (o *obj) SetMaterial(m Material) {
	o.m = m
}

func (o obj) WorldToObject(worldPoint ray.Vector) (point ray.Vector) {
	if o.p != nil {
		worldPoint = o.p.WorldToObject(worldPoint)
	}

	return o.
		tInv.
		MultiplyByVector(worldPoint)
}

func (o obj) NormalToWorld(normaVector ray.Vector) (resultVector ray.Vector) {
	resultVector = o.
		tInv.
		Transpose().
		MultiplyByVector(normaVector).
		SetW(0).
		Normalize()

	if o.p != nil {
		resultVector = o.p.
			NormalToWorld(resultVector)
	}
	return resultVector
}
