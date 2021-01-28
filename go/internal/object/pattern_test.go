package object_test

import (
	"testing"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
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
