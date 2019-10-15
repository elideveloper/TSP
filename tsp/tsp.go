package tsp

import (
	"math/rand"
	"strconv"
)

// поскольку тип ключа byte, то максимум 256 пунктов назначений возможно

// надо ограничить на 10М мапу исследованных, возможно далее сделать оптимизацию
// чтобы хранились паттерны долгих путей и за это убирались баллы, или полностью отбрасывались

// надо чтобы воркеры параллельно находили лучшие роуты, затем синхронились и создавали пул из лучших,
// затем каждый воркер работал с этим пулом лучших

// в конце концов можно будет прикрепить http api функцию, которая принимает массив входных данных =)
// и запустить это в контейнере на удаленном серваке

// DataManager is a class representing input destinations
type DataManager struct {
	InterNamesIndexes map[byte]int
	DistancesMatrix   [][]float64
	InternalNames     []byte
	NamesMap          map[byte]string
	exploredRoutes    map[string]bool
}

// NewDataManager is a constructor of DataManager
func NewDataManager(inputMatrix [][]string) DataManager {
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
	ds := DataManager{
		InterNamesIndexes: nameIndexes,
		DistancesMatrix:   distMatrix,
		NamesMap:          namesMap,
		InternalNames:     internalNames,
		exploredRoutes:    make(map[string]bool),
	}
	return ds
}

// GetDistance computes distance between two destinations using internal names
func (d DataManager) GetDistance(a, b byte) float64 {
	return d.DistancesMatrix[d.InterNamesIndexes[a]][d.InterNamesIndexes[b]]
}

// GetInternalNames returns a new slice of internal names
func (d DataManager) GetInternalNames() []byte {
	names := make([]byte, len(d.InternalNames))
	copy(names, d.InternalNames)
	return names
}

// GetRandomRoute returns random string of all destinations
func (d DataManager) GetRandomRoute() []byte {
	destinations := d.GetInternalNames()
	for {
		for i := len(destinations) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			destinations[i], destinations[j] = destinations[j], destinations[i]
		}

		if !d.exploredRoutes[string(destinations)] {
			d.exploredRoutes[string(destinations)] = true
			return destinations
		}
	}
}
