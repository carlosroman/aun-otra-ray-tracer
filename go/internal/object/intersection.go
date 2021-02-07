package object

type Intersection struct {
	T   float64 // Intersection value
	Obj Object  // Intersected object
}

type Intersections []Intersection

var NoHit = Intersection{}

func Hit(intersections Intersections) Intersection {
	for i := range intersections {
		if 0 <= intersections[i].T {
			return intersections[i]
		}
	}
	return NoHit
}
