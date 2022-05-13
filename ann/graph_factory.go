package ann

import "errors"

var (
	// ErrInvalidPath is returned when the path that has been given is not valid (inexistent/not writable)
	ErrInvalidPath = errors.New("Supplied argument 'path' is not valid")
)

type GraphFactory struct {
}

func (gf *GraphFactory) New(path string) (GraphInterface, error) {
	// TODO create file in path

	return &Graph{
		nextVertexId: 0,
		vertices:     make(map[uint64]*Vertex), // maps vertex id to the actual vertex
		edges:        make(map[uint64][]*Edge), // maps vertex id to its edges
	}, nil
}

func (gf *GraphFactory) Open(path string) (GraphInterface, error) {
	return nil, errors.New("not implemented")
}
