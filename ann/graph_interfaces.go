package ann

type GraphFactoryInterface interface {

	// Initialize empty graph
	New(path string) (GraphInterface, error)

	// Loads graph from disk into memory
	Open(path string) (*GraphInterface, error)
}

// the interface for the graph
type GraphInterface interface {

	// Search for the approximate k nearest neighbours of the object in the graph.
	// m is the number of multi searches being performed.
	NNSearch(object ObjectInterface, m uint16, k uint16) ([]ObjectInterface, error)

	// Insert a new object into the graph.
	// The new object will be linked to the f approximate nearest neighbours.
	// w is the number of multi searches
	NNInsert(object ObjectInterface, f uint16, w uint16) error

	// Saves graph to disk and frees memory
	Close() error
}

// the interface for the object to store in the graph
type ObjectInterface interface {

	// Get distance between this object and another.
	// Different implementations of ObjectInterface might not be compatible and throw an error.
	calculateDistance(object ObjectInterface) (float64, error)
}
