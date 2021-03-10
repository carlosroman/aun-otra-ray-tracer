package object_test

import (
	"fmt"
	"testing"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTriangle(t *testing.T) {
	tr := object.NewTriangle(
		ray.NewPoint(0, 1, 0), ray.NewPoint(-1, 0, 0), ray.NewPoint(1, 0, 0))
	assert.Equal(t, "e1: -1,-1,0,0, e2: 1,-1,0,0, normal: 0,-0,-1,0", fmt.Sprintf("%v", tr))
}

func TestTriangle_LocalNormalAt(t *testing.T) {

	tr := object.NewTriangle(
		ray.NewPoint(0, 1, 0), ray.NewPoint(-1, 0, 0), ray.NewPoint(1, 0, 0))

	en := ray.Cross(
		ray.NewPoint(1, 0, 0).Subtract(ray.NewPoint(0, 1, 0)),
		ray.NewPoint(-1, 0, 0).Subtract(ray.NewPoint(0, 1, 0))).
		Normalize()

	testCases := []struct {
		name           string
		point          ray.Vector
		expectedNormal ray.Vector
	}{
		{
			name:           "one",
			point:          ray.NewPoint(0, 0.5, 0),
			expectedNormal: en,
		},
		{
			name:           "two",
			point:          ray.NewPoint(-0.5, 0.75, 0),
			expectedNormal: en,
		},
		{
			name:           "three",
			point:          ray.NewPoint(0.5, 0.25, 0),
			expectedNormal: en,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actual := tr.LocalNormalAt(tt.point, object.NoHit)
			assertVec(t, tt.expectedNormal, actual)
		})
	}
}

func TestTriangle_LocalIntersect(t *testing.T) {
	tr := object.NewTriangle(
		ray.NewPoint(0, 1, 0), ray.NewPoint(-1, 0, 0), ray.NewPoint(1, 0, 0))

	testCases := []struct {
		name       string
		expectedXs object.Intersections
		ray        ray.Ray
	}{
		{
			name:       "Intersecting a ray parallel to the triangle",
			expectedXs: nil,
			ray:        ray.NewRayAt(ray.NewPoint(0, -1, -2), ray.NewVec(0, 1, 0)),
		},
		{
			name:       " A ray misses the p1-p3 edge",
			expectedXs: nil,
			ray:        ray.NewRayAt(ray.NewPoint(1, 1, -2), ray.NewVec(0, 0, 1)),
		},
		{
			name:       "A ray misses the p1-p2 edge",
			expectedXs: nil,
			ray:        ray.NewRayAt(ray.NewPoint(-1, -1, -2), ray.NewVec(0, 0, 1)),
		},
		{
			name:       "A ray misses the p2-p3 edge",
			expectedXs: nil,
			ray:        ray.NewRayAt(ray.NewPoint(0, -1, -2), ray.NewVec(0, 0, 1)),
		},
		{
			name: "A ray strikes a triangle",
			expectedXs: object.Intersections{
				{T: 2, Obj: tr}},
			ray: ray.NewRayAt(ray.NewPoint(0, 0.5, -2), ray.NewVec(0, 0, 1)),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actual := tr.LocalIntersect(tt.ray)
			require.Len(t, actual, len(tt.expectedXs))
			for i := range tt.expectedXs {
				assert.Equal(t, tt.expectedXs[i].T, actual[i].T, fmt.Sprintf("Checking xs[%v].T", i))
				assert.Equal(t, tt.expectedXs[i].Obj, actual[i].Obj, fmt.Sprintf("Checking xs[%v].Obj", i))
			}
		})
	}
}
