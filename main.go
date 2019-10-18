package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/elideveloper/TSP/eval"

	"github.com/elideveloper/TSP/ga"
	"github.com/elideveloper/TSP/tsp"
)

// отдельный поток с получением неповторяющегося рандомного роута
// который постоянно не дожидаясь воркеров, создает маршрут и готов его отдать на исследование

// несколько потоков воркеров, которые через канал принимают маршруты
// и проделывают нужные операции для получения лучших маршрутов

// по окончанию времени/итераций, среди лучших маршрутов воркеров
// выбирается лучший из лучших (или просто получаем все)

func main() {

	rand.Seed(time.Now().UnixNano())

	routesChan := make(chan tsp.Route)

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

	dm := tsp.NewDataManager(fields)

	g := ga.NewGA(10)

	go func() {
		for {
			routesChan <- dm.GetRandomRoute()
		}
	}()

	for i := 0; i < 4; i++ {
		go g.Worker(eval.Evaluate, dm, routesChan)
	}

	time.Sleep(time.Second * 10)

}
