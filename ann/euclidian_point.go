package ann

import (
	"errors"
	"math"
)

type EuclidianPoint struct {
	dimension   uint16
	coordinates []float64
}

func NewPoint(dimension uint16, coordinates []float64) (*EuclidianPoint, error) {
	if dimension == 0 {
		return nil, errors.New("dimension must be greater than 0")
	}

	if len(coordinates) != int(dimension) {
		return nil, errors.New("size of coordinates must match dimension")
	}

	point := &EuclidianPoint{
		dimension:   dimension,
		coordinates: make([]float64, dimension),
	}

	copy(point.coordinates, coordinates)

	return point, nil
}

func (p EuclidianPoint) calculateDistance(object *ObjectInterface) float64 {
	other, ok := (*object).(*EuclidianPoint)
	if !ok {
		panic("unable to calculate distance to object that is not a EuclidianPoint")
	}

	if other.dimension != p.dimension {
		panic("unable to calculate distance between points of different dimension")
	}

	var distance float64 = 0
	for i := 0; i < int(p.dimension); i++ {
		distance += math.Sqrt(math.Pow(other.coordinates[i]-p.coordinates[i], 2))
	}

	return distance
}
