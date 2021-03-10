package object_test

import (
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carlosroman/aun-otra-ray-trace/go/internal/object"
	"github.com/carlosroman/aun-otra-ray-trace/go/internal/ray"
)

var (
	utahTeapot         = path.Join("..", "..", "..", "test", "models", "utah-teapot.obj")
	utahTeapotVerts    = 3200
	utahTeapotLow      = path.Join("..", "..", "..", "test", "models", "utah-teapot-low.obj")
	utahTeapotLowVerts = 128
)

func TestNewWavefrontObj(t *testing.T) {

	testCases := []struct {
		name        string
		fileContent string
		assertFunc  func(t *testing.T, obj object.WavefrontObj)
	}{
		{
			name: "Ignoring unrecognized lines",
			assertFunc: func(t *testing.T, obj object.WavefrontObj) {
				assert.Empty(t, obj.Vertices)
			},
			fileContent: `There was a young lady named Bright
who traveled much faster than light.
She set out one day
in a relative way,
and came back the previous nigh
`,
		},
		{
			name: "Vertex records",
			assertFunc: func(t *testing.T, obj object.WavefrontObj) {
				require.Len(t, obj.Vertices, 4)

				assertVec(t, obj.Vertices[0], ray.NewPoint(-1, 1, 0))
				assertVec(t, obj.Vertices[1], ray.NewPoint(-1, 0.5, 0))
				assertVec(t, obj.Vertices[2], ray.NewPoint(1, 0, 0))
				assertVec(t, obj.Vertices[3], ray.NewPoint(1, 1, 0))
			},
			fileContent: `v -1 1 0
v -1.0000 0.5000 0.0000
v 1 0 0
v 1 1 0
`,
		},
		{
			name: "Parsing triangle faces",
			assertFunc: func(t *testing.T, obj object.WavefrontObj) {
				require.Len(t, obj.Group.Children, 2)

				t1 := object.NewTriangle(obj.Vertices[0], obj.Vertices[1], obj.Vertices[2])
				t1.SetParent(&obj.Group)
				assert.Equal(t, t1, obj.Group.Children[0])
				t2 := object.NewTriangle(obj.Vertices[0], obj.Vertices[2], obj.Vertices[3])
				t2.SetParent(&obj.Group)
				assert.Equal(t, t2, obj.Group.Children[1])
			},
			fileContent: `v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0

f 1 2 3
f 1 3 4
`,
		},
		{
			name: "Parsing triangle faces with texture coordinate indices and normal indices",
			assertFunc: func(t *testing.T, obj object.WavefrontObj) {
				require.Len(t, obj.Group.Children, 2)

				t1 := object.NewTriangle(obj.Vertices[0], obj.Vertices[1], obj.Vertices[2])
				t1.SetParent(&obj.Group)
				assert.Equal(t, t1, obj.Group.Children[0])
				t2 := object.NewTriangle(obj.Vertices[0], obj.Vertices[2], obj.Vertices[3])
				t2.SetParent(&obj.Group)
				assert.Equal(t, t2, obj.Group.Children[1])
			},
			fileContent: `v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0

f 1/1/1 2/2/2 3/3/3
f 1//1 3//3 4//4
`,
		},
		{
			name: "Triangulating square",
			assertFunc: func(t *testing.T, obj object.WavefrontObj) {
				require.Len(t, obj.Group.Children, 2)

				t1 := object.NewTriangle(obj.Vertices[1-1], obj.Vertices[2-1], obj.Vertices[3-1])
				t1.SetParent(&obj.Group)
				assert.Equal(t, t1, obj.Group.Children[0], "t1")
				t2 := object.NewTriangle(obj.Vertices[1-1], obj.Vertices[3-1], obj.Vertices[4-1])
				t2.SetParent(&obj.Group)
				assert.Equal(t, t2, obj.Group.Children[1], "t2")
			},
			fileContent: `v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0
v 0 2 0

f 1 2 3 4
`,
		},
		{
			name: "Triangulating polygons",
			assertFunc: func(t *testing.T, obj object.WavefrontObj) {
				require.Len(t, obj.Group.Children, 3)

				t1 := object.NewTriangle(obj.Vertices[1-1], obj.Vertices[2-1], obj.Vertices[3-1])
				t1.SetParent(&obj.Group)
				assert.Equal(t, t1, obj.Group.Children[0])
				t2 := object.NewTriangle(obj.Vertices[1-1], obj.Vertices[3-1], obj.Vertices[4-1])
				t2.SetParent(&obj.Group)
				assert.Equal(t, t2, obj.Group.Children[1])
				t3 := object.NewTriangle(obj.Vertices[1-1], obj.Vertices[4-1], obj.Vertices[5-1])
				t3.SetParent(&obj.Group)
				assert.Equal(t, t3, obj.Group.Children[2])
			},
			fileContent: `v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0
v 0 2 0

f 1 2 3 4 5
`,
		},
		{
			name: "Faces with normals",
			assertFunc: func(t *testing.T, obj object.WavefrontObj) {
				require.Len(t, obj.Group.Children, 2)
				require.Len(t, obj.Normals, 3)

				t1 := object.NewSmoothTriangle(obj.Vertices[1-1], obj.Vertices[2-1], obj.Vertices[3-1],
					obj.Normals[3-1], obj.Normals[1-1], obj.Normals[2-1],
				)
				t1.SetParent(&obj.Group)
				assert.Equal(t, t1, obj.Group.Children[0])
				t2 := object.NewSmoothTriangle(obj.Vertices[3-1], obj.Vertices[2-1], obj.Vertices[1-1],
					obj.Normals[2-1], obj.Normals[1-1], obj.Normals[3-1],
				)
				t2.SetParent(&obj.Group)
				assert.Equal(t, t2, obj.Group.Children[1])

			},
			fileContent: `v 0 1 0
v -1 0 0
v 1 0 0

vn -1 0 0
vn 1 0 0
vn 0 1 0

f 1//3 2//1 3//2
f 3/0/2 2/102/1 1/14/3
`,
		},
		{
			name: "Triangles in groups",
			assertFunc: func(t *testing.T, obj object.WavefrontObj) {
				require.Len(t, obj.Group.Children, 2)

				t1 := object.NewTriangle(obj.Vertices[1-1], obj.Vertices[2-1], obj.Vertices[3-1])
				g1, ok := obj.Group.Children[0].(*object.Group)
				require.True(t, ok, "Can not assign G1 as object.Group")
				t1.SetParent(g1) // Setting p
				require.Len(t, g1.Children, 1)
				assert.Equal(t, t1, g1.Children[0], "G1 child")

				t2 := object.NewTriangle(obj.Vertices[1-1], obj.Vertices[3-1], obj.Vertices[4-1])
				g2, ok := obj.Group.Children[1].(*object.Group)
				require.True(t, ok, "Can not assign G2 as object.Group")
				t2.SetParent(g2)
				require.Len(t, g2.Children, 1)
				assert.Equal(t, t2, g2.Children[0], "G2 child")
			},
			fileContent: `v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0

g FirstGroup
f 1 2 3
g SecondGroup
f 1 3 4
`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			bufferString := bytes.NewBufferString(tt.fileContent)
			obj, err := object.NewWavefrontObj(bufferString)
			require.NoError(t, err)
			tt.assertFunc(t, obj)
		})
	}
}

func TestWavefrontObj_Object(t *testing.T) {

	content := `v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0

f 1 2 3`

	bufferString := bytes.NewBufferString(content)
	obj, err := object.NewWavefrontObj(bufferString)
	require.NoError(t, err)
	actual := obj.Object()

	t1 := object.NewTriangle(obj.Vertices[1-1], obj.Vertices[2-1], obj.Vertices[3-1])
	g1, ok := actual.(*object.Group)
	require.True(t, ok, "Can not assign G1 as object.Group")
	t1.SetParent(g1) // Setting p
	require.Len(t, g1.Children, 1)
	assert.Equal(t, t1, g1.Children[0], "G1 child")
}

var benchObj object.Object

func BenchmarkNewWavefrontObj(b *testing.B) {

	testCases := []struct {
		name string
		path string
	}{
		{
			name: "utah teapot",
			path: utahTeapot,
		},
		{
			name: "utah teapot low res",
			path: utahTeapotLow,
		},
	}
	for _, tt := range testCases {
		b.Run(tt.name, func(b *testing.B) {
			var outObj object.Object
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				f, err := os.Open(tt.path)
				require.NoError(b, err)
				obj, err := object.NewWavefrontObj(f)
				require.NoError(b, err)
				outObj = obj.Object()
			}
			benchObj = outObj
			assert.NotNil(b, benchObj)
		})
	}
}
