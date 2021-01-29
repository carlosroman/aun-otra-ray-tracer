package object_test

import (
	"testing"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
	"github.com/stretchr/testify/assert"
)

func TestNewStripePattern(t *testing.T) {
	p := object.NewStripePattern(object.White, object.Black)
	assertColorEqual(t, object.White, p.A)
	assertColorEqual(t, object.Black, p.B)
}

func TestStripePattern_At(t *testing.T) {

	type args struct {
		point    ray.Vector
		expected object.RGB
	}

	testCases := []struct {
		name     string
		expected []args
	}{
		{
			name: "A stripe pattern is constant in y",
			expected: []args{
				{
					point:    ray.NewPoint(0, 0, 0),
					expected: object.White,
				},
				{
					point:    ray.NewPoint(0, 1, 0),
					expected: object.White,
				},
				{
					point:    ray.NewPoint(0, 2, 0),
					expected: object.White,
				},
			},
		},
		{
			name: "A stripe pattern is constant in z",
			expected: []args{
				{
					point:    ray.NewPoint(0, 0, 0),
					expected: object.White,
				},
				{
					point:    ray.NewPoint(0, 0, 1),
					expected: object.White,
				},
				{
					point:    ray.NewPoint(0, 0, 2),
					expected: object.White,
				},
			},
		},
		{
			name: "A stripe pattern alternates in x",
			expected: []args{
				{
					point:    ray.NewPoint(0, 0, 0),
					expected: object.White,
				},
				{
					point:    ray.NewPoint(0.9, 0, 0),
					expected: object.White,
				},
				{
					point:    ray.NewPoint(1, 0, 0),
					expected: object.Black,
				},
				{
					point:    ray.NewPoint(-0.1, 0, 0),
					expected: object.Black,
				},
				{
					point:    ray.NewPoint(-1, 0, 0),
					expected: object.Black,
				},
				{
					point:    ray.NewPoint(-1.1, 0, 0),
					expected: object.White,
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			pattern := object.NewStripePattern(object.White, object.Black)
			for _, arg := range tt.expected {
				t.Log("Expecting: ", arg)
				actual := pattern.At(arg.point)
				assertColorEqual(t, arg.expected, actual)
			}
		})
	}
}

func TestStripePattern_AtObj(t *testing.T) {
	testCases := []struct {
		name             string
		point            ray.Vector
		expected         object.RGB
		pattern          object.StripePattern
		objTransform     ray.Matrix
		patternTransform ray.Matrix
	}{
		{
			name:         "Stripes with an object transformation",
			point:        ray.NewPoint(1.5, 0, 0),
			pattern:      object.NewStripePattern(object.White, object.Black),
			objTransform: ray.Scaling(2, 2, 2),
			expected:     object.White,
		},
		{
			name:             "Stripes with a pattern transformation",
			point:            ray.NewPoint(1.5, 0, 0),
			pattern:          object.NewStripePattern(object.White, object.Black),
			patternTransform: ray.Scaling(2, 2, 2),
			expected:         object.White,
		},
		{
			name:             "Stripes with both an object and a pattern transformation",
			point:            ray.NewPoint(2.5, 0, 0),
			pattern:          object.NewStripePattern(object.White, object.Black),
			objTransform:     ray.Scaling(2, 2, 2),
			patternTransform: ray.Translation(0.5, 0, 0),
			expected:         object.White,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			obj := object.NewSphere(ray.NewPoint(0, 0, 0), 1)
			if tt.objTransform != nil {
				obj.SetTransform(tt.objTransform)
			}

			if tt.patternTransform != nil {
				tt.pattern.Transform = tt.patternTransform
			}
			actual := tt.pattern.AtObj(obj, tt.point)
			assertColorEqual(t, tt.expected, actual)
		})
	}
}

func TestStripePattern_IsEmpty(t *testing.T) {

	testCases := []struct {
		name     string
		pattern  object.StripePattern
		expected bool
	}{
		{
			name:     "None empty",
			pattern:  object.NewStripePattern(object.White, object.Black),
			expected: false,
		},
		{
			name:     "Empty",
			pattern:  object.EmptyStripePattern(),
			expected: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.pattern.IsEmpty())
		})
	}
}

func TestNewTestPattern(t *testing.T) {

	testCases := []struct {
		name             string
		point            ray.Vector
		shapeTransform   ray.Matrix
		patternTransform ray.Matrix
		expected         object.RGB
	}{
		{
			name:           "A pattern with an object transformation",
			shapeTransform: ray.Scaling(2, 2, 2),
			point:          ray.NewPoint(2, 3, 4),
			expected:       object.NewColor(1, 1.5, 2),
		},
		{
			name:             "A pattern with a pattern transformation",
			patternTransform: ray.Scaling(2, 2, 2),
			point:            ray.NewPoint(2, 3, 4),
			expected:         object.NewColor(1, 1.5, 2),
		},
		{
			name:             "A pattern with both an object and a pattern transformation",
			patternTransform: ray.Translation(0.5, 1, 1.5),
			shapeTransform:   ray.Scaling(2, 2, 2),
			point:            ray.NewPoint(2.5, 3, 3.5),
			expected:         object.NewColor(0.75, 0.5, 0.25),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			shape := object.NewSphere(ray.ZeroPoint, 1)

			if tt.shapeTransform != nil {
				shape.SetTransform(tt.shapeTransform)
			}

			pattern := object.NewTestPattern()
			if tt.patternTransform != nil {
				pattern.Transform = tt.patternTransform
			}

			actual := pattern.AtObj(shape, tt.point)
			assertColorEqual(t, tt.expected, actual)
		})
	}

}
