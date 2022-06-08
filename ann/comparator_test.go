package ann

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLarger(t *testing.T) {
	assert := assert.New(t)
	vertex := Vertex{1, []float64{0}}
	vertexComparator := &VertexComparator{vertex: &vertex}

	vertex1 := Vertex{1, []float64{1}}
	vertex2 := Vertex{1, []float64{2}}

	// Distance 1 < Distance 2 --> -1
	assert.Equal(-1, vertexComparator.compare(&vertex1, &vertex2))

	// Distance 1 > Distance 2 --> -1
	assert.Equal(1, vertexComparator.compare(&vertex2, &vertex1))

	assert.Equal(0, vertexComparator.compare(&vertex1, &vertex1))
}
