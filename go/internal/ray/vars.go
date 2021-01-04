package ray

import (
	"errors"
)

const (
	defaultIdentityMatrixSize = 4
)

var (
	ZeroVector       = NewVec(0, 0, 0)
	ZeroPoint        = NewPoint(0, 0, 0)
	NonInvertibleErr = errors.New("non-invertible")
)

type Axis int

const (
	X Axis = iota
	Y
	Z
)
