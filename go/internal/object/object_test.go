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
		n := object.NormalAt(object.Intersection{Obj: s}, ray.NewPoint(0, 1.70711, -0.70711))
		assertVec(t, ray.NewVec(0, 0.70711, -0.70711), n)
	})

	t.Run("Computing the normal on a transformed shape", func(t *testing.T) {
		s := object.NewTestShape()
		require.NoError(t, s.SetTransform(ray.Scaling(1, 0.5, 1).Multiply(ray.Rotation(ray.Z, math.Pi/5))))
		n := object.NormalAt(object.Intersection{Obj: s}, ray.NewPoint(0, math.Sqrt(2)/2, -math.Sqrt(2)/2))
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

func TestObj_Material(t *testing.T) {
	m := object.DefaultMaterial()
	m.Color = object.Red
	testCases := []struct {
		name     string
		obj      object.Object
		expected object.Material
	}{
		{
			name: "No parent",
			obj: object.NewTestShape(
				object.WithMaterial(m),
			),
			expected: m,
		},
		{
			name: "Has parent",
			obj: object.NewGroup(
				object.WithMaterial(m),
				object.WithChildren(object.NewTestShape()),
			).Object(),
			expected: m,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.obj.Material())
		})
	}
}

func assertIntersection(t *testing.T, expectedT float64, expectedObj object.Object, actual object.Intersection, msg ...string) {
	assert.InDelta(t, expectedT, actual.T, 0.00001, msg)
	assert.Equal(t, expectedObj, actual.Obj, msg)
}

func TestObj_WorldToObject(t *testing.T) {
	g1 := object.NewGroup()
	require.NoError(t, g1.SetTransform(ray.Rotation(ray.Y, math.Pi/2)))

	g2 := object.NewGroup()
	require.NoError(t, g2.SetTransform(ray.Scaling(2, 2, 2)))
	g1.AddChild(&g2)

	s := object.DefaultSphere()
	require.NoError(t, s.SetTransform(ray.Translation(5, 0, 0)))
	g2.AddChild(s)

	p := s.WorldToObject(ray.NewPoint(-2, 0, -10))
	assertVec(t, ray.NewPoint(0, 0, -1), p)
}

func TestObj_NormalToWorld(t *testing.T) {
	g1 := object.NewGroup()
	require.NoError(t, g1.SetTransform(ray.Rotation(ray.Y, math.Pi/2)))

	g2 := object.NewGroup()
	require.NoError(t, g2.SetTransform(ray.Scaling(1, 2, 3)))
	g1.AddChild(&g2)

	s := object.DefaultSphere()
	require.NoError(t, s.SetTransform(ray.Translation(5, 0, 0)))
	g2.AddChild(s)

	p := s.NormalToWorld(ray.NewVec(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3))
	assertVec(t, ray.NewVec(0.28571, 0.42857, -0.85714), p)
}
