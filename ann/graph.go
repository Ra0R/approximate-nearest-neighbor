package ann

import "errors"

type Graph struct {
}

func (g *Graph) NNSearch(object ObjectInterface, m uint16, k uint16) ([]ObjectInterface, error) {
	return nil, errors.New("not implemented")
}

func (g *Graph) NNInsert(object ObjectInterface, f uint16, w uint16) error {
	return errors.New("not implemented")
}
