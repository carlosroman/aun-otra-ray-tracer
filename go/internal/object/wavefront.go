package object

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/carlosroman/aun-otra-ray-tracer/go/internal/ray"
)

type WavefrontObj struct {
	Vertices,
	Normals []ray.Vector
	Group Group
}

func (w *WavefrontObj) Object() Object {
	return &w.Group
}

var pointsReg = regexp.MustCompile("(?P<Point>\\d+)(/(?P<Texture>\\d{*})/(?P<Normal>\\d+))?")

func NewWavefrontObj(reader io.Reader, opts ...Option) (wv WavefrontObj, err error) {
	var vertices, normals []ray.Vector
	defaultGroup := NewGroup(opts...)
	scanner := bufio.NewScanner(reader)
	re := regexp.MustCompile("\\s+")
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "v ") {
			point, err := parsePoint(re, text)
			if err != nil {
				return wv, err
			}
			vertices = append(vertices, point)
		}
		if strings.HasPrefix(text, "vn ") {
			point, err := parsePoint(re, text)
			if err != nil {
				return wv, err
			}
			normals = append(normals, point)
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
				groupsToAdd, err = createTriangle(points, vertices, normals)
			case 5:
				groupsToAdd, err = createTrianglesFromSquare(points, vertices, normals)
			case 6:
				groupsToAdd, err = createTrianglesFromPentagon(points, vertices, normals)
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
	wv.Normals = normals
	wv.Group = defaultGroup
	return wv, nil
}

func parsePoint(re *regexp.Regexp, text string) (ray.Vector, error) {
	points := re.Split(text, 4)
	x, err := strconv.ParseFloat(points[1], 64)
	if err != nil {
		return nil, err
	}
	y, err := strconv.ParseFloat(points[2], 64)
	if err != nil {
		return nil, err
	}
	z, err := strconv.ParseFloat(points[3], 64)
	if err != nil {
		return nil, err
	}
	point := ray.NewPoint(x, y, z)
	return point, nil
}

func createTrianglesFromSquare(points []string, vertices, normals []ray.Vector) ([]Object, error) {
	objs := make([]Object, 2)
	p1, _, p1n, err := parseFaceElement(points[1])
	if err != nil {
		return nil, err
	}

	for i := 2; i < len(points)-1; i++ {
		p2, _, p2n, err := parseFaceElement(points[i])
		if err != nil {
			return nil, err
		}
		p3, _, p3n, err := parseFaceElement(points[i+1])
		if err != nil {
			return nil, err
		}

		if len(normals) > 0 && p1n >= 0 && p2n >= 0 && p3n >= 0 {
			objs[i-2] = NewSmoothTriangle(
				vertices[p1-1],
				vertices[p2-1],
				vertices[p3-1],
				normals[p1n-1],
				normals[p2n-1],
				normals[p3n-1],
			)
		} else {
			objs[i-2] = NewTriangle(
				vertices[p1-1],
				vertices[p2-1],
				vertices[p3-1],
			)
		}
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

func createTrianglesFromPentagon(points []string, vertices, normals []ray.Vector) ([]Object, error) {
	objs := make([]Object, 3)
	p1, _, p1n, err := parseFaceElement(points[1])
	if err != nil {
		return nil, err
	}

	for i := 2; i < len(points)-1; i++ {
		p2, _, p2n, err := parseFaceElement(points[i])
		if err != nil {
			return nil, err
		}
		p3, _, p3n, err := parseFaceElement(points[i+1])
		if err != nil {
			return nil, err
		}

		if len(normals) > 0 && p1n >= 0 && p2n >= 0 && p3n >= 0 {
			objs[i-2] = NewSmoothTriangle(
				vertices[p1-1],
				vertices[p2-1],
				vertices[p3-1],
				normals[p1n-1],
				normals[p2n-1],
				normals[p3n-1],
			)
		} else {
			objs[i-2] = NewTriangle(
				vertices[p1-1],
				vertices[p2-1],
				vertices[p3-1],
			)
		}
	}
	return objs, nil
}

func createTriangle(points []string, vertices, normals []ray.Vector) ([]Object, error) {
	p1, _, p1n, err := parseFaceElement(points[1])
	if err != nil {
		return nil, err
	}
	p2, _, p2n, err := parseFaceElement(points[2])
	if err != nil {
		return nil, err
	}
	p3, _, p3n, err := parseFaceElement(points[3])
	if err != nil {
		return nil, err
	}

	if len(normals) > 0 && p1n >= 0 && p2n >= 0 && p3n >= 0 {
		return []Object{NewSmoothTriangle(
			vertices[p1-1],
			vertices[p2-1],
			vertices[p3-1],
			normals[p1n-1],
			normals[p2n-1],
			normals[p3n-1],
		)}, nil
	}

	return []Object{NewTriangle(
		vertices[p1-1],
		vertices[p2-1],
		vertices[p3-1],
	)}, nil
}
