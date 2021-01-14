package scene_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/scene"
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
	}{
		{
			name:              "when an intersection occurs on the outside",
			r:                 ray.NewRayAt(ray.NewPoint(0, 0, -5), ray.NewVec(0, 0, 1)),
			intersectT:        4,
			expectedObject:    object.NewSphere(ray.ZeroPoint, 1),
			expectedPoint:     ray.NewPoint(0, 0, -1),
			expectedEyev:      ray.NewVec(0, 0, -1),
			expectedNormalv:   ray.NewVec(0, 0, -1),
			expectedOverPoint: ray.NewPoint(0, 0, -1.00000001),
		},
		{
			name:           "when an intersection occurs on the inside",
			r:              ray.NewRayAt(ray.NewPoint(0, 0, 0), ray.NewVec(0, 0, 1)),
			intersectT:     1,
			expectedObject: object.NewSphere(ray.ZeroPoint, 1),
			expectedPoint:  ray.NewPoint(0, 0, 1),
			expectedEyev:   ray.NewVec(0, 0, -1),
			//  normal would have been (0, 0, 1), but is inverted
			expectedNormalv:   ray.NewVec(0, 0, -1),
			expectedOverPoint: ray.NewPoint(0, 0, 0.99999999),
			expectedInside:    true,
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
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedObjectTranslation != nil {
				tt.expectedObject.SetTransform(tt.expectedObjectTranslation)
			}
			// Given
			i := scene.Intersection{
				T:   tt.intersectT,
				Obj: tt.expectedObject,
			}
			// When
			actual := scene.PrepareComputations(i, tt.r)

			// Then
			assert.Equal(t, tt.intersectT, actual.Intersect())
			assert.Equal(t, tt.expectedObject, actual.Object())
			assert.Equal(t, tt.expectedPoint, actual.Point())
			assert.Equal(t, tt.expectedEyev, actual.Eyev())
			assert.Equal(t, tt.expectedNormalv, actual.Normalv())
			assert.Equal(t, tt.expectedInside, actual.Inside())
			assert.Equal(t, tt.expectedOverPoint, actual.OverPoint())

			assert.True(t, actual.Point().GetZ() > actual.OverPoint().GetZ(),
				fmt.Sprintf("Got %v > %v", actual.Point().GetZ(), actual.OverPoint().GetZ()))

			// Does not seem to always be true :/
			//assert.True(t, actual.OverPoint().GetZ() < -0.00000001/2,
			//	fmt.Sprintf("Got %v < %v", actual.OverPoint().GetZ(), -0.00000001/2))
		})
	}
}
