package ann

import "github.com/emirpasic/gods/utils"

type VertexComparator struct {
	vertex *Vertex
}

func (c *VertexComparator) compare(a, b interface{}) int {
	v, ok := a.(*Vertex)
	if !ok {
		panic("unable to compare non vertex object")
	}
	w, ok := b.(*Vertex)
	if !ok {
		panic("unable to compare non vertex object")
	}

	d1 := c.vertex.calculateDistance(v)
	d2 := c.vertex.calculateDistance(w)

	return utils.Float64Comparator(d1, d2)
}
