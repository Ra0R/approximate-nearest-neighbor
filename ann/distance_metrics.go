package ann

import (
	"math"
	"strings"
)

func getDistanceFunctionByName(name string) DistanceFunction {
	name = strings.ToLower(name)
	switch {
	case name == "euclidean":
		return euclideanDistance
	case name == "euclideanfast":
		return euclideanDistanceFast
	}

	return nil
}

func euclideanDistance(v *Vertex, w *Vertex) float64 {
	var distance float64 = 0
	for i := 0; i < len(v.Coordinates); i++ {
		distance += math.Sqrt(math.Pow(v.Coordinates[i]-w.Coordinates[i], 2))
	}
	return distance
}

// https://stackoverflow.com/a/58134438/14204586, might be faster without taking sqrt and it should not change behaviour
func euclideanDistanceFast(v *Vertex, w *Vertex) float64 {
	var distance float64 = 0
	for i := 0; i < len(v.Coordinates); i++ {
		distance += math.Pow(v.Coordinates[i]-w.Coordinates[i], 2)
	}
	return distance
}
