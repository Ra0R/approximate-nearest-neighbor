package ann

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newVertex(dimension int, coordinates []float64) Vertex {
	testPoint, _ := NewPoint(uint16(dimension), coordinates)
	var testObject ObjectInterface = testPoint

	vertex := Vertex{
		id:     0,
		object: &testObject,
	}

	return vertex
}

func TestLarger(t *testing.T) {
	assert := assert.New(t)
	vertex := newVertex(1, []float64{0})
	vertexComparator := &VertexComparator{vertex: &vertex}

	vertex1 := newVertex(1, []float64{1})
	vertex2 := newVertex(1, []float64{2})

	// Distance 1 < Distance 2 --> -1
	assert.Equal(-1, vertexComparator.compare(&vertex1, &vertex2))

	// Distance 1 > Distance 2 --> -1
	assert.Equal(1, vertexComparator.compare(&vertex2, &vertex1))

	assert.Equal(0, vertexComparator.compare(&vertex1, &vertex1))
}
