package object_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

func TestNewTestShape(t *testing.T) {

	t.Run("Intersecting a scaled shape with a ray", func(t *testing.T) {
		s := object.NewTestShape()
		r := ray.NewRayAt(ray.NewPoint(0, 0, -5), ray.NewVec(0, 0, 1))
		require.NoError(t, s.SetTransform(ray.Scaling(2, 2, 2)))

		xs := object.Intersect(s, r)
		assert.Nil(t, xs)
	})

	t.Run("Computing the normal on a translated shape", func(t *testing.T) {
		s := object.NewTestShape()
		require.NoError(t, s.SetTransform(ray.Translation(0, 1, 0)))
		n := object.NormalAt(s, ray.NewPoint(0, 1.70711, -0.70711))
		assertVec(t, ray.NewVec(0, 0.70711, -0.70711), n)
	})

	t.Run("Computing the normal on a transformed shape", func(t *testing.T) {
		s := object.NewTestShape()
		require.NoError(t, s.SetTransform(ray.Scaling(1, 0.5, 1).Multiply(ray.Rotation(ray.Z, math.Pi/5))))
		n := object.NormalAt(s, ray.NewPoint(0, math.Sqrt(2)/2, -math.Sqrt(2)/2))
		assertVec(t, ray.NewVec(0, 0.97014, -0.24254), n)
	})
}

func TestObj_Parent(t *testing.T) {
	s := object.NewTestShape()
	require.Nil(t, s.Parent())
	g := object.NewGroup()
	s.SetParent(&g)
	actual := s.Parent()
	assert.NotNil(t, actual)
	assert.Equal(t, &g, actual)
}

func assertIntersection(t *testing.T, expectedT float64, expectedObj object.Object, actual object.Intersection, msg ...string) {
	assert.InDelta(t, expectedT, actual.T, 0.00001, msg)
	assert.Equal(t, expectedObj, actual.Obj, msg)
}
