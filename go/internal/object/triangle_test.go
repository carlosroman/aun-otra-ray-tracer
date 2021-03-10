package object_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

func TestNewTriangle(t *testing.T) {
	tr := getTestTriangle()
	assert.Equal(t, "e1: -1,-1,0,0, e2: 1,-1,0,0, n1: 0,-0,-1,0, n2: 0,-0,-1,0, n3: 0,-0,-1,0", fmt.Sprintf("%v", tr))
}

func TestNewSmoothTriangle(t *testing.T) {
	tr := object.NewSmoothTriangle(
		ray.NewPoint(0, 1, 0), ray.NewPoint(-1, 0, 0), ray.NewPoint(1, 0, 0),
		ray.NewVec(0, 2, 0), ray.NewVec(-2, 0, 0), ray.NewVec(2, 0, 0),
	)
	assert.Equal(t, "e1: -1,-1,0,0, e2: 1,-1,0,0, n1: 0,2,0,0, n2: -2,0,0,0, n3: 2,0,0,0", fmt.Sprintf("%v", tr))
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
	testCases := []struct {
		name       string
		expectedXs object.Intersections
		ray        ray.Ray
		tri        object.Object
	}{
		{
			name:       "Intersecting a ray parallel to the triangle",
			tri:        getTestTriangle(),
			expectedXs: nil,
			ray:        ray.NewRayAt(ray.NewPoint(0, -1, -2), ray.NewVec(0, 1, 0)),
		},
		{
			name:       " A ray misses the p1-p3 edge",
			tri:        getTestTriangle(),
			expectedXs: nil,
			ray:        ray.NewRayAt(ray.NewPoint(1, 1, -2), ray.NewVec(0, 0, 1)),
		},
		{
			name:       "A ray misses the p1-p2 edge",
			tri:        getTestTriangle(),
			expectedXs: nil,
			ray:        ray.NewRayAt(ray.NewPoint(-1, -1, -2), ray.NewVec(0, 0, 1)),
		},
		{
			name:       "A ray misses the p2-p3 edge",
			tri:        getTestTriangle(),
			expectedXs: nil,
			ray:        ray.NewRayAt(ray.NewPoint(0, -1, -2), ray.NewVec(0, 0, 1)),
		},
		{
			name: "A ray strikes a triangle",
			tri:  getTestTriangle(),
			expectedXs: object.Intersections{
				{T: 2, Obj: getTestTriangle(), U: 0.25, V: 0.25}},
			ray: ray.NewRayAt(ray.NewPoint(0, 0.5, -2), ray.NewVec(0, 0, 1)),
		},
		{
			name: "A ray strikes a smooth triangle",
			tri:  getTestSmoothTriangle(),
			expectedXs: object.Intersections{
				{T: 2, Obj: getTestSmoothTriangle(), U: 0.45, V: 0.25}},
			ray: ray.NewRayAt(ray.NewPoint(-0.2, 0.3, -2), ray.NewVec(0, 0, 1)),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.tri.LocalIntersect(tt.ray)
			require.Len(t, actual, len(tt.expectedXs))
			for i := range tt.expectedXs {
				assert.Equal(t, tt.expectedXs[i].T, actual[i].T, fmt.Sprintf("Checking xs[%v].T", i))
				assert.Equal(t, tt.expectedXs[i].Obj, actual[i].Obj, fmt.Sprintf("Checking xs[%v].Obj", i))
				assert.InDelta(t, tt.expectedXs[i].U, actual[i].U, 0.00001, fmt.Sprintf("Checking xs[%v].U", i))
				assert.InDelta(t, tt.expectedXs[i].V, actual[i].V, 0.00001, fmt.Sprintf("Checking xs[%v].V", i))
			}
		})
	}
}

func getTestTriangle() object.Object {
	return object.NewTriangle(
		ray.NewPoint(0, 1, 0), ray.NewPoint(-1, 0, 0), ray.NewPoint(1, 0, 0))
}

func getTestSmoothTriangle() object.Object {
	return object.NewSmoothTriangle(
		ray.NewPoint(0, 1, 0), ray.NewPoint(-1, 0, 0), ray.NewPoint(1, 0, 0),
		ray.NewPoint(0, 1, 0), ray.NewPoint(-1, 0, 0), ray.NewPoint(1, 0, 0),
	)
}

func TestNormalAtSmoothTriangle(t *testing.T) {
	tri := getTestSmoothTriangle()
	xs := object.Intersection{T: 1, Obj: tri, U: 0.45, V: 0.25}
	n := object.NormalAt(xs, ray.NewPoint(0, 0, 0))
	assertVec(t, ray.NewVec(-0.5547, 0.83205, 0), n)
}
