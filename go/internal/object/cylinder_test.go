package object_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/ray"
)

func TestCylinder_LocalIntersect(t *testing.T) {
	testCases := []struct {
		name      string
		origin    ray.Vector
		direction ray.Vector
		t0, t1    float64
		miss      bool
	}{
		{
			name:      "A ray misses a cylinder 1",
			origin:    ray.NewPoint(1, 0, 0),
			direction: ray.NewVec(0, 1, 0),
			miss:      true,
		},
		{
			name:      "A ray misses a cylinder 2",
			origin:    ray.NewPoint(0, 0, 0),
			direction: ray.NewVec(0, 1, 0),
			miss:      true,
		},
		{
			name:      "A ray misses a cylinder 3",
			origin:    ray.NewPoint(0, 0, -5),
			direction: ray.NewVec(1, 1, 1),
			miss:      true,
		},
		{
			name:      "A ray strikes a cylinder 1",
			origin:    ray.NewPoint(1, 0, -5),
			direction: ray.NewVec(0, 0, 1),
			t0:        5,
			t1:        5,
		},
		{
			name:      "A ray strikes a cylinder 2",
			origin:    ray.NewPoint(0, 0, -5),
			direction: ray.NewVec(0, 0, 1),
			t0:        4,
			t1:        6,
		},
		{
			name:      "A ray strikes a cylinder 3",
			origin:    ray.NewPoint(0.5, 0, -5),
			direction: ray.NewVec(0.1, 1, 1),
			t0:        6.80798,
			t1:        7.08872,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			c := object.DefaultCylinder()
			direction := tt.direction.Normalize()
			r := ray.NewRayAt(tt.origin, direction)
			xs := c.LocalIntersect(r)
			if tt.miss {
				assert.Len(t, xs, 0)
				return
			}
			require.Len(t, xs, 2)
			assertIntersection(t, tt.t0, c, xs[0], "t1")
			assertIntersection(t, tt.t1, c, xs[1], "t2")
		})
	}

	nextTestCases := []struct {
		name          string
		point         ray.Vector
		direction     ray.Vector
		closed        bool
		expectedCount int
	}{
		{
			name:          "Intersecting a constrained cylinder 1",
			point:         ray.NewPoint(0, 1.5, 0),
			direction:     ray.NewVec(0.1, 1, 0),
			expectedCount: 0,
		},
		{
			name:          "Intersecting a constrained cylinder 2",
			point:         ray.NewPoint(0, 3, -5),
			direction:     ray.NewVec(0, 0, 1),
			expectedCount: 0,
		},
		{
			name:          "Intersecting a constrained cylinder 3",
			point:         ray.NewPoint(0, 0, -5),
			direction:     ray.NewVec(0, 0, 1),
			expectedCount: 0,
		},
		{
			name:          "Intersecting a constrained cylinder 4",
			point:         ray.NewPoint(0, 2, -5),
			direction:     ray.NewVec(0, 0, 1),
			expectedCount: 0,
		},
		{
			name:          "Intersecting a constrained cylinder 5",
			point:         ray.NewPoint(0, 1, -5),
			direction:     ray.NewVec(0, 0, 1),
			expectedCount: 0,
		},
		{
			name:          "Intersecting a constrained cylinder 6",
			point:         ray.NewPoint(0, 1.5, -2),
			direction:     ray.NewVec(0, 0, 1),
			expectedCount: 2,
		},
		{
			name:          "Intersecting the caps of a closed cylinder 1",
			point:         ray.NewPoint(0, 3, 0),
			direction:     ray.NewVec(0, -1, 0),
			expectedCount: 2,
			closed:        true,
		},
		{
			name:          "Intersecting the caps of a closed cylinder 2",
			point:         ray.NewPoint(0, 3, -2),
			direction:     ray.NewVec(0, -1, 2),
			expectedCount: 2,
			closed:        true,
		},
		{
			name:          "Intersecting the caps of a closed cylinder 3",
			point:         ray.NewPoint(0, 4, -2),
			direction:     ray.NewVec(0, -1, 1),
			expectedCount: 2,
			closed:        true,
		},
		{
			name:          "Intersecting the caps of a closed cylinder 4",
			point:         ray.NewPoint(0, 0, -2),
			direction:     ray.NewVec(0, 1, 2),
			expectedCount: 2,
			closed:        true,
		},
		{
			name:          "Intersecting the caps of a closed cylinder 5",
			point:         ray.NewPoint(0, -1, -2),
			direction:     ray.NewVec(0, 1, 1),
			expectedCount: 2,
			closed:        true,
		},
	}

	for _, tt := range nextTestCases {
		t.Run(tt.name, func(t *testing.T) {
			c := object.NewCylinder(1, 2, tt.closed)
			direction := tt.direction.Normalize()
			r := ray.NewRayAt(tt.point, direction)
			xs := c.LocalIntersect(r)
			require.Len(t, xs, tt.expectedCount)
		})
	}
}

func TestCylinder_LocalNormalAt(t *testing.T) {

	testCases := []struct {
		name           string
		point          ray.Vector
		closed         bool
		expectedNormal ray.Vector
	}{
		{
			name:           "Normal vector on a cylinder 1",
			point:          ray.NewPoint(1, 0, 0),
			expectedNormal: ray.NewVec(1, 0, 0),
		},
		{
			name:           "Normal vector on a cylinder 2",
			point:          ray.NewPoint(0, 5, -1),
			expectedNormal: ray.NewVec(0, 0, -1),
		},
		{
			name:           "Normal vector on a cylinder 3",
			point:          ray.NewPoint(0, -2, 1),
			expectedNormal: ray.NewVec(0, 0, 1),
		},
		{
			name:           "Normal vector on a cylinder 4",
			point:          ray.NewPoint(-1, 1, 0),
			expectedNormal: ray.NewVec(-1, 0, 0),
		},
		{
			name:           "The normal vector on a cylinder's end caps 1",
			closed:         true,
			point:          ray.NewPoint(0, 1, 0),
			expectedNormal: ray.NewVec(0, -1, 0),
		},
		{
			name:           "The normal vector on a cylinder's end caps 2",
			closed:         true,
			point:          ray.NewPoint(0.5, 1, 0),
			expectedNormal: ray.NewVec(0, -1, 0),
		},
		{
			name:           "The normal vector on a cylinder's end caps 3",
			closed:         true,
			point:          ray.NewPoint(0, 1, 0.5),
			expectedNormal: ray.NewVec(0, -1, 0),
		},
		{
			name:           "The normal vector on a cylinder's end caps 4",
			closed:         true,
			point:          ray.NewPoint(0, 2, 0),
			expectedNormal: ray.NewVec(0, 1, 0),
		},
		{
			name:           "The normal vector on a cylinder's end caps 5",
			closed:         true,
			point:          ray.NewPoint(0.5, 2, 0),
			expectedNormal: ray.NewVec(0, 1, 0),
		},
		{
			name:           "The normal vector on a cylinder's end caps 6",
			closed:         true,
			point:          ray.NewPoint(0, 2, 0.5),
			expectedNormal: ray.NewVec(0, 1, 0),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var c object.Object
			if tt.closed {
				c = object.NewCylinder(1, 2, tt.closed)
			} else {
				c = object.DefaultCylinder()
			}
			assert.Equal(t, tt.expectedNormal, c.LocalNormalAt(tt.point, object.NoHit))
		})
	}
}
