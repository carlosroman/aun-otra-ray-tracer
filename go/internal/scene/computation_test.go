package scene_test

import (
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
	}{
		{
			name:            "when an intersection occurs on the outside",
			r:               ray.NewRayAt(ray.NewPoint(0, 0, -5), ray.NewVec(0, 0, 1)),
			intersectT:      4,
			expectedObject:  object.NewSphere(ray.ZeroPoint, 1),
			expectedPoint:   ray.NewPoint(0, 0, -1),
			expectedEyev:    ray.NewVec(0, 0, -1),
			expectedNormalv: ray.NewVec(0, 0, -1),
		}, {
			name:           "when an intersection occurs on the inside",
			r:              ray.NewRayAt(ray.NewPoint(0, 0, 0), ray.NewVec(0, 0, 1)),
			intersectT:     1,
			expectedObject: object.NewSphere(ray.ZeroPoint, 1),
			expectedPoint:  ray.NewPoint(0, 0, 1),
			expectedEyev:   ray.NewVec(0, 0, -1),
			//  normal would have been (0, 0, 1), but is inverted
			expectedNormalv: ray.NewVec(0, 0, -1),
			expectedInside:  true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
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
		})
	}
}
