package ray

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Dot(t *testing.T) {
	a := NewVec(1, 2, 3)
	b := NewVec(2, 3, 4)
	actual := Dot(a, b)
	assert.Equal(t, float64(20), actual)
}

func Test_Cross(t *testing.T) {
	a := NewVec(1, 2, 3)
	b := NewVec(2, 3, 4)

	assertVec(t, NewVec(-1, 2, -1), Cross(a, b))
	assertVec(t, NewVec(1, -2, 1), Cross(b, a))
}
