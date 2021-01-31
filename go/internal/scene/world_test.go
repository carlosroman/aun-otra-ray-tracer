package scene_test

import (
	"math"
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
		{
			name:              "Shading with a reflective material",
			r:                 ray.NewRayAt(ray.NewPoint(0, 0, -3), ray.NewVec(0, -math.Sqrt(2)/2, math.Sqrt(2)/2)),
			intersectionT:     math.Sqrt(2),
			intersectionShape: 2,
			expected:          object.NewColor(0.87677, 0.92436, 0.82918),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			w := scene.DefaultWorld()
			if tt.light != emptyLight {
				w.AddLight(tt.light)
			}

			if tt.intersectionShape+1 > len(w.Objects()) {
				shape := object.NewPlane()
				material := shape.Material()
				material.Reflective = 0.5
				shape.SetTransform(ray.Translation(0, -1, 0))
				shape.SetMaterial(material)
				w.AddObject(shape)
			}
			shape := w.Objects()[tt.intersectionShape]
			comps := scene.PrepareComputations(scene.Intersection{
				T:   tt.intersectionT,
				Obj: shape,
			}, tt.r)
			actual := scene.ShadeHit(w, comps, 1)
			assertColorEqual(t, tt.expected, actual)
		})
	}

}

func TestReflectedColor(t *testing.T) {
	t.Run("The reflected color for a nonreflective material", func(t *testing.T) {
		w := scene.DefaultWorld()
		shape := w.Objects()[1]
		material := shape.Material()
		material.Ambient = 1
		shape.SetMaterial(material)
		i := scene.Intersection{T: 1, Obj: shape}
		r := ray.NewRayAt(ray.NewPoint(0, 0, 0), ray.NewVec(0, 0, 1))
		comps := scene.PrepareComputations(i, r)
		actual := scene.ReflectedColor(w, comps, 1)
		assertColorEqual(t, object.NewColor(0, 0, 0), actual)
	})

	t.Run("The reflected color for a reflective material", func(t *testing.T) {
		w := scene.DefaultWorld()
		r := ray.NewRayAt(ray.NewPoint(0, 0, -3), ray.NewVec(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
		shape := object.NewPlane()
		material := shape.Material()
		material.Reflective = 0.5
		shape.SetTransform(ray.Translation(0, -1, 0))
		shape.SetMaterial(material)
		w.AddObject(shape)
		i := scene.Intersection{T: math.Sqrt(2), Obj: shape}
		comps := scene.PrepareComputations(i, r)
		actual := scene.ReflectedColor(w, comps, 1)
		assertColorEqual(t, object.NewColor(0.19032, 0.2379, 0.14274), actual)
	})

	t.Run("The reflected color at the maximum recursive depth", func(t *testing.T) {
		w := scene.DefaultWorld()

		r := ray.NewRayAt(ray.NewPoint(0, 0, -3), ray.NewVec(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
		shape := object.NewPlane()
		material := shape.Material()
		material.Reflective = 0.5
		shape.SetTransform(ray.Translation(0, -1, 0))
		shape.SetMaterial(material)
		w.AddObject(shape)

		i := scene.Intersection{T: math.Sqrt(2), Obj: shape}
		comps := scene.PrepareComputations(i, r)
		actual := scene.ReflectedColor(w, comps, 0)
		assertColorEqual(t, object.NewColor(0, 0, 0), actual)
	})
}

func TestRefractedColor(t *testing.T) {
	testCases := []struct {
		name                          string
		r                             ray.Ray
		expectedColor                 object.RGB
		transparency, refractiveIndex float64
		remaining                     int
		intersectionTs                []float64
		idx                           int
	}{
		{
			name:           "The refracted color with an opaque surface",
			r:              ray.NewRayAt(ray.NewPoint(0, 0, -5), ray.NewVec(0, 0, 1)),
			expectedColor:  object.NewColor(0, 0, 0),
			remaining:      5,
			intersectionTs: []float64{4, 6},
			idx:            0,
		},
		{
			name:            "The refracted color at the maximum recursive depth",
			r:               ray.NewRayAt(ray.NewPoint(0, 0, -5), ray.NewVec(0, 0, 1)),
			expectedColor:   object.NewColor(0, 0, 0),
			transparency:    1.0,
			refractiveIndex: 1.5,
			remaining:       0,
			intersectionTs:  []float64{4, 6},
			idx:             0,
		},
		{
			name:            "The refracted color at the maximum recursive depth",
			r:               ray.NewRayAt(ray.NewPoint(0, 0, math.Sqrt(2)/2), ray.NewVec(0, 1, 0)),
			expectedColor:   object.NewColor(0, 0, 0),
			transparency:    1.0,
			refractiveIndex: 1.5,
			remaining:       5,
			intersectionTs:  []float64{-math.Sqrt(2) / 2, math.Sqrt(2) / 2},
			idx:             1, // NOTE: this time you're inside the sphere, so you need to look at the second intersection, xs[1], not xs[0]
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			w := scene.DefaultWorld()
			shape := w.Objects()[0]
			shapeM := shape.Material()
			if tt.refractiveIndex != 0 {
				shapeM.RefractiveIndex = tt.refractiveIndex
			}
			if tt.transparency != 0 {
				shapeM.Transparency = tt.transparency
			}
			shape.SetMaterial(shapeM)
			xs := scene.Intersections{
				{T: tt.intersectionTs[0], Obj: shape},
				{T: tt.intersectionTs[1], Obj: shape},
			}
			comps := scene.PrepareComputations(xs[tt.idx], tt.r, xs...)
			c := scene.RefractedColor(w, comps, tt.remaining)
			assertColorEqual(t, tt.expectedColor, c)
		})
	}

	t.Run("The refracted color with a refracted ray", func(t *testing.T) {

		w := scene.DefaultWorld()

		a := w.Objects()[0]
		aM := a.Material()
		aM.Ambient = 1.0
		aM.Pattern = object.NewTestPattern()
		a.SetMaterial(aM)

		b := w.Objects()[1]
		bM := b.Material()
		bM.Transparency = 1.0
		bM.RefractiveIndex = 1.5
		b.SetMaterial(bM)

		r := ray.NewRayAt(ray.NewPoint(0, 0, 0.1), ray.NewVec(0, 1, 0))

		xs := scene.Intersections{
			{T: -0.9899, Obj: a},
			{T: -0.4899, Obj: b},
			{T: 0.4899, Obj: b},
			{T: 0.9899, Obj: a},
		}

		comps := scene.PrepareComputations(xs[2], r, xs...)
		c := scene.RefractedColor(w, comps, 5)
		t.Log(c)
		assertColorEqual(t, object.NewColor(0, 0.99888, 0.04725), c)
	})
}

func TestShadeHitWithATransparentMaterial(t *testing.T) {
	w := scene.DefaultWorld()

	floor := object.NewPlane()
	floor.SetTransform(ray.Translation(0, -1, 0))
	floorM := floor.Material()
	floorM.Transparency = 0.5
	floorM.RefractiveIndex = 1.5
	floor.SetMaterial(floorM)
	w.AddObject(floor)

	ball := object.DefaultSphere()
	ball.SetTransform(ray.Translation(0, -3.5, -0.5))
	ballM := ball.Material()
	ballM.Color = object.NewColor(1, 0, 0)
	ballM.Ambient = 0.5
	ball.SetMaterial(ballM)
	w.AddObject(ball)

	r := ray.NewRayAt(ray.NewPoint(0, 0, -3), ray.NewVec(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	xs := scene.Intersections{
		{T: math.Sqrt(2), Obj: floor},
	}

	comps := scene.PrepareComputations(xs[0], r, xs...)
	color := scene.ShadeHit(w, comps, 5)
	assertColorEqual(t, object.NewColor(0.93642, 0.68642, 0.68642), color)
}

func TestShadeHitWithAnIntersectionInShadow(t *testing.T) {
	w := scene.NewWorld()
	w.AddLight(object.NewPointLight(ray.NewPoint(0, 0, -10), object.NewColor(1, 1, 1)))
	s1 := object.NewSphere(ray.ZeroPoint, 1)
	s2 := object.NewSphere(ray.ZeroPoint, 1)
	s2.SetTransform(ray.Translation(0, 0, 10))
	w.AddObjects(s1, s2)
	r := ray.NewRayAt(ray.NewPoint(0, 0, 5), ray.NewVec(0, 0, 1))
	i := scene.Intersection{T: 4, Obj: s2}
	comps := scene.PrepareComputations(i, r)
	c := scene.ShadeHit(w, comps, 1)
	assertColorEqual(t, object.NewColor(0.1, 0.1, 0.1), c)
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

			actual := w.ColorAt(tt.r, 1)
			assertColorEqual(t, tt.expected, actual)
		})
	}

	t.Run("With mutually reflective surfaces", func(t *testing.T) {
		w := scene.DefaultWorld()
		w.AddLight(object.NewPointLight(ray.ZeroPoint, object.NewColor(1, 1, 1)))
		lower := object.NewPlane()
		mLower := lower.Material()
		mLower.Reflective = 1
		lower.SetMaterial(mLower)
		lower.SetTransform(ray.Translation(0, -1, 0))
		w.AddObject(lower)
		upper := object.NewPlane()
		mUpper := upper.Material()
		mUpper.Reflective = 1
		upper.SetMaterial(mUpper)
		upper.SetTransform(ray.Translation(0, 1, 0))
		w.AddObject(upper)

		r := ray.NewRayAt(ray.ZeroPoint, ray.NewVec(0, 1, 0))
		actual := w.ColorAt(r, 100)
		t.Log(actual)
		assert.NotNil(t, actual)
	})
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

func TestWorld_IsShadowed(t *testing.T) {

	testCases := []struct {
		name       string
		point      ray.Vector
		isShadowed bool
	}{
		{
			name:       "There is no shadow when nothing is collinear with point and light",
			point:      ray.NewPoint(0, 10, 0),
			isShadowed: false,
		},
		{
			name:       "The shadow when an object is between the point and the light",
			point:      ray.NewPoint(10, -10, 10),
			isShadowed: true,
		},
		{
			name:       "There is no shadow when an object is behind the light",
			point:      ray.NewPoint(-20, 20, -20),
			isShadowed: false,
		},
		{
			name:       "There is no shadow when an object is behind the point",
			point:      ray.NewPoint(-2, 2, -2),
			isShadowed: false,
		},
	}

	w := scene.DefaultWorld()
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actual := w.IsShadowed(tt.point)
			assert.Equal(t, tt.isShadowed, actual)
		})
	}
}

func assertColorEqual(t *testing.T, expected object.RGB, actual object.RGB) {
	assert.InDelta(t, expected.R, actual.R, 0.0001, "R")
	assert.InDelta(t, expected.G, actual.G, 0.0001, "G")
	assert.InDelta(t, expected.B, actual.B, 0.0001, "B")
}
