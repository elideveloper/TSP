package main

import (
	"TSP/tsp"
	"bufio"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// поскольку тип ключа byte, то максимум 256 пунктов назначений возможно

type Destinations struct {
	InterNamesIndexes map[byte]int
	DistancesMatrix   [][]float64
	InternalNames     []byte
	NamesMap          map[byte]string
}

func NewDestinations(inputMatrix [][]string) Destinations {
	l := len(inputMatrix[0])
	nameIndexes := make(map[byte]int)
	internalNames := make([]byte, l)
	namesMap := make(map[byte]string)
	var startingInnerValue byte
	for i := 0; i < l; i++ {
		namesMap[startingInnerValue] = inputMatrix[0][i]
		internalNames[i] = startingInnerValue
		nameIndexes[startingInnerValue] = i
		startingInnerValue++
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
		InterNamesIndexes: nameIndexes,
		DistancesMatrix:   distMatrix,
		NamesMap:          namesMap,
		InternalNames:     internalNames,
	}
	return ds
}

func (d Destinations) getDistance(a, b byte) float64 {
	return d.DistancesMatrix[d.InterNamesIndexes[a]][d.InterNamesIndexes[b]]
}

func (d Destinations) getInternalNames() []byte {
	names := make([]byte, len(d.InternalNames))
	copy(names, d.InternalNames)
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

	dests := NewDestinations(fields)
	fmt.Println(dests.DistancesMatrix)
	fmt.Println(dests.InterNamesIndexes)

	fmt.Println(dests.getDistance(0, 1))

	fmt.Println(dests.getInternalNames())

	for i := 0; i < 10; i++ {
		fmt.Println(tsp.GetRandomRoute(dests.getInternalNames()))
	}

	fmt.Println(len(tsp.ExploredRoutes))

}
