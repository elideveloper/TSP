package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
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

	genSize := 40
	numWorkers := 20

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

	wg := sync.WaitGroup{}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			gA.Worker(eval.Evaluate, dm, routesChan)
			wg.Done()
		}(&wg)
	}
	wg.Wait()
	time.Sleep(time.Second * 1)
	//gA.PrintParents(eval.Evaluate, dm)

	route := gA.GetBestFoundRoute(eval.Evaluate, dm)
	fmt.Println(route, eval.Evaluate(route, dm))

	// GA operations on parents
	// and make a new generation

	//time.Sleep(time.Second * 2)

}
