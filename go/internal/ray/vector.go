package ray

type Point Vector

type Vector interface {
	GetX() float64
	SetX(x float64) Vector
	GetY() float64
	SetY(y float64) Vector
	GetZ() float64
	SetZ(z float64) Vector
	GetW() float64
	SetW(w float64) Vector
	Multiply(by float64) Vector
	Multiplied(by Matrix) Vector
	Add(vec Vector) Vector
	Divide(by float64) Vector
	Subtract(vec Vector) Vector
	Magnitude() float64
	Normalize() Vector
	Negate() Vector
	Reflect(normal Vector) Vector
	Dot(n Vector) float64
}

func Dot(a, b Vector) float64 {
	return (a.GetX() * b.GetX()) +
		(a.GetY() * b.GetY()) +
		(a.GetZ() * b.GetZ()) +
		(a.GetW() * b.GetW())
}

func Cross(a, b Vector) Vector {
	return NewVec(
		(a.GetY()*b.GetZ())-(a.GetZ()*b.GetY()),
		(a.GetZ()*b.GetX())-(a.GetX()*b.GetZ()),
		(a.GetX()*b.GetY())-(a.GetY()*b.GetX()))
}

func Translation(x, y, z float64) Matrix {
	t := IdentityMatrix(defaultIdentityMatrixSize, defaultIdentityMatrixSize)
	t[0][3] = x
	t[1][3] = y
	t[2][3] = z
	return t
}

func Scaling(x, y, z float64) Matrix {
	t := IdentityMatrix(defaultIdentityMatrixSize, defaultIdentityMatrixSize)
	t[0][0] = x
	t[1][1] = y
	t[2][2] = z
	return t
}
