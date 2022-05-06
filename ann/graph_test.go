package ann

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGraph(t *testing.T) {
	assert := assert.New(t)
	factory := GraphFactory{}
	graph, err := factory.New()

	assert.NoError(err)
	assert.NotNil(graph)
}

func TestInsertionOnEmptyGraph(t *testing.T) {
	assert := assert.New(t)
	factory := GraphFactory{}
	graph, err := factory.New()
	assert.NoError(err)
	assert.NotNil(graph)

	var dimension uint16 = 1
	coordinates := make([]float64, 1)
	coordinates[0] = 0.0
	point, err := NewPoint(dimension, coordinates)
	assert.NoError(err)
	assert.NotNil(point)

	err = graph.NNInsert(point, 3, 1)
	assert.NoError(err)
}
