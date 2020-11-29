package ray

type Vector interface {
	GetX() float64
	GetY() float64
	GetZ() float64
}

func Dot(a, b Vector) float64 {
	return (a.GetX() * b.GetX()) +
		(a.GetY() * b.GetY()) +
		(a.GetZ() * b.GetZ())
}

func Cross(a, b Vector) Vector {
	return NewVec(
		(a.GetY()*b.GetZ())-(a.GetZ()*b.GetY()),
		(a.GetZ()*b.GetX())-(a.GetX()*b.GetZ()),
		(a.GetX()*b.GetY())-(a.GetY()*b.GetX()))
}
