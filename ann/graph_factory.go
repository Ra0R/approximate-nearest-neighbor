package ann

import "errors"

type GraphFactory struct {
}

func (gf *GraphFactory) New() (GraphInterface, error) {
	return nil, errors.New("not implemented")
}
