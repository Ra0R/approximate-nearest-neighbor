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

type GraphConfig struct {
	NextVertexId         uint64
	Path                 string
	DistanceFunctionName string
}

type GraphFactory struct {
}

func (gf *GraphFactory) New(path string, distanceFunctionName string) (GraphInterface, error) {

	// Create file if it not exists
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	distanceFunction := getDistanceFunctionByName(distanceFunctionName)

	if distanceFunction == nil {
		panic("no distance function found " + distanceFunctionName)
	}

	return &Graph{
		Path:                 path,
		nextVertexId:         0,
		vertices:             make(map[uint64]*Vertex), // maps vertex id to the actual vertex
		edges:                make(map[uint64][]*Edge), // maps vertex id to its edges
		distanceFunctionName: distanceFunctionName,
		distanceFunction:     distanceFunction,
	}, nil
}

func (gf *GraphFactory) Open(path string) (GraphInterface, error) {
	verticeFile, err := os.OpenFile(path+"vertices.ann", os.O_RDWR, 0644)

	defer verticeFile.Close()
	if err != nil {
		return nil, err
	}

	edgeFile, err := os.OpenFile(path+"edges.ann", os.O_RDWR, 0644)
	defer edgeFile.Close()
	if err != nil {
		return nil, err
	}

	configFile, err := os.OpenFile(path+"graph.config", os.O_RDWR, 0644)
	defer configFile.Close()
	if err != nil {
		return nil, err
	}

	var graph Graph

	decoder := gob.NewDecoder(edgeFile)
	err = decoder.Decode(&graph.edges)

	decoder = gob.NewDecoder(verticeFile)
	err = decoder.Decode(&graph.vertices)

	var graphConfig GraphConfig
	decoder = gob.NewDecoder(configFile)
	err = decoder.Decode(&graphConfig)

	graph.distanceFunction = getDistanceFunctionByName(graphConfig.DistanceFunctionName)
	graph.nextVertexId = graphConfig.NextVertexId
	graph.Path = graphConfig.Path

	return &graph, err
}

func (gf *GraphFactory) Delete(path string) error {
	err := os.Remove(path + "vertices.ann")
	if err != nil {
		return err
	}
	err = os.Remove(path + "edges.ann")
	if err != nil {
		return err
	}

	err = os.Remove(path + "graph.config")

	if err != nil {
		return err
	}
	return err
}
