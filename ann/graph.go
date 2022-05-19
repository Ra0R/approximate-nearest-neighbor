package ann

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/emirpasic/gods/sets/treeset"
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
		if v == nil {
			break
		}
		ret[i] = v.object
	}

	return ret, nil
}

func (g *Graph) internalSearch(q *Vertex, m uint16, k uint16) ([]*Vertex, error) {
	comparator := &VertexComparator{
		vertex: q,
	}
	tempRes := treeset.NewWith(comparator.compare)
	candidates := treeset.NewWith(comparator.compare)
	visitedSet := treeset.NewWith(comparator.compare)
	result := treeset.NewWith(comparator.compare)

	for i := uint16(0); i < m; i++ {
		// put random entry point in candidates
		r, err := g.getVertex(g.getRandomVertexId())
		if err != nil {
			return nil, err
		}

		candidates.Add(r)
		tempRes.Clear()
		// repeat
		for {
			// get element c closest from candidates to q
			if candidates.Size() == 0 {
				break
			}
			iter := candidates.Iterator()
			iter.First()
			c := iter.Value().(*Vertex)
			// remove c from candiadates
			candidates.Remove(c)

			// check stop condition
			// if c is further than k-th element from result, then break repeat
			if result.Size() == int(k) {
				iter = result.Iterator()
				iter.Last()
				furthestResult := iter.Value().(*Vertex)
				if q.calculateDistance(c) > q.calculateDistance(furthestResult) {
					break
				}
			}

			// update list of candidates
			friends, err := g.getFriends(c.id)
			if err != nil {
				return nil, err
			}
			for _, e := range friends {
				if e != nil && !visitedSet.Contains(e) {
					visitedSet.Add(e)
					candidates.Add(e)
					tempRes.Add(e)
				}
			}
		}

		for _, r := range tempRes.Values() {
			result.Add(r)
		}
	}

	kNearestVertices := make([]*Vertex, k)
	for i, r := range result.Values() {
		if i == int(k) {
			break
		}
		kNearestVertices[i] = r.(*Vertex)
	}

	return kNearestVertices, nil
}

// Returns a random uint64 in [0, nextVertexId)
func (g *Graph) getRandomVertexId() uint64 {
	max := int64(g.nextVertexId)
	if max > 0 {
		return uint64(rand.Int63n(max))
	} else {
		return 0
	}
}

func (g *Graph) getFriends(vertexId uint64) ([]*Vertex, error) {
	result := make([]*Vertex, 0)
	for _, edge := range g.edges[vertexId] {
		if edge == nil {
			break
		}
		if edge.v == vertexId {
			w, err := g.getVertex(edge.w)
			if err != nil {
				return nil, err
			}
			result = append(result, w)
		} else {
			v, err := g.getVertex(edge.v)
			if err != nil {
				return nil, err
			}
			result = append(result, v)
		}
	}

	return result, nil
}

func (g *Graph) addEdge(v, w uint64) {
	e := &Edge{
		v: v,
		w: w,
	}

	g.edges[v] = append(g.edges[v], e)

	e2 := &Edge{
		v: w,
		w: v,
	}
	g.edges[w] = append(g.edges[w], e2)
}

func (g *Graph) addVertex(vertex *Vertex) {
	g.vertices[vertex.id] = vertex
	g.nextVertexId++
}

func (g *Graph) NNInsert(object ObjectInterface, f uint16, w uint16) error {
	v := &Vertex{
		id:     g.nextVertexId,
		object: &object,
	}

	g.edges[v.id] = make([]*Edge, 0)

	// Empty Graph
	if g.nextVertexId == 0 {
		g.addVertex(v)
		return nil
	}

	if v.id <= uint64(f) {
		// Connect previous nodes to newly inserted node
		for i := uint64(0); i < v.id; i++ {
			if v.id != i { // Do not add reflexiv edges
				g.addEdge(v.id, uint64(i))
			}
		}
	} else {
		nearestVertices, err := g.internalSearch(v, w, f)
		if err != nil {
			return err
		}

		for _, w := range nearestVertices {
			g.addEdge(v.id, w.id)
		}
	}

	g.addVertex(v)

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

func (v *Vertex) calculateDistance(w *Vertex) float64 {
	return (*v.object).calculateDistance(w.object)
}

func (g *Graph) String() {
	s := ""
	for i := 0; i < len(g.vertices); i++ {
		s += "Vertice " + strconv.FormatUint(g.vertices[uint64(i)].id, 10)
		s += " -> "
		near := g.edges[g.vertices[uint64(i)].id]
		for j := 0; j < len(near); j++ {
			s += strconv.FormatUint(near[j].w, 10) + " "
		}
		s += "\n"
	}
	fmt.Println(s)
}
