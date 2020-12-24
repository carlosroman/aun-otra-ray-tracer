package ray

import (
	"errors"
)

const (
	defaultIdentityMatrixSize = 4
)

var (
	Zero             = NewVec(0, 0, 0)
	NonInvertibleErr = errors.New("non-invertible")
)
