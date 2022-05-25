package ann

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	DEFAULT_PATH = ".\\testgraph.ann"
)

func TestNewGraph(t *testing.T) {
	assert := assert.New(t)
	factory := GraphFactory{}
	graph, err := factory.New(DEFAULT_PATH)

	assert.NoError(err)
	assert.NotNil(graph)
}

func TestClose(t *testing.T) {
	assert := assert.New(t)

	factory := GraphFactory{}
	graph, err := factory.New(DEFAULT_PATH)

	assert.NoError(err)
	assert.NotNil(graph)
	err = graph.Close()
	assert.NoError(err)
	err = factory.Delete(DEFAULT_PATH)
	assert.NoError(err)
}

func TestOpen_NoPath_Fails(t *testing.T) {
	assert := assert.New(t)

	const invalidPath = ""

	factory := GraphFactory{}
	graph, err := factory.Open(invalidPath)

	assert.Nil(graph, "Creation should fail")
	assert.Error(err)
}

func TestSaveGraphReOpen(t *testing.T) {
	assert := assert.New(t)

	factory := GraphFactory{}
	graph, err := factory.New(DEFAULT_PATH)

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
	assert.NoError(err)

	graph.Close()

	// Reopen graph from disk
	graph2, err := factory.Open(DEFAULT_PATH)

	// TODO Manual encoding
	graph2.String()

	err = factory.Delete(DEFAULT_PATH)
	assert.NoError(err)
}

func TestInsertionOnEmptyGraph(t *testing.T) {
	assert := assert.New(t)
	factory := GraphFactory{}

	graph, err := factory.New(DEFAULT_PATH)
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

	graph.String()

	err = graph.Close()
	assert.NoError(err)

	err = factory.Delete(DEFAULT_PATH)
	assert.NoError(err)
}

func TestGetNearestNeighbor(t *testing.T) {
	assert := assert.New(t)
	factory := GraphFactory{}
	graph, err := factory.New(DEFAULT_PATH)
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
	assert.NoError(err)

	// Insert Point3 = (2.0)
	point3, err := NewPoint(dimension, []float64{2.0})
	assert.NoError(err)
	err = graph.NNInsert(point3, 3, 1)
	assert.NoError(err)

	// Search nearest point of Point1, should return Point2 (assuming Euclidean metric)
	nearestNeighbors, err := graph.NNSearch(point1, 1, 3)
	assert.NoError(err)
	assert.NotNil(nearestNeighbors)
	assert.Equal(0.0, point2.calculateDistance(nearestNeighbors[1]))

	graph.String()

	err = graph.Close()
	assert.NoError(err)

	err = factory.Delete(DEFAULT_PATH)
	assert.NoError(err)
}

func TestGetNearestNeighbors(t *testing.T) {
	assert := assert.New(t)
	factory := GraphFactory{}

	graph, err := factory.New(DEFAULT_PATH)
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

	// Insert Point4 = (3.0)
	point4, err := NewPoint(dimension, []float64{3.0})
	assert.NoError(err)
	err = graph.NNInsert(point4, 3, 1)

	// Insert Point5 = (4.0)
	point5, err := NewPoint(dimension, []float64{3.0})
	assert.NoError(err)
	err = graph.NNInsert(point5, 3, 1)

	// Insert Point6 = (5.0)
	point6, err := NewPoint(dimension, []float64{3.0})
	assert.NoError(err)
	err = graph.NNInsert(point6, 3, 1)

	// Search nearest neighbors of point1
	nearestNeighbors, err := graph.NNSearch(point1, 5, 3)
	assert.NoError(err)
	assert.NotNil(nearestNeighbors)

	// Assuming that they are ordered by distance (ascending)
	graph.String()
	assert.Equal(point1.calculateDistance(nearestNeighbors[0]), 0.0)
	assert.Equal(point2.calculateDistance(nearestNeighbors[1]), 0.0)
	assert.Equal(point3.calculateDistance(nearestNeighbors[2]), 0.0)

	err = graph.Close()
	assert.NoError(err)

	err = factory.Delete(DEFAULT_PATH)
	assert.NoError(err)
}

func TestNNInsert(t *testing.T) {
	assert := assert.New(t)
	factory := GraphFactory{}

	graph, err := factory.New(DEFAULT_PATH)
	assert.NoError(err)
	assert.NotNil(graph)

	// Insert Point1 = (0.0)
	var dimension uint16 = 1
	coordinates := make([]float64, 1)
	coordinates[0] = 0.0
	point1, err := NewPoint(dimension, coordinates)
	assert.NoError(err)

	err = graph.NNInsert(point1, 3, 1)
	assert.NoError(err)
	graph.String()
}
