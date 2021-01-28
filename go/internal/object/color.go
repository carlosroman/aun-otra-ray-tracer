package object

type RGB struct {
	R, G, B float64
}

var (
	Black = NewColor(0, 0, 0)
	White = NewColor(1, 1, 1)
)

func NewColor(red, green, blue float64) RGB {
	return RGB{
		R: red,
		G: green,
		B: blue,
	}
}

func (r RGB) Add(c RGB) RGB {
	return NewColor(
		r.R+c.R,
		r.G+c.G,
		r.B+c.B,
	)
}

func (r RGB) Subtract(c RGB) RGB {
	return NewColor(
		r.R-c.R,
		r.G-c.G,
		r.B-c.B,
	)
}

// Technically Hadamard product or Schur product
func (r RGB) Multiply(c RGB) RGB {
	return NewColor(
		r.R*c.R,
		r.G*c.G,
		r.B*c.B,
	)
}

func (r RGB) MultiplyBy(scalar float64) RGB {
	return NewColor(
		r.R*scalar,
		r.G*scalar,
		r.B*scalar,
	)
}

func (r RGB) Plus(c RGB) RGB {
	return NewColor(
		r.R+c.R,
		r.G+c.G,
		r.B+c.B,
	)
}
