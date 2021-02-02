package object

import (
	"testing"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCube_LocalIntersect(t *testing.T) {
	c := NewCube()
	testCases := []struct {
		name      string
		origin    ray.Vector
		direction ray.Vector
		t1, t2    float64
		miss      bool
	}{
		{
			name:      "+x",
			origin:    ray.NewPoint(5, 0.5, 0),
			direction: ray.NewVec(-1, 0, 0),
			t1:        4,
			t2:        6,
		},
		{
			name:      "-x",
			origin:    ray.NewPoint(-5, 0.5, 0),
			direction: ray.NewVec(1, 0, 0),
			t1:        4,
			t2:        6,
		},
		{
			name:      "+y",
			origin:    ray.NewPoint(0.5, 5, 0),
			direction: ray.NewVec(0, -1, 0),
			t1:        4,
			t2:        6,
		},
		{
			name:      "-y",
			origin:    ray.NewPoint(0.5, -5, 0),
			direction: ray.NewVec(0, 1, 0),
			t1:        4,
			t2:        6,
		},
		{
			name:      "+z",
			origin:    ray.NewPoint(0.5, 0, 5),
			direction: ray.NewVec(0, 0, -1),
			t1:        4,
			t2:        6,
		},
		{
			name:      "-z",
			origin:    ray.NewPoint(0.5, 0, -5),
			direction: ray.NewVec(0, 0, 1),
			t1:        4,
			t2:        6,
		},
		{
			name:      "inside",
			origin:    ray.NewPoint(0, 0.5, 0),
			direction: ray.NewVec(0, 0, 1),
			t1:        -1,
			t2:        1,
		},
		{
			name:      "A ray misses a cube 1",
			miss:      true,
			origin:    ray.NewPoint(-2, 0, 0),
			direction: ray.NewVec(0.2673, 0.5345, 0.8018),
		},
		{
			name:      "A ray misses a cube 2",
			miss:      true,
			origin:    ray.NewPoint(0, -2, 0),
			direction: ray.NewVec(0.8018, 0.2673, 0.5345),
		},
		{
			name:      "A ray misses a cube 3",
			miss:      true,
			origin:    ray.NewPoint(0, 0, -2),
			direction: ray.NewVec(0.5345, 0.8018, 0.2673),
		},
		{
			name:      "A ray misses a cube 4",
			miss:      true,
			origin:    ray.NewPoint(2, 0, 2),
			direction: ray.NewVec(0, 0, -1),
		},
		{
			name:      "A ray misses a cube 5",
			miss:      true,
			origin:    ray.NewPoint(0, 2, 2),
			direction: ray.NewVec(0, -1, 0),
		},
		{
			name:      "A ray misses a cube 6",
			miss:      true,
			origin:    ray.NewPoint(2, 2, 0),
			direction: ray.NewVec(-1, 0, 0),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := ray.NewRayAt(tt.origin, tt.direction)
			xs := c.LocalIntersect(r)
			if tt.miss {
				assert.Len(t, xs, 0)
				return
			}
			require.Len(t, xs, 2)
			assert.Equal(t, tt.t1, xs[0])
			assert.Equal(t, tt.t2, xs[1])
		})
	}
}

func TestCube_LocalNormalAt(t *testing.T) {
	c := NewCube()
	testCases := []struct {
		name           string
		point          ray.Vector
		expectedNormal ray.Vector
	}{
		{
			name:           "1, 0.5, -0.8",
			point:          ray.NewPoint(1, 0.5, -0.8),
			expectedNormal: ray.NewVec(1, 0, 0),
		},
		{
			name:           "-1, -0.2, 0.9",
			point:          ray.NewPoint(-1, -0.2, 0.9),
			expectedNormal: ray.NewVec(-1, 0, 0),
		},
		{
			name:           "-0.4, 1, -0.1",
			point:          ray.NewPoint(-0.4, 1, -0.1),
			expectedNormal: ray.NewVec(0, 1, 0),
		},
		{
			name:           "0.3, -1, -0.7",
			point:          ray.NewPoint(0.3, -1, -0.7),
			expectedNormal: ray.NewVec(0, -1, 0),
		},
		{
			name:           "-0.6, 0.3, 1",
			point:          ray.NewPoint(-0.6, 0.3, 1),
			expectedNormal: ray.NewVec(0, 0, 1),
		},
		{
			name:           "0.4, 0.4, -1",
			point:          ray.NewPoint(0.4, 0.4, -1),
			expectedNormal: ray.NewVec(0, 0, -1),
		},
		{
			name:           "1, 1, 1",
			point:          ray.NewPoint(1, 1, 1),
			expectedNormal: ray.NewVec(1, 0, 0),
		},
		{
			name:           "-1, -1, -1",
			point:          ray.NewPoint(-1, -1, -1),
			expectedNormal: ray.NewVec(-1, 0, 0),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedNormal, c.LocalNormalAt(tt.point))
		})
	}
}
