package ann

type VertexComparator struct {
	vertex           *Vertex
	distanceFunction DistanceFunction
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

	d1 := c.distanceFunction(c.vertex, v)
	d2 := c.distanceFunction(c.vertex, w)

	return float64Comparator(d1, d2)
}

func float64Comparator(a, b float64) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}
