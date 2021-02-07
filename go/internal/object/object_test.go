package object_test

import (
	"math"
	"testing"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
