package ann

import (
	"errors"
	"fmt"
)

type Graph struct {
	nextVertexId uint64
	vertices     map[uint64]*Vertex
	edges        map[uint64][]*Edge
}

type Vertex struct {
	id     uint64
	object *ObjectInterface
}

type Edge struct {
	v uint64
	w uint64
}

func (g *Graph) NNSearch(object ObjectInterface, m uint16, k uint16) ([]*ObjectInterface, error) {
	searchVertex := &Vertex{
		id:     0,
		object: &object,
	}
	vertices, err := g.internalSearch(searchVertex, m, k)
	if err != nil {
		return nil, err
	}

	ret := make([]*ObjectInterface, k)
	for i, v := range vertices {
		ret[i] = v.object
	}

	return ret, nil
}

func (g *Graph) internalSearch(vertex *Vertex, m uint16, k uint16) ([]*Vertex, error) {
	return nil, errors.New("not implemented")
}

func (g *Graph) NNInsert(object ObjectInterface, f uint16, w uint16) error {
	v := &Vertex{
		id:     g.nextVertexId,
		object: &object,
	}

	g.edges[v.id] = make([]*Edge, f)

	if v.id <= uint64(f) {
		for i := 0; i < int(f); i++ {
			e := &Edge{
				v: v.id,
				w: uint64(i),
			}
			g.edges[v.id] = append(g.edges[v.id], e)
			g.edges[uint64(i)] = append(g.edges[uint64(i)], e)
		}
	} else {
		nearestVertices, err := g.internalSearch(v, w, f)
		if err != nil {
			return err
		}

		for _, w := range nearestVertices {
			e := &Edge{
				v: v.id,
				w: w.id,
			}
			g.edges[v.id] = append(g.edges[v.id], e)
			g.edges[w.id] = append(g.edges[w.id], e)
		}
	}

	g.vertices[v.id] = v
	g.nextVertexId++

	return nil
}

func (g *Graph) Close() error {
	// TODO: implement this
	return nil
}

func (g *Graph) getVertex(id uint64) (*Vertex, error) {
	v := g.vertices[id]
	if v != nil {
		return v, nil
	}
	return nil, fmt.Errorf("Vertex %d not found in graph", id)
}
