package ann

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGraph(t *testing.T) {
	assert := assert.New(t)
	factory := GraphFactory{}
	path := "."
	graph, err := factory.New(path)

	assert.NoError(err)
	assert.NotNil(graph)
}

func TestOpen(t *testing.T) {
	assert := assert.New(t)

	factory := GraphFactory{}
	path := "."
	graph, err := factory.Open(path)

	assert.NoError(err)
	assert.NotNil(graph)
	graph.Close()
}

func TestClose(t *testing.T) {
	assert := assert.New(t)

	factory := GraphFactory{}
	path := "."
	graph, err := factory.New(path)

	assert.NoError(err)
	assert.NotNil(graph)
	err = graph.Close()
	assert.NoError(err)
}
func TestOpen_NoPath_Fails(t *testing.T) {
	assert := assert.New(t)

	const invalidPath = ""

	factory := GraphFactory{}
	graph, err := factory.Open(invalidPath)

	assert.Nil(graph, "Creation should fail")
	assert.EqualError(err, ErrInvalidPath.Error())
}

func TestInsertionOnEmptyGraph(t *testing.T) {
	assert := assert.New(t)
	factory := GraphFactory{}
	path := "."
	graph, err := factory.New(path)
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

	err = graph.Close()
	assert.NoError(err)
}

func TestGetNearestNeighbor(t *testing.T) {
	assert := assert.New(t)
	factory := GraphFactory{}
	path := "."
	graph, err := factory.New(path)
	assert.NoError(err)
	assert.NotNil(graph)

	// Insert Point1 = (0.0)
	var dimension uint16 = 1
	coordinates := make([]float64, 1)
	coordinates[0] = 0.0
	point1, err := NewPoint(dimension, coordinates)
	assert.NoError(err)
	assert.NotNil(point1)
	err = graph.NNInsert(point1, 3, 1)
	assert.NoError(err)

	// Insert Point2 = (1.0)
	point2, err := NewPoint(dimension, []float64{1.0})
	assert.NoError(err)
	err = graph.NNInsert(point2, 3, 1)

	// Insert Point3 = (2.0)
	point3, err := NewPoint(dimension, []float64{2.0})
	assert.NoError(err)
	err = graph.NNInsert(point3, 3, 1)

	// Search nearest point of Point1, should return Point2 (assuming Euclidean metric)
	nearestNeighbors, err := graph.NNSearch(point1, 3, 1)
	assert.NoError(err)
	assert.NotNil(nearestNeighbors)
	assert.Equal(nearestNeighbors[0], point2)

	err = graph.Close()
	assert.NoError(err)
}

func TestGetNearestNeighbors(t *testing.T) {
	assert := assert.New(t)
	factory := GraphFactory{}
	path := "."
	graph, err := factory.New(path)
	assert.NoError(err)
	assert.NotNil(graph)

	// Insert Point1 = (0.0)
	var dimension uint16 = 1
	coordinates := make([]float64, 1)
	coordinates[0] = 0.0
	point1, err := NewPoint(dimension, coordinates)
	assert.NoError(err)
	assert.NotNil(point1)
	err = graph.NNInsert(point1, 3, 1)
	assert.NoError(err)

	// Insert Point2 = (1.0)
	point2, err := NewPoint(dimension, []float64{1.0})
	assert.NoError(err)
	err = graph.NNInsert(point2, 3, 1)

	// Insert Point3 = (2.0)
	point3, err := NewPoint(dimension, []float64{2.0})
	assert.NoError(err)
	err = graph.NNInsert(point3, 3, 1)

	// Search nearest neighbors of point1
	nearestNeighbors, err := graph.NNSearch(point1, 3, 2)
	assert.NoError(err)
	assert.NotNil(nearestNeighbors)
	// Assuming that they are ordered by distance (ascending)
	assert.Equal(nearestNeighbors[0], point2)
	assert.Equal(nearestNeighbors[1], point3)

	err = graph.Close()
	assert.NoError(err)
}
