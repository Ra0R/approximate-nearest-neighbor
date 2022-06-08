package ann

type GraphFactoryInterface interface {

	// Initialize empty graph
	New(path string, distanceFunction func(*Vertex) float64) (GraphInterface, error)

	// Loads graph from disk into memory
	Open(path string) (*GraphInterface, error)

	//  Delete existing graph
	Delete(path string) error
}

// the interface for the graph
type GraphInterface interface {

	// Search for the approximate k nearest neighbours of the object in the graph.
	// m is the number of multi searches being performed.
	NNSearch(vertex *Vertex, m uint16, k uint16) ([]*Vertex, error)

	// Insert a new object into the graph.
	// The new object will be linked to the f approximate nearest neighbours.
	// w is the number of multi searches
	NNInsert(vertex *Vertex, f uint16, w uint16) error

	// Returns the distance between to vertex depending on the graphs distance function
	CalculateDistance(v *Vertex, w *Vertex) float64

	// Saves graph to disk and frees memory
	Close() error

	String() string
}
