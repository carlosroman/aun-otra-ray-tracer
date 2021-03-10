package object_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

func TestNewPlane(t *testing.T) {
	plane := object.NewPlane()
	assert.Equal(t, ray.DefaultIdentityMatrix(), plane.Transform())
	assert.Equal(t, object.DefaultMaterial(), plane.Material())
}

func TestPlane_LocalNormalAt(t *testing.T) {
	plane := object.NewPlane()
	n1 := plane.LocalNormalAt(ray.NewPoint(0, 0, 0), object.NoHit)
	n2 := plane.LocalNormalAt(ray.NewPoint(10, 0, -10), object.NoHit)
	n3 := plane.LocalNormalAt(ray.NewPoint(-5, 0, 150), object.NoHit)

	assertVec(t, ray.NewVec(0, 1, 0), n1)
	assertVec(t, ray.NewVec(0, 1, 0), n2)
	assertVec(t, ray.NewVec(0, 1, 0), n3)
}

func TestPlane_LocalIntersect(t *testing.T) {

	p := object.NewPlane()
	testCases := []struct {
		name     string
		ray      ray.Ray
		actualXS []float64
	}{
		{
			name: "Intersect with a ray parallel to the plane",
			ray:  ray.NewRayAt(ray.NewPoint(0, 10, 0), ray.NewVec(0, 0, 1)),
		},
		{
			name: "Intersect with a coplanar ray",
			ray:  ray.NewRayAt(ray.NewPoint(0, 0, 0), ray.NewVec(0, 0, 1)),
		},
		{
			name:     " A ray intersecting a plane from above",
			ray:      ray.NewRayAt(ray.NewPoint(0, 1, 0), ray.NewVec(0, -1, 0)),
			actualXS: []float64{1},
		},
		{
			name:     "A ray intersecting a plane from below",
			ray:      ray.NewRayAt(ray.NewPoint(0, -1, 0), ray.NewVec(0, 1, 0)),
			actualXS: []float64{1},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actualXS := p.LocalIntersect(tt.ray)
			if len(tt.actualXS) < 1 {
				assert.Empty(t, actualXS)
				return
			}
			assert.Len(t, actualXS, len(tt.actualXS))
		})
	}
}
