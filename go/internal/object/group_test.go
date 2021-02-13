package object_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

func TestNewGroup(t *testing.T) {
	g := object.NewGroup()
	assert.Equal(t, ray.DefaultIdentityMatrix(), g.Transform())
	assert.Equal(t, ray.DefaultIdentityMatrixInverse(), g.TransformInverse())
	assert.Empty(t, g.Children)
}

func TestGroup_AddChild(t *testing.T) {
	g := object.NewGroup()
	s := object.NewTestShape()
	g.AddChild(s)
	assert.Equal(t, &g, s.Parent())
	assert.NotNil(t, g.Children)
	assert.Equal(t, []object.Object{s}, g.Children)
}

func TestGroup_LocalIntersect(t *testing.T) {
	t.Run("Intersecting a ray with an empty group", func(t *testing.T) {
		// Given
		g := object.NewGroup()

		// And
		r := ray.NewRayAt(ray.NewPoint(0, 0, 0), ray.NewVec(0, 0, 1))

		// When
		xs := g.LocalIntersect(r)

		// Then
		assert.Empty(t, xs)
	})

	t.Run("Intersecting a ray with a nonempty group", func(t *testing.T) {
		// Given
		g := object.NewGroup()
		s1 := object.DefaultSphere()
		s2 := object.DefaultSphere()
		require.NoError(t, s2.SetTransform(ray.Translation(0, 0, -3)))
		s3 := object.DefaultSphere()
		require.NoError(t, s3.SetTransform(ray.Translation(5, 0, 0)))
		g.AddChild(s1, s2, s3)

		// When
		r := ray.NewRayAt(ray.NewPoint(0, 0, -5), ray.NewVec(0, 0, 1))
		xs := g.LocalIntersect(r)
		require.Len(t, xs, 4)
		assert.Equal(t, s2, xs[0].Obj, "xs 0")
		assert.Equal(t, s2, xs[1].Obj, "xs 1")
		assert.Equal(t, s1, xs[2].Obj, "xs 2")
		assert.Equal(t, s1, xs[3].Obj, "xs 3")
	})

	t.Run("Intersecting a transformed group", func(t *testing.T) {
		// Given
		g := object.NewGroup()
		require.NoError(t, g.SetTransform(ray.Scaling(2, 2, 2)))
		s := object.DefaultSphere()
		require.NoError(t, s.SetTransform(ray.Translation(5, 0, 0)))
		g.AddChild(s)
		r := ray.NewRayAt(ray.NewPoint(10, 0, -10), ray.NewVec(0, 0, 1))
		xs := object.Intersect(&g, r)
		assert.Len(t, xs, 2)
	})
}
