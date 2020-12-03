package ray_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

func Test_Dot(t *testing.T) {
	a := ray.NewVec(1, 2, 3)
	b := ray.NewVec(2, 3, 4)
	actual := ray.Dot(a, b)
	assert.Equal(t, float64(20), actual)
}

func Test_Cross(t *testing.T) {
	a := ray.NewVec(1, 2, 3)
	b := ray.NewVec(2, 3, 4)

	assertVec(t, ray.NewVec(-1, 2, -1), ray.Cross(a, b))
	assertVec(t, ray.NewVec(1, -2, 1), ray.Cross(b, a))
}