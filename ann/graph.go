package ann

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/emirpasic/gods/sets/treeset"
)

type DistanceFunction func(*Vertex, *Vertex) float64

type Graph struct {
	Path                 string
	nextVertexId         uint64
	vertices             map[uint64]*Vertex
	edges                map[uint64][]*Edge
	distanceFunctionName string
	distanceFunction     DistanceFunction
}

type Vertex struct {
	Id          uint64
	Coordinates []float64
}

type Edge struct {
	V uint64
	W uint64
}

func (g *Graph) NNSearch(searchVertex *Vertex, m uint16, k uint16) ([]*Vertex, error) {

	vertices, err := g.internalSearch(searchVertex, m, k)
	if err != nil {
		return nil, err
	}

	return vertices, nil
}

func (g *Graph) internalSearch(q *Vertex, m uint16, k uint16) ([]*Vertex, error) {
	comparator := &VertexComparator{
		vertex:           q,
		distanceFunction: g.distanceFunction,
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

		candidates.Clear()
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
			// remove c from candidates
			candidates.Remove(c)

			// check stop condition
			// if c is further than k-th element from result, then break repeat
			if tempRes.Size() >= int(k) {
				kthResult := tempRes.Values()[k-1].(*Vertex)
				if q.calculateDistance(g, c) > q.calculateDistance(g, kthResult) {
					break
				}
			}

			// update list of candidates
			friends, err := g.getFriends(c.Id)
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

		for j, r := range tempRes.Values() {
			if j > int(k) {
				break
			}
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
		if edge.V == vertexId {
			w, err := g.getVertex(edge.W)
			if err != nil {
				return nil, err
			}
			result = append(result, w)
		} else {
			v, err := g.getVertex(edge.V)
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
		V: v,
		W: w,
	}

	g.edges[v] = append(g.edges[v], e)

	e2 := &Edge{
		V: w,
		W: v,
	}
	g.edges[w] = append(g.edges[w], e2)
}

func (g *Graph) addVertex(vertex *Vertex) {
	g.vertices[vertex.Id] = vertex
	g.nextVertexId++
}

func (g *Graph) NNInsert(v *Vertex, f uint16, w uint16) error {
	v.Id = g.nextVertexId

	g.edges[v.Id] = make([]*Edge, 0)

	// Empty Graph
	if g.nextVertexId == 0 {
		g.addVertex(v)
		return nil
	}

	if v.Id <= uint64(f) {
		// Connect previous nodes to newly inserted node
		for i := uint64(0); i < v.Id; i++ {
			if v.Id != i { // Do not add reflexiv edges
				g.addEdge(v.Id, uint64(i))
			}
		}
	} else {
		nearestVertices, err := g.internalSearch(v, w, f)
		if err != nil {
			return err
		}

		for _, w := range nearestVertices {
			if w == nil {
				break
			}
			g.addEdge(v.Id, w.Id)
		}
	}

	g.addVertex(v)

	return nil
}

func (g *Graph) Close() error {
	// Persist vertices
	file, err := os.OpenFile(g.Path+"vertices.ann", os.O_RDWR|os.O_CREATE, 0644)

	defer file.Close()
	if err != nil {
		return err
	}

	file2, err := os.OpenFile(g.Path+"edges.ann", os.O_RDWR|os.O_CREATE, 0644)

	defer file2.Close()
	if err != nil {
		return err
	}

	file3, err := os.OpenFile(g.Path+"graph.config", os.O_RDWR|os.O_CREATE, 0644)

	defer file3.Close()
	if err != nil {
		return err
	}
	// Persist vertices
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(g.vertices)

	// Persist edges
	encoder = gob.NewEncoder(file2)
	err = encoder.Encode(g.edges)

	// Persist config
	config := GraphConfig{
		NextVertexId:         g.nextVertexId,
		Path:                 g.Path,
		DistanceFunctionName: g.distanceFunctionName}

	encoder = gob.NewEncoder(file3)
	err = encoder.Encode(config)

	return err
}

func (g *Graph) getVertex(id uint64) (*Vertex, error) {
	v := g.vertices[id]
	if v != nil {
		return v, nil
	}
	return nil, fmt.Errorf("Vertex %d not found in graph", id)
}

func (g *Graph) CalculateDistance(v, w *Vertex) float64 {
	return g.distanceFunction(v, w)
}

func (v *Vertex) calculateDistance(g *Graph, w *Vertex) float64 {
	return g.distanceFunction(v, w)
}

func (g *Graph) String() string {
	s := ""
	for i := 0; i < len(g.vertices); i++ {
		s += "Vertice " + strconv.FormatUint(g.vertices[uint64(i)].Id, 10)
		s += " -> "
		near := g.edges[g.vertices[uint64(i)].Id]
		for j := 0; j < len(near); j++ {
			s += strconv.FormatUint(near[j].W, 10) + " "
		}
		s += "\n"
	}
	fmt.Println(s)

	return s
}
