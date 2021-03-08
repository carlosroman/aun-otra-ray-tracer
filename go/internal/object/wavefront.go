package object

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

type WavefrontObj struct {
	Vertices []ray.Vector
	Group    Group
}

func (w *WavefrontObj) Object() Object {
	return &w.Group
}

var pointsReg = regexp.MustCompile("(?P<Point>\\d+)(\\/(?P<Texture>\\d{*})\\/(?P<Normal>\\d+))?")

func NewWavefrontObj(reader io.Reader) (wv WavefrontObj, err error) {
	var vertices []ray.Vector
	defaultGroup := NewGroup()
	scanner := bufio.NewScanner(reader)
	re := regexp.MustCompile("\\s+")
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "v ") {
			points := re.Split(text, 4)
			x, err := strconv.ParseFloat(points[1], 64)
			if err != nil {
				return wv, err
			}
			y, err := strconv.ParseFloat(points[2], 64)
			if err != nil {
				return wv, err
			}
			z, err := strconv.ParseFloat(points[3], 64)
			if err != nil {
				return wv, err
			}
			vertices = append(vertices, ray.NewPoint(x, y, z))
		}
		if strings.HasPrefix(text, "g ") {
			if vertices == nil {
				return wv, errors.New("no vertices found")
			}
			group := NewGroup()
			defaultGroup.AddChild(&group)
		}
		if strings.HasPrefix(text, "f ") {
			if vertices == nil {
				return wv, errors.New("no vertices found")
			}
			points := re.Split(strings.TrimSpace(text), -1)
			var groupsToAdd []Object
			switch len(points) {
			case 4:
				groupsToAdd, err = createTriangle(points, vertices)
			case 5:
				groupsToAdd, err = createTrianglesFromSquare(points, vertices)
			case 6:
				groupsToAdd, err = createTrianglesFromPentagon(points, vertices)
			default:
				return WavefrontObj{}, errors.New(fmt.Sprintf("do not know how to handle polygon with %v points", len(points)-1))
			}
			if err != nil {
				return wv, err
			}
			if defaultGroup.Children != nil {
				if g, ok := defaultGroup.Children[len(defaultGroup.Children)-1].(*Group); ok {
					g.AddChild(groupsToAdd...)
				} else {
					defaultGroup.AddChild(groupsToAdd...)
				}
			} else {
				defaultGroup.AddChild(groupsToAdd...)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return wv, err
	}
	wv.Vertices = vertices
	wv.Group = defaultGroup
	return wv, nil
}

func createTrianglesFromSquare(points []string, vertices []ray.Vector) ([]Object, error) {
	objs := make([]Object, 2)
	p1, _, _, err := parseFaceElement(points[1])
	if err != nil {
		return nil, err
	}

	for i := 2; i < len(points)-1; i++ {
		p2, _, _, err := parseFaceElement(points[i])
		if err != nil {
			return nil, err
		}
		p3, _, _, err := parseFaceElement(points[i+1])
		if err != nil {
			return nil, err
		}
		objs[i-2] = NewTriangle(
			vertices[p1-1],
			vertices[p2-1],
			vertices[p3-1],
		)
	}
	return objs, nil
}

func parseFaceElement(face string) (v, vt, vn int, err error) {
	sub := pointsReg.FindAllStringSubmatch(face, 3)
	v, err = strconv.Atoi(sub[0][0])
	if len(sub) == 1 {
		vt = -1
		vn = v
		return
	}
	if len(sub) == 2 {
		vt = -1
		vn, err = strconv.Atoi(sub[1][0])
		return
	}
	vt, err = strconv.Atoi(sub[1][0])
	vn, err = strconv.Atoi(sub[2][0])
	return
}

func createTrianglesFromPentagon(points []string, vertices []ray.Vector) ([]Object, error) {
	objs := make([]Object, 3)
	p1, _, _, err := parseFaceElement(points[1])
	if err != nil {
		return nil, err
	}

	for i := 2; i < len(points)-1; i++ {
		p2, _, _, err := parseFaceElement(points[i])
		if err != nil {
			return nil, err
		}
		p3, _, _, err := parseFaceElement(points[i+1])
		if err != nil {
			return nil, err
		}
		objs[i-2] = NewTriangle(
			vertices[p1-1],
			vertices[p2-1],
			vertices[p3-1],
		)
	}
	return objs, nil
}

func createTriangle(points []string, vertices []ray.Vector) ([]Object, error) {
	p1, _, _, err := parseFaceElement(points[1])
	if err != nil {
		return nil, err
	}
	p2, _, _, err := parseFaceElement(points[2])
	if err != nil {
		return nil, err
	}
	p3, _, _, err := parseFaceElement(points[3])
	if err != nil {
		return nil, err
	}

	return []Object{NewTriangle(
		vertices[p1-1],
		vertices[p2-1],
		vertices[p3-1],
	)}, nil
}
