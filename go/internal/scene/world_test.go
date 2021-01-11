package scene_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/scene"
)

func TestNewWorld(t *testing.T) {
	world := scene.NewWorld()
	assert.Empty(t, world.Objects())
}

func TestWorld_AddObject(t *testing.T) {
	world := scene.NewWorld()
	s1, s2 := createTestWorld(world)
	assert.Len(t, world.Objects(), 2)
	assert.Contains(t, world.Objects(), s1)
	assert.Contains(t, world.Objects(), s2)
}

func TestWorld_AddObjects(t *testing.T) {
	world := scene.NewWorld()
	s1, s2 := createTestWorld(scene.NewWorld())
	world.AddObjects(s1, s2)
	assert.Len(t, world.Objects(), 2)
	assert.Contains(t, world.Objects(), s1)
	assert.Contains(t, world.Objects(), s2)
}

func TestIntersect(t *testing.T) {
	world := scene.NewWorld()
	s1, s2 := createTestWorld(world)
	r := ray.NewRayAt(ray.NewPoint(0, 0, -5), ray.NewVec(0, 0, 1))
	intersects := scene.Intersect(world, r)
	require.NotEmpty(t, intersects)
	require.Len(t, intersects, 4)

	assert.Equal(t, float64(4), intersects[0].T)
	assert.Equal(t, s1, intersects[0].Obj)

	assert.Equal(t, 4.5, intersects[1].T)
	assert.Equal(t, s2, intersects[1].Obj)

	assert.Equal(t, 5.5, intersects[2].T)
	assert.Equal(t, s2, intersects[2].Obj)

	assert.Equal(t, float64(6), intersects[3].T)
	assert.Equal(t, s1, intersects[3].Obj)
}

func createTestWorld(world scene.World) (object.Object, object.Object) {
	s1 := object.NewSphere(ray.ZeroPoint, 1)
	world.AddObject(s1)
	s2 := object.NewSphere(ray.ZeroPoint, 0.5)
	world.AddObject(s2)
	return s1, s2
}

func TestShadeHit(t *testing.T) {

	var emptyLight object.PointLight
	testCases := []struct {
		name              string
		r                 ray.Ray
		intersectionT     float64
		intersectionShape int
		expected          object.RGB
		light             object.PointLight
	}{
		{
			name:              "Shading an intersection",
			r:                 ray.NewRayAt(ray.NewPoint(0, 0, -5), ray.NewVec(0, 0, 1)),
			intersectionT:     4,
			intersectionShape: 0,
			expected:          object.NewColor(0.38066, 0.47583, 0.2855),
		},
		{
			name:              "Shading an intersection from the inside",
			r:                 ray.NewRayAt(ray.NewPoint(0, 0, 0), ray.NewVec(0, 0, 1)),
			intersectionT:     0.5,
			intersectionShape: 1,
			light:             object.NewPointLight(ray.NewPoint(0, 0.25, 0), object.NewColor(1, 1, 1)),
			expected:          object.NewColor(0.90498, 0.90498, 0.90498),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			w := scene.DefaultWorld()
			if tt.light != emptyLight {
				w.AddLight(tt.light)
			}

			shape := w.Objects()[tt.intersectionShape]
			comps := scene.PrepareComputations(scene.Intersection{
				T:   tt.intersectionT,
				Obj: shape,
			}, tt.r)
			actual := scene.ShadeHit(w, comps)
			assertColorEqual(t, tt.expected, actual)
		})
	}

}

func TestWorld_ColorAt(t *testing.T) {

	testCases := []struct {
		name     string
		r        ray.Ray
		expected object.RGB
		ambient  float64
	}{
		{
			name:     "The color when a ray misses",
			r:        ray.NewRayAt(ray.NewPoint(0, 0, -5), ray.NewVec(0, 1, 0)),
			expected: object.RGB{},
		},
		{
			name:     "The color when a ray hits",
			r:        ray.NewRayAt(ray.NewPoint(0, 0, -5), ray.NewVec(0, 0, 1)),
			expected: object.NewColor(0.38066, 0.47583, 0.2855),
		},
		{
			name:     "The color with an intersection behind the ray",
			r:        ray.NewRayAt(ray.NewPoint(0, 0, 0.75), ray.NewVec(0, 0, -1)),
			ambient:  1,
			expected: object.NewColor(1, 1, 1),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Given default world
			w := scene.DefaultWorld()
			if tt.ambient > 0 {
				for i := range w.Objects() {
					obj := w.Objects()[i]
					material := obj.Material()
					material.Ambient = tt.ambient
					w.Objects()[i].SetMaterial(material)
				}
			}

			actual := w.ColorAt(tt.r)
			assertColorEqual(t, tt.expected, actual)
		})
	}
}

func TestHit(t *testing.T) {
	testCases := []struct {
		name      string
		hitPoints scene.Intersections
		expected  scene.Intersection
	}{
		{
			name: "All positive",
			hitPoints: scene.Intersections{{
				T: 1,
			}, {
				T: 2,
			}, {
				T: 3,
			}},
			expected: scene.Intersection{T: 1},
		},
		{
			name: "Some negative and positive",
			hitPoints: scene.Intersections{{
				T: -1,
			}, {
				T: 1,
			},
			},
			expected: scene.Intersection{T: 1},
		},
		{
			name: "All negative",
			hitPoints: scene.Intersections{{
				T: -2,
			}, {
				T: -1,
			},
			},
			expected: scene.Intersection{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := scene.Hit(testCase.hitPoints)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func assertColorEqual(t *testing.T, expected object.RGB, actual object.RGB) {
	assert.InDelta(t, expected.R, actual.R, 0.00001, "R")
	assert.InDelta(t, expected.G, actual.G, 0.00001, "G")
	assert.InDelta(t, expected.B, actual.B, 0.00001, "B")
}
