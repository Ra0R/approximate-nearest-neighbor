package ann

import "errors"

var (
	// ErrInvalidPath is returned when the path that has been given is not valid (inexistent/not writable)
	ErrInvalidPath = errors.New("Supplied argument 'path' is not valid")
)

type GraphFactory struct {
}

func (gf *GraphFactory) New() (GraphInterface, error) {
	return nil, errors.New("not implemented")
}

func (gf *GraphFactory) Open(path string) (GraphInterface, error) {
	return nil, errors.New("not implemented")
}
