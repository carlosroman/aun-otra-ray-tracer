package scene_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/ray"
	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/scene"
)

func TestPrepareComputations(t *testing.T) {
	testCases := []struct {
		name                                         string
		r                                            ray.Ray
		expectedInside                               bool
		expectedNormalv, expectedEyev, expectedPoint ray.Vector
		expectedObject                               object.Object
		intersectT                                   float64
		expectedObjectTranslation                    ray.Matrix
		expectedOverPoint                            ray.Vector
		expectedUnderPoint                           ray.Vector
		xsU, xsV                                     float64
	}{
		{
			name:               "when an intersection occurs on the outside",
			r:                  ray.NewRayAt(ray.NewPoint(0, 0, -5), ray.NewVec(0, 0, 1)),
			intersectT:         4,
			expectedObject:     object.NewSphere(ray.ZeroPoint, 1),
			expectedPoint:      ray.NewPoint(0, 0, -1),
			expectedEyev:       ray.NewVec(0, 0, -1),
			expectedNormalv:    ray.NewVec(0, 0, -1),
			expectedOverPoint:  ray.NewPoint(0, 0, -1.00000001),
			expectedUnderPoint: ray.NewPoint(0, 0, -0.99999999),
		},
		{
			name:           "when an intersection occurs on the inside",
			r:              ray.NewRayAt(ray.NewPoint(0, 0, 0), ray.NewVec(0, 0, 1)),
			intersectT:     1,
			expectedObject: object.NewSphere(ray.ZeroPoint, 1),
			expectedPoint:  ray.NewPoint(0, 0, 1),
			expectedEyev:   ray.NewVec(0, 0, -1),
			//  normal would have been (0, 0, 1), but is inverted
			expectedNormalv:    ray.NewVec(0, 0, -1),
			expectedOverPoint:  ray.NewPoint(0, 0, 0.99999999),
			expectedUnderPoint: ray.NewPoint(0, 0, 1.00000001),
			expectedInside:     true,
		},
		{
			name:                      "The hit should offset the point",
			r:                         ray.NewRayAt(ray.NewPoint(0, 0, -5), ray.NewVec(0, 0, 1)),
			intersectT:                5,
			expectedObject:            object.NewSphere(ray.ZeroPoint, 1),
			expectedObjectTranslation: ray.Translation(0, 0, 1),
			expectedPoint:             ray.NewPoint(0, 0, 0),
			expectedEyev:              ray.NewVec(0, 0, -1),
			expectedNormalv:           ray.NewVec(0, 0, -1),
			expectedOverPoint:         ray.NewPoint(0, 0, -0.00000001),
			expectedUnderPoint:        ray.NewPoint(0, 0, 0.00000001),
		},
		{
			name:                      "The under point is offset below the surface",
			r:                         ray.NewRayAt(ray.NewPoint(0, 0, -5), ray.NewVec(0, 0, 1)),
			intersectT:                5,
			expectedObject:            object.DefaultGlassSphere(),
			expectedObjectTranslation: ray.Translation(0, 0, 1),
			expectedPoint:             ray.NewPoint(0, 0, 0),
			expectedEyev:              ray.NewVec(0, 0, -1),
			expectedNormalv:           ray.NewVec(0, 0, -1),
			expectedOverPoint:         ray.NewPoint(0, 0, -0.00000001),
			expectedUnderPoint:        ray.NewPoint(0, 0, 0.00000001),
		},
		{
			name:       "Preparing the normal on a smooth triangle",
			r:          ray.NewRayAt(ray.NewPoint(-0.2, 0.3, -2), ray.NewVec(0, 0, 1)),
			intersectT: 5,
			expectedObject: object.NewSmoothTriangle(
				ray.NewPoint(0, 1, 0), ray.NewPoint(-1, 0, 0), ray.NewPoint(1, 0, 0),
				ray.NewPoint(0, 1, 0), ray.NewPoint(-1, 0, 0), ray.NewPoint(1, 0, 0),
			),
			expectedObjectTranslation: ray.Translation(0, 0, 1),
			expectedPoint:             ray.NewPoint(-0.2, 0.3, 3),
			expectedEyev:              ray.NewVec(0, 0, -1),
			expectedNormalv:           ray.NewVec(-0.5547, 0.83205, 0),
			expectedOverPoint:         ray.NewPoint(-0.2, 0.3, 3),
			expectedUnderPoint:        ray.NewPoint(-0.2, 0.3, 3),
			xsU:                       0.45,
			xsV:                       0.25,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedObjectTranslation != nil {
				require.NoError(t, tt.expectedObject.SetTransform(tt.expectedObjectTranslation))
			}
			// Given
			i := object.Intersection{
				T:   tt.intersectT,
				Obj: tt.expectedObject,
				U:   tt.xsU,
				V:   tt.xsV,
			}
			// When
			actual := scene.PrepareComputations(i, tt.r)

			// Then
			assert.Equal(t, tt.intersectT, actual.Intersect())
			assert.Equal(t, tt.expectedObject, actual.Object())
			assertVectorEqual(t, tt.expectedPoint, actual.Point())
			assertVectorEqual(t, tt.expectedEyev, actual.Eyev())
			assertVectorEqual(t, tt.expectedNormalv, actual.Normalv())
			assert.Equal(t, tt.expectedInside, actual.Inside())
			assertVectorEqual(t, tt.expectedOverPoint, actual.OverPoint())
			assertVectorEqual(t, tt.expectedUnderPoint, actual.UnderPoint())
			assert.True(t, actual.Point().GetZ() <= actual.UnderPoint().GetZ())

			assert.True(t, actual.Point().GetZ() >= actual.OverPoint().GetZ(),
				fmt.Sprintf("Got %v > %v", actual.Point().GetZ(), actual.OverPoint().GetZ()))

			// Does not seem to always be true :/
			//  assert.True(t, actual.OverPoint().GetZ() < -0.00000001/2,
			//	fmt.Sprintf("Got %v < %v", actual.OverPoint().GetZ(), -0.00000001/2))
		})
	}

	t.Run("Finding n1 and n2 at various intersections", func(t *testing.T) {
		// GIVEN
		a := object.DefaultGlassSphere()
		require.NoError(t, a.SetTransform(ray.Scaling(2, 2, 2)))
		aM := a.Material()
		aM.RefractiveIndex = 1.5
		a.SetMaterial(aM)

		b := object.DefaultGlassSphere()
		require.NoError(t, b.SetTransform(ray.Translation(0, 0, -0.25)))
		bM := b.Material()
		bM.RefractiveIndex = 2.0
		b.SetMaterial(bM)

		c := object.DefaultGlassSphere()
		require.NoError(t, c.SetTransform(ray.Translation(0, 0, 0.25)))
		cM := c.Material()
		cM.RefractiveIndex = 2.5
		c.SetMaterial(cM)

		r := ray.NewRayAt(ray.NewPoint(0, 0, -4), ray.NewVec(0, 0, 1))
		xs := object.Intersections{
			object.Intersection{T: 2, Obj: a},
			object.Intersection{T: 2.75, Obj: b},
			object.Intersection{T: 3.25, Obj: c},
			object.Intersection{T: 4.75, Obj: b},
			object.Intersection{T: 5.25, Obj: c},
			object.Intersection{T: 6, Obj: a},
		}

		// WHEN
		examples := []struct {
			index  int
			n1, n2 float64
		}{
			{index: 0, n1: 1.0, n2: 1.5},
			{index: 1, n1: 1.5, n2: 2.0},
			{index: 2, n1: 2.0, n2: 2.5},
			{index: 3, n1: 2.5, n2: 2.5},
			{index: 4, n1: 2.5, n2: 1.5},
			{index: 5, n1: 1.5, n2: 1.0},
		}

		for _, example := range examples {
			comps := scene.PrepareComputations(xs[example.index], r, xs...)
			assert.NotNil(t, comps)
			assert.Equal(t, example.n1, comps.N1(), "%v N1 did not match", example.index)
			assert.Equal(t, example.n2, comps.N2(), "%v N2 did not match", example.index)
		}
	})
}

func TestComputation_Reflectv(t *testing.T) {
	shape := object.NewPlane()
	r := ray.NewRayAt(
		ray.NewPoint(0, 1, -1),
		ray.NewVec(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	i := object.Intersection{
		T:   math.Sqrt(2),
		Obj: shape,
	}
	actual := scene.PrepareComputations(i, r)
	assert.Equal(t, ray.NewVec(0, math.Sqrt(2)/2, math.Sqrt(2)/2), actual.Reflectv())
}

func TestSchlick(t *testing.T) {

	shape := object.DefaultGlassSphere()
	testCases := []struct {
		name     string
		xs       object.Intersections
		expected float64
		r        ray.Ray
		idx      int
	}{
		{
			name: "The Schlick approximation under total internal reflection",
			xs: object.Intersections{
				{T: -math.Sqrt(2) / 2, Obj: shape},
				{T: math.Sqrt(2) / 2, Obj: shape},
			},
			r:        ray.NewRayAt(ray.NewPoint(0, 0, math.Sqrt(2)/2), ray.NewVec(0, 1, 0)),
			idx:      1,
			expected: 1.0,
		},
		{
			name: "The Schlick approximation with a perpendicular viewing angle",
			xs: object.Intersections{
				{T: -1, Obj: shape},
				{T: 1, Obj: shape},
			},
			r:        ray.NewRayAt(ray.NewPoint(0, 0, 0), ray.NewVec(0, 1, 0)),
			idx:      1,
			expected: 0.04,
		},
		{
			name: "The Schlick approximation with small angle and n2 > n1",
			xs: object.Intersections{
				{T: 1.8589, Obj: shape},
			},
			r:        ray.NewRayAt(ray.NewPoint(0, 0.99, -2), ray.NewVec(0, 0, 1)),
			idx:      0,
			expected: 0.48873,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			comps := scene.PrepareComputations(tt.xs[tt.idx], tt.r, tt.xs...)
			reflectance := scene.Schlick(comps)
			assert.InDelta(t, tt.expected, reflectance, 0.00001)
		})
	}
}
