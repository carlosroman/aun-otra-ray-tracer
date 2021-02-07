package object

import (
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

const (
	epsilon = 1e-5
)

type Object interface {
	LocalIntersect(ray ray.Ray) []float64
	LocalNormalAt(worldPoint ray.Vector) ray.Vector

	Transform() ray.Matrix
	TransformInverse() ray.Matrix
	SetTransform(t ray.Matrix) error

	Material() (m Material)
	SetMaterial(m Material)
}

type BasicObject struct {
	Transform ray.Matrix
	Material  Material
}

type obj struct {
	t    ray.Matrix
	tInv ray.Matrix
	m    Material
}

func (o obj) LocalNormalAt(r ray.Vector) ray.Vector {
	return r
}

func (o obj) LocalIntersect(r ray.Ray) []float64 {
	return nil
}

func NewTestShape() Object {
	s := obj{}
	_ = s.SetTransform(ray.DefaultIdentityMatrix())
	m := DefaultMaterial()
	m.Ambient = 1
	s.SetMaterial(m)
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
	return o.m
}

func (o *obj) SetMaterial(m Material) {
	o.m = m
}
