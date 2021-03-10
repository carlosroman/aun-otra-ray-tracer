package object

import (
	"sort"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type Group struct {
	obj
	Children []Object
}

func (g *Group) AddChild(children ...Object) {
	g.Children = append(g.Children, children...)
	for i := range children {
		children[i].SetParent(g)
	}
}

func NewGroup(opts ...Option) Group {
	g := Group{}

	_ = g.SetTransform(ray.DefaultIdentityMatrix())
	g.SetMaterial(DefaultMaterial())
	for i := range opts {
		opts[i].Apply(&g)
	}
	return g
}

func (g Group) LocalNormalAt(worldPoint ray.Vector, _ Intersection) ray.Vector {
	return worldPoint
}

func (g Group) LocalIntersect(r ray.Ray) (xs Intersections) {
	for i := range g.Children {
		is := Intersect(g.Children[i], r)
		xs = append(xs, is...)
	}

	sort.SliceStable(xs, func(i, j int) bool {
		return xs[i].T < xs[j].T
	})
	return xs
}
