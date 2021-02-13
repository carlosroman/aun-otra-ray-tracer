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

func NewGroup() Group {
	g := Group{}
	_ = g.SetTransform(ray.DefaultIdentityMatrix())
	return g
}

func (g Group) LocalNormalAt(r ray.Vector) ray.Vector {
	return r
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
