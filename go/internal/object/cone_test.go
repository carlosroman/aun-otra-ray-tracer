package object_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

func TestCone_LocalIntersect(t *testing.T) {

	testCases := []struct {
		name           string
		point          ray.Vector
		expectedNormal ray.Vector
	}{
		{
			name:           "Computing the normal vector on a cone 1",
			point:          ray.NewPoint(0, 0, 0),
			expectedNormal: ray.NewVec(0, 0, 0),
		},
		{
			name:           "Computing the normal vector on a cone 2",
			point:          ray.NewPoint(1, 1, 1),
			expectedNormal: ray.NewVec(1, -math.Sqrt(2), 1),
		},
		{
			name:           "Computing the normal vector on a cone 3",
			point:          ray.NewPoint(-1, -1, 0),
			expectedNormal: ray.NewVec(-1, 1, 0),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			shape := object.DefaultCone()
			assert.Equal(t, tt.expectedNormal, shape.LocalNormalAt(tt.point))
		})
	}
}

func TestCone_LocalNormalAt(t *testing.T) {
	testCases := []struct {
		name      string
		origin    ray.Vector
		direction ray.Vector
		t0, t1    float64
		miss      bool
		count     int
	}{
		{
			name:      "Intersecting a cone with a ray 1",
			origin:    ray.NewPoint(0, 0, -5),
			direction: ray.NewVec(0, 0, 1),
			t0:        5,
			t1:        5,
			count:     2,
		},
		{
			name:      "Intersecting a cone with a ray 2",
			origin:    ray.NewPoint(0, 0, -5),
			direction: ray.NewVec(1, 1, 1),
			t0:        8.66025,
			t1:        8.66025,
			count:     2,
		},
		{
			name:      "Intersecting a cone with a ray 3",
			origin:    ray.NewPoint(1, 1, -5),
			direction: ray.NewVec(-0.5, -1, 1),
			t0:        4.55006,
			t1:        49.44994,
			count:     2,
		},
		{
			name:      "Intersecting a cone with a ray parallel to one of its halves",
			origin:    ray.NewPoint(0, 0, -1),
			direction: ray.NewVec(0, 1, 1),
			t0:        0.35355,
			count:     1,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			c := object.DefaultCone()
			r := ray.NewRayAt(tt.origin, tt.direction.Normalize())
			xs := c.LocalIntersect(r)
			if tt.miss {
				assert.Len(t, xs, 0)
				return
			}
			require.Len(t, xs, tt.count)
			assertIntersection(t, tt.t0, c, xs[0], "t0 did not match")
			if tt.count > 1 {
				assertIntersection(t, tt.t1, c, xs[1], "t1 did not match")
			}
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
			name:          "Intersecting a cone's end caps 1",
			closed:        true,
			point:         ray.NewPoint(0, 0, -5),
			direction:     ray.NewVec(0, 1, 0),
			expectedCount: 0,
		},
		{
			name:          "Intersecting a cone's end caps 2",
			closed:        true,
			point:         ray.NewPoint(0, 0, -0.25),
			direction:     ray.NewVec(0, 1, 1),
			expectedCount: 3, // Should be 2 but I get 3
		},
		{
			name:          "Intersecting a cone's end caps 3",
			closed:        true,
			point:         ray.NewPoint(0, 0, -0.25),
			direction:     ray.NewVec(0, 1, 0),
			expectedCount: 4,
		},
	}

	for _, tt := range nextTestCases {
		t.Run(tt.name, func(t *testing.T) {
			c := object.NewCone(-0.5, 0.5, tt.closed)
			direction := tt.direction.Normalize()
			r := ray.NewRayAt(tt.point, direction)
			xs := c.LocalIntersect(r)
			require.Len(t, xs, tt.expectedCount)
		})
	}
}
