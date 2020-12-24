package object_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

func TestSphere_Intersect(t *testing.T) {
	testCases := []struct {
		name      string
		ray       ray.Ray
		sphere    object.Object
		transform ray.Matrix
		expected  []float64
	}{
		{
			name:     "center",
			ray:      ray.NewRayAt(ray.NewPoint(0, 0, -5), ray.NewVec(0, 0, 1)),
			sphere:   object.NewSphere(ray.NewPoint(0, 0, 0), 1),
			expected: []float64{4.0, 6.0},
		},
		{
			name:     "tangent",
			ray:      ray.NewRayAt(ray.NewPoint(0, 1, -5), ray.NewVec(0, 0, 1)),
			sphere:   object.NewSphere(ray.NewPoint(0, 0, 0), 1),
			expected: []float64{5.0, 5.0},
		},
		{
			name:   "misses",
			ray:    ray.NewRayAt(ray.NewPoint(0, 2, -5), ray.NewVec(0, 0, 1)),
			sphere: object.NewSphere(ray.NewPoint(0, 0, 0), 1),
		},
		{
			name:     "inside",
			ray:      ray.NewRayAt(ray.NewPoint(0, 0, 0), ray.NewVec(0, 0, 1)),
			sphere:   object.NewSphere(ray.NewPoint(0, 0, 0), 1),
			expected: []float64{-1.0, 1.0},
		},
		{
			name:     "behind",
			ray:      ray.NewRayAt(ray.NewPoint(0, 0, 5), ray.NewVec(0, 0, 1)),
			sphere:   object.NewSphere(ray.NewPoint(0, 0, 0), 1),
			expected: []float64{-6.0, -4.0},
		},
		{
			name:      "scaled",
			ray:       ray.NewRayAt(ray.NewPoint(0, 0, -5), ray.NewVec(0, 0, 1)),
			sphere:    object.NewSphere(ray.NewPoint(0, 0, 0), 1),
			transform: ray.Scaling(2, 2, 2),
			expected:  []float64{3, 7},
		},
		{
			name:      "translated",
			ray:       ray.NewRayAt(ray.NewPoint(0, 0, -5), ray.NewVec(0, 0, 1)),
			sphere:    object.NewSphere(ray.NewPoint(0, 0, 0), 1),
			transform: ray.Translation(5, 0, 0),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.transform != nil {
				tt.sphere.SetTransform(tt.transform)
			}
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

func TestSphere_Transform(t *testing.T) {
	obj := object.NewSphere(ray.NewPoint(0, 0, 0), 1)
	translation := ray.Translation(2, 3, 4)
	obj.SetTransform(translation)
	assert.Equal(t, translation, obj.Transform())
}

func TestNewSphere(t *testing.T) {
	obj := object.NewSphere(ray.NewPoint(0, 0, 0), 1)
	assert.Equal(t, ray.DefaultIdentityMatrix(), obj.Transform())
}

func TestSphere_NormalAt(t *testing.T) {

	testCases := []struct {
		name     string
		point    ray.Vector
		sphere   object.Object
		expected ray.Vector
	}{
		{
			name:     "X axis",
			point:    ray.NewVec(1, 0, 0),
			sphere:   object.NewSphere(ray.NewPoint(0, 0, 0), 1),
			expected: ray.NewVec(1, 0, 0),
		},
		{
			name:     "Y axis",
			point:    ray.NewVec(0, 1, 0),
			sphere:   object.NewSphere(ray.NewPoint(0, 0, 0), 1),
			expected: ray.NewVec(0, 1, 0),
		},
		{
			name:     "Z axis",
			point:    ray.NewVec(0, 0, 1),
			sphere:   object.NewSphere(ray.NewPoint(0, 0, 0), 1),
			expected: ray.NewVec(0, 0, 1),
		},
		{
			name:     "nonaxial",
			point:    ray.NewVec(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3),
			sphere:   object.NewSphere(ray.NewPoint(0, 0, 0), 1),
			expected: ray.NewVec(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.sphere.NormalAt(tt.point)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.expected, actual.Normalize())
		})
	}
}
