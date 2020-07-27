package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/elideveloper/TSP/eval"

	"github.com/elideveloper/TSP/ga"
	"github.com/elideveloper/TSP/tsp"
)

const fileName = "destinations.csv"

// отдельный поток с получением неповторяющегося рандомного роута
// который постоянно не дожидаясь воркеров, создает маршрут и готов его отдать на исследование

// несколько потоков воркеров, которые через канал принимают маршруты
// и проделывают нужные операции для получения лучших маршрутов

// по окончанию времени/итераций, среди лучших маршрутов воркеров
// выбирается лучший из лучших (или просто получаем все)

func main() {

	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "logger: ", log.Lshortfile)
	)

	genSize := 20
	numWorkers := 10
	numGenerations := 1000

	rand.Seed(time.Now().UnixNano())

	routesChan := make(chan tsp.Route)

	csvFile, err := os.Open(fileName)
	if err != nil {
		logger.Fatal(err)
		return
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	fields, err := reader.ReadAll()
	if err != nil {
		logger.Fatal(err)
		return
	}
	csvFile.Close()

	dm := tsp.NewDataManager(fields)

	gA := ga.NewGA(genSize, numWorkers)

	go func() {
		for {
			routesChan <- dm.GetUnqRandomRoute()
		}
	}()

	gA.RunSearch(eval.Evaluate, dm, routesChan, numGenerations)

	route := gA.GetBestFoundRoute(eval.Evaluate, dm)
	fmt.Println(route, eval.Evaluate(route, dm))

}
