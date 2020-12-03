package object_test

import (
	"testing"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

func TestSphere_Intersect(t *testing.T) {
	testCases := []struct {
		name     string
		ray      ray.Ray
		sphere   object.Object
		expected []float64
	}{
		{
			name:     "center",
			ray:      ray.NewRayAt(ray.NewVec(0, 0, -5), ray.NewVec(0, 0, 1)),
			sphere:   object.NewSphere(ray.NewVec(0, 0, 0), 1),
			expected: []float64{4.0, 6.0},
		},
		{
			name:     "tangent",
			ray:      ray.NewRayAt(ray.NewVec(0, 1, -5), ray.NewVec(0, 0, 1)),
			sphere:   object.NewSphere(ray.NewVec(0, 0, 0), 1),
			expected: []float64{5.0, 5.0},
		},
		{
			name:   "misses",
			ray:    ray.NewRayAt(ray.NewVec(0, 2, -5), ray.NewVec(0, 0, 1)),
			sphere: object.NewSphere(ray.NewVec(0, 0, 0), 1),
		},
		{
			name:     "inside",
			ray:      ray.NewRayAt(ray.NewVec(0, 0, 0), ray.NewVec(0, 0, 1)),
			sphere:   object.NewSphere(ray.NewVec(0, 0, 0), 1),
			expected: []float64{-1.0, 1.0},
		},
		{
			name:     "behind",
			ray:      ray.NewRayAt(ray.NewVec(0, 0, 5), ray.NewVec(0, 0, 1)),
			sphere:   object.NewSphere(ray.NewVec(0, 0, 0), 1),
			expected: []float64{-6.0, -4.0},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			intersects := tt.sphere.Intersect(tt.ray)
			if len(tt.expected) < 1 {
				require.Empty(t, intersects)
				return
			}

			assert.NotEmpty(t, intersects)
			require.Len(t, intersects, len(tt.expected))
			for i := range tt.expected {
				assert.Equal(t, tt.expected[i], intersects[i])
			}
		})
	}
}
