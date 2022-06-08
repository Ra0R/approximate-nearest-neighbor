package ann

import (
	"encoding/gob"
	"errors"
	"os"
)

var (
	// ErrInvalidPath is returned when the path that has been given is not valid (inexistent/not writable)
	ErrInvalidPath = errors.New("supplied argument 'path' is not valid")
)

type GraphFactory struct {
}

func (gf *GraphFactory) New(path string, distanceFunctionName string) (GraphInterface, error) {

	// Create file if it not exists
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	var distanceFunction func(*Vertex, *Vertex) float64
	if distanceFunctionName == "euclidean" {
		distanceFunction = euclideanDistance
	}

	return &Graph{
		Path:             path,
		nextVertexId:     0,
		vertices:         make(map[uint64]*Vertex), // maps vertex id to the actual vertex
		edges:            make(map[uint64][]*Edge), // maps vertex id to its edges
		distanceFunction: distanceFunction,
	}, nil
}

func (gf *GraphFactory) Open(path string) (GraphInterface, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0660)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	var graph Graph
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&graph)
	return &graph, err
}

func (gf *GraphFactory) Delete(path string) error {
	err := os.Remove(path)
	return err
}
