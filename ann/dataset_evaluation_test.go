package ann

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/bsm/mlmetrics"
	"github.com/stretchr/testify/assert"
)

func loadDataset(filePath string, skipHeader bool, nameRow int, coordinateStart int, coordinatesEnd int, delimiter rune) []EuclidianPoint {
	dataset := []EuclidianPoint{}

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Unable to read input file %s", err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = delimiter
	i := 0
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic("Error reading file" + err.Error())
		}

		var row []float64

		if i == 0 && skipHeader {
			i++
			continue
		}

		name := ""
		for j, s := range rec {
			name += s

			if j < coordinateStart || j >= coordinatesEnd {
				continue
			}
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				row = append(row, f)
			} else {
				fmt.Println(err)
			}
		}
		l := len(row)
		dataset = append(dataset, EuclidianPoint{name, uint16(l), row})
	}

	return dataset
}

func randomDataset(size int, dim int) []EuclidianPoint {
	var dataset []EuclidianPoint
	for point := 0; point < size; point++ {
		coordinates := make([]float64, dim)
		for i := 0; i < dim; i++ {
			coordinates[i] = rand.Float64()
		}
		dataset = append(dataset, EuclidianPoint{"", uint16(dim), coordinates[:]})
	}

	return dataset
}

func (lhs *EuclidianPoint) distance(rhs *EuclidianPoint) float64 {
	distance := float64(0.0)
	for i := 0; i < int(lhs.dimension); i++ {
		distance += math.Sqrt(math.Pow(rhs.coordinates[i]-lhs.coordinates[i], 2))
	}

	return distance
}
func findPointByName(dataset []EuclidianPoint, name string) int {
	r := -1

	for j := range dataset {
		if dataset[j].name == name {
			return j
		}
	}

	return r
}

// Might need to increase test parameter "-timout" > 30s, when running test
func TestInsertDatasetIntoGraph(t *testing.T) {
	assert := assert.New(t)

	//dataset := loadDataset("data.csv", true, -1, 0, 14, ',')
	dataset := loadDataset("data_ch.csv", false, 2, 9, 11, ';')

	// Setup Graph
	factory := GraphFactory{}
	path := "."
	graph, err := factory.New(path)
	assert.NoError(err)
	assert.NotNil(graph)
	var f uint16 = 4
	var w uint16 = 2
	timer := time.Now()
	for i, j := range rand.Perm(len(dataset)) {
		if i > 0 && i%100 == 0 {
			fmt.Println("Inserting 100 elements took " + time.Since(timer).String())
			timer = time.Now()
		}
		point := dataset[j]
		graph.NNInsert(point, f, w)
	}
	_ = graph.String()
}

type neighbor struct {
	nodeId   int
	distance float64
}

func TestCalculateRecallAccuracyConfusionMatrix(t *testing.T) {

	const k uint16 = 10 // Evaluate on k nearest neighbors

	// Insertion parameters
	const f uint16 = 4
	const w uint16 = 3

	var m uint16 = w
	assert := assert.New(t)
	//dataset := loadDataset("data.csv", true, -1, 0, 14, ',')
	dataset := loadDataset("data_ch.csv", false, 2, 9, 11, ';')
	fmt.Println("Lenght of dataset ", len(dataset))

	// Calculate distance from every point to every point
	// and store k nearest neighbors
	test_nearest_neighbors := make(map[int][]neighbor)
	for i, point1 := range dataset {

		// Calculate distances to all points
		var distances []neighbor

		for j, point2 := range dataset {
			distances = append(distances, neighbor{j, point1.distance(&point2)})
		}

		// Sort by distance
		sort.Slice(distances, func(i, j int) bool {
			return distances[i].distance < distances[j].distance
		})

		// Take first k
		test_nearest_neighbors[i] = distances[:k]

		if i > 0 && i%100 == 0 {
			fmt.Println("Calculated exact NN for " + strconv.Itoa(i) + " points")
		}
	}

	// Setup Graph
	factory := GraphFactory{}
	path := "."
	graph, err := factory.New(path)
	assert.NoError(err)
	assert.NotNil(graph)

	for i, j := range rand.Perm(len(dataset)) {
		point := dataset[j]
		graph.NNInsert(point, f, w)
		if i > 0 && i%100 == 0 {
			fmt.Println("Inserted " + strconv.Itoa(i) + " points to graph")
		}
	}

	// Evaluate nearest neighbors (ids)
	predictionsO := make(map[int][]ObjectInterface)
	for j, point := range dataset {

		preds, _ := graph.NNSearch(point, m, k)

		var predO []ObjectInterface
		for _, pred := range preds {
			if pred != nil {
				predO = append(predO, *pred)
			} else {
				assert.Fail("Unable to get k approximate nearest neighbors")
			}
		}
		predictionsO[j] = predO
		if j > 0 && j%100 == 0 {
			fmt.Println("Got " + strconv.Itoa(j) + " predictions")
		}
	}

	_ = graph.String()

	mat := mlmetrics.NewConfusionMatrix()
	for i := range test_nearest_neighbors {
		for j := range test_nearest_neighbors[i] {
			if len(predictionsO[i]) <= j {
				continue
			}
			x := findPointByName(dataset, predictionsO[i][j].(EuclidianPoint).name)
			mat.Observe(test_nearest_neighbors[i][j].nodeId, x)
		}
		if i > 0 && i%100 == 0 {
			fmt.Println("Observed " + strconv.Itoa(i) + " points")
		}
	}

	// print matrix
	for i := 0; i < mat.Order(); i++ {
		fmt.Println(mat.Row(i))
	}

	// print metrics
	fmt.Println()
	fmt.Printf("accuracy : %.3f\n", mat.Accuracy())
}

// Method used to create to graph for the presentation
func TestCalculateRecallAccuracyConfusionMatrix_Presentation(t *testing.T) {

	assert := assert.New(t)

	var testsetW = [4]uint16{2, 3, 4}
	var testsetF = [7]uint16{10, 11, 12, 13, 14, 16, 17}
	var results = [7][4]float64{}

	const k uint16 = 10 // Evaluate on k nearest neighbors
	dataset := loadDataset("data_ch.csv", false, 2, 9, 11, ';')
	fmt.Println("Lenght of dataset ", len(dataset))

	// Calculate distance from every point to every point
	// and store k nearest neighbors
	test_nearest_neighbors := make(map[int][]neighbor)
	for i, point1 := range dataset {

		// Calculate distances to all points
		var distances []neighbor

		for j, point2 := range dataset {
			distances = append(distances, neighbor{j, point1.distance(&point2)})
		}

		// Sort by distance
		sort.Slice(distances, func(i, j int) bool {
			return distances[i].distance < distances[j].distance
		})

		// Take first k
		test_nearest_neighbors[i] = distances[:k]

		// Sort by id for evaluation
		sort.Slice(test_nearest_neighbors[i], func(k, l int) bool {
			return test_nearest_neighbors[i][k].nodeId < test_nearest_neighbors[i][l].nodeId
		})

	}

	for w_i := range testsetW {
		m := testsetW[w_i]
		w := testsetW[w_i]
		for f_i := range testsetF {
			f := testsetF[f_i]
			// Setup Graph
			factory := GraphFactory{}
			path := "."
			graph, err := factory.New(path)
			assert.NoError(err)
			assert.NotNil(graph)

			for j := range rand.Perm(len(dataset)) {
				point := dataset[j]
				graph.NNInsert(point, f, w)

			}

			// Evaluate nearest neighbors (ids)

			predictionsO := make(map[int][]ObjectInterface)
			for j, point := range dataset {

				preds, _ := graph.NNSearch(point, m, k)

				var predO []ObjectInterface
				for _, pred := range preds {
					if pred != nil {
						predO = append(predO, *pred)
					} else {
						assert.Fail("Unable to get k approximate nearest neighbors")
					}
				}
				predictionsO[j] = predO
				if j > 0 && j%100 == 0 {
					fmt.Println("Got " + strconv.Itoa(j) + " predictions")
				}
			}

			mat := mlmetrics.NewConfusionMatrix()
			for i := range test_nearest_neighbors {

				for j := range test_nearest_neighbors[i] {
					if len(predictionsO[i]) <= j {
						continue
					}
					x := findPointByName(dataset, predictionsO[i][j].(EuclidianPoint).name)

					mat.Observe(test_nearest_neighbors[i][j].nodeId, x)
				}

			}

			results[f_i][w_i] = mat.Accuracy()
		}
	}
	fmt.Println(results)
}
