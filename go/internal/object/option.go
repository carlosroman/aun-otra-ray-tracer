package object

import "github.com/carlosroman/aun-otra-ray-tracer/go/internal/ray"

type OptionFunc func(object Object)

func (of OptionFunc) Apply(o Object) {
	of(o)
}

type Option interface {
	Apply(object Object)
}

func WithTransform(t ray.Matrix) Option {
	return OptionFunc(func(o Object) {
		_ = o.SetTransform(t)
	})
}

func WithMaterial(m Material) Option {
	return OptionFunc(func(o Object) {
		o.SetMaterial(m)
	})
}
