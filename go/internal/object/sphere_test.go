package object_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

func TestSphere_LocalIntersect(t *testing.T) {

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
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.transform != nil {
				tt.sphere.SetTransform(tt.transform)
			}
			intersects := tt.sphere.LocalIntersect(tt.ray)
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

func TestIntersect(t *testing.T) {
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
			name:   "misses",
			ray:    ray.NewRayAt(ray.NewPoint(0, 2, -5), ray.NewVec(0, 0, 1)),
			sphere: object.NewSphere(ray.NewPoint(0, 0, 0), 1),
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
			intersects := object.Intersect(tt.sphere, tt.ray)
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
	obj := object.DefaultSphere()
	translation := ray.Translation(2, 3, 4)
	obj.SetTransform(translation)
	assert.Equal(t, translation, obj.Transform())
}

func TestSphere_Material(t *testing.T) {
	obj := object.DefaultSphere()
	material := object.Material{
		Color:     object.NewColor(1, 2, 3),
		Ambient:   0,
		Diffuse:   0,
		Specular:  0,
		Shininess: 0,
	}
	obj.SetMaterial(material)
	assert.Equal(t, material, obj.Material())
}

func TestNewSphere(t *testing.T) {
	obj := object.DefaultSphere()
	assert.Equal(t, ray.DefaultIdentityMatrix(), obj.Transform())
}

func TestSphere_NormalAt(t *testing.T) {

	testCases := []struct {
		name      string
		point     ray.Vector
		sphere    object.Object
		expected  ray.Vector
		transform ray.Matrix
	}{
		{
			name:     "X axis",
			point:    ray.NewPoint(1, 0, 0),
			sphere:   object.NewSphere(ray.NewPoint(0, 0, 0), 1),
			expected: ray.NewVec(1, 0, 0),
		},
		{
			name:     "Y axis",
			point:    ray.NewPoint(0, 1, 0),
			sphere:   object.NewSphere(ray.NewPoint(0, 0, 0), 1),
			expected: ray.NewVec(0, 1, 0),
		},
		{
			name:     "Z axis",
			point:    ray.NewPoint(0, 0, 1),
			sphere:   object.NewSphere(ray.NewPoint(0, 0, 0), 1),
			expected: ray.NewVec(0, 0, 1),
		},
		{
			name:     "nonaxial",
			point:    ray.NewPoint(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3),
			sphere:   object.NewSphere(ray.NewPoint(0, 0, 0), 1),
			expected: ray.NewVec(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3),
		},
		{
			name:      "translation",
			point:     ray.NewPoint(0, 1.70711, -0.70711),
			sphere:    object.NewSphere(ray.NewPoint(0, 0, 0), 1),
			expected:  ray.NewVec(0, 0.70711, -0.70711),
			transform: ray.Translation(0, 1, 0),
		},
		{
			name:      "transformed",
			point:     ray.NewPoint(0, math.Sqrt(2)/2, -math.Sqrt(2)/2),
			sphere:    object.NewSphere(ray.NewPoint(0, 0, 0), 1),
			expected:  ray.NewVec(0, 0.97014, -0.24254),
			transform: ray.Scaling(1, 0.5, 1).Multiply(ray.Rotation(ray.Z, math.Pi/5)),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.transform != nil {
				tt.sphere.SetTransform(tt.transform)
			}
			actual := tt.sphere.NormalAt(tt.point)
			assertVec(t, tt.expected, actual)
			assertVec(t, tt.expected, actual.Normalize())
		})
	}
}

func assertVec(t *testing.T, expected, actual ray.Vector) {
	assert.InDelta(t, expected.GetX(), actual.GetX(), 0.00001, "X did not match")
	assert.InDelta(t, expected.GetY(), actual.GetY(), 0.00001, "Y did not match")
	assert.InDelta(t, expected.GetZ(), actual.GetZ(), 0.00001, "Z did not match")
	assert.InDelta(t, expected.GetW(), actual.GetW(), 0.00001, "W did not match")
}
