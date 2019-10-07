package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Destinations struct {
	NameIndex       map[string]int
	DistancesMatrix [][]float64
	Names           []string
}

func NewDestinations(inputMatrix [][]string) Destinations {
	l := len(inputMatrix[0])
	nameIndexes := make(map[string]int, l)
	names := make([]string, l)
	for i := 0; i < l; i++ {
		nameIndexes[inputMatrix[0][i]] = i
		names[i] = inputMatrix[0][i]
	}
	var err error
	distMatrix := make([][]float64, l)
	for i := 0; i < l; i++ {
		distMatrix[i] = make([]float64, l)
		for j := 0; j < l; j++ {
			distMatrix[i][j], err = strconv.ParseFloat(inputMatrix[i+1][j], 64)
			if err != nil {
				panic(err)
			}
		}
	}
	for i := 0; i < l; i++ {
		for j := 0; j < l; j++ {
			distMatrix[j][i] = distMatrix[i][j]
		}
	}
	ds := Destinations{
		NameIndex:       nameIndexes,
		DistancesMatrix: distMatrix,
		Names:           names,
	}
	return ds
}

func (d Destinations) getDistance(a, b string) float64 {
	return d.DistancesMatrix[d.NameIndex[a]][d.NameIndex[b]]
}

func (d Destinations) getNames() []string {
	names := make([]string, len(d.Names))
	copy(names, d.Names)
	return names
}

func main() {

	rand.Seed(time.Now().UnixNano())

	csvFile, err := os.Open("destinations.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	fields, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	csvFile.Close()

	fmt.Println(len(fields))
	fmt.Println(fields)

	dests := NewDestinations(fields)
	fmt.Println(dests.DistancesMatrix)
	fmt.Println(dests.NameIndex)

	fmt.Println(dests.getDistance("15", "6"))

	fmt.Println(dests.getNames())

	// basicRoute := []string{
	// 	"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15"}

	// for i := 0; i < 10; i++ {
	// 	tsp.GetRandomRoute(basicRoute)
	// }

	// fmt.Println(len(tsp.ExploredRoutes))

}
