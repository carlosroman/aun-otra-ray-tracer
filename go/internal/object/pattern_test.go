package object_test

import (
	"math"
	"testing"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
	"github.com/stretchr/testify/assert"
)

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
		pattern          object.Pattern
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
		pattern  object.Pattern
		expected bool
	}{
		{
			name:     "None empty",
			pattern:  object.NewStripePattern(object.White, object.Black),
			expected: true,
		},
		{
			name:     "Empty",
			pattern:  object.EmptyPattern(),
			expected: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.pattern.IsNotEmpty)
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

func TestNewGradientPattern(t *testing.T) {
	pattern := object.NewGradientPattern(object.White, object.Black)

	testCases := []struct {
		name     string
		point    ray.Vector
		expected object.RGB
	}{
		{
			name:     "0, 0, 0",
			point:    ray.NewPoint(0, 0, 0),
			expected: object.White,
		},
		{
			name:     "0.25, 0, 0",
			point:    ray.NewPoint(0.25, 0, 0),
			expected: object.NewColor(0.75, 0.75, 0.75),
		},
		{
			name:     "0.5, 0, 0",
			point:    ray.NewPoint(0.5, 0, 0),
			expected: object.NewColor(0.5, 0.5, 0.5),
		},
		{
			name:     "0.75, 0, 0",
			point:    ray.NewPoint(0.75, 0, 0),
			expected: object.NewColor(0.25, 0.25, 0.25),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actual := pattern.At(tt.point)
			assertColorEqual(t, tt.expected, actual)
		})
	}
}

func TestNewRingPattern(t *testing.T) {
	pattern := object.NewRingPattern(object.White, object.Black)
	testCases := []struct {
		name     string
		point    ray.Vector
		expected object.RGB
	}{
		{
			name:     "0, 0, 0",
			point:    ray.NewPoint(0, 0, 0),
			expected: object.White,
		},
		{
			name:     "1, 0, 0",
			point:    ray.NewPoint(1, 0, 0),
			expected: object.Black,
		},
		{
			name:     "0, 0, 1",
			point:    ray.NewPoint(0, 0, 1),
			expected: object.Black,
		},
		{
			name:     "√2/2, 0, √2/2",
			point:    ray.NewPoint(math.Sqrt(2)/2, 0, math.Sqrt(2)/2),
			expected: object.Black,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actual := pattern.At(tt.point)
			assertColorEqual(t, tt.expected, actual)
		})
	}
}

func TestNewCheckerPattern(t *testing.T) {
	pattern := object.NewCheckerPattern(object.White, object.Black)

	type args struct {
		point    ray.Vector
		expected object.RGB
	}
	testCases := []struct {
		name     string
		expected []args
	}{
		{
			name: "Checkers should repeat in x",
			expected: []args{
				{
					point:    ray.NewPoint(0, 0, 0),
					expected: object.White,
				},
				{
					point:    ray.NewPoint(0.99, 0, 0),
					expected: object.White,
				},
				{
					point:    ray.NewPoint(1.01, 0, 0),
					expected: object.Black,
				},
			},
		},
		{
			name: "Checkers should repeat in y",
			expected: []args{
				{
					point:    ray.NewPoint(0, 0, 0),
					expected: object.White,
				},
				{
					point:    ray.NewPoint(0, 0.99, 0),
					expected: object.White,
				},
				{
					point:    ray.NewPoint(0, 1.01, 0),
					expected: object.Black,
				},
			},
		},
		{
			name: "Checkers should repeat in z",
			expected: []args{
				{
					point:    ray.NewPoint(0, 0, 0),
					expected: object.White,
				},
				{
					point:    ray.NewPoint(0, 0, 0.99),
					expected: object.White,
				},
				{
					point:    ray.NewPoint(0, 0, 1.01),
					expected: object.Black,
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			for _, arg := range tt.expected {
				actual := pattern.At(arg.point)
				assertColorEqual(t, arg.expected, actual)
			}
		})
	}
}
