package main

import (
	"TSP/tsp"
	"bufio"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// отдельный поток с получением неповторяющегося рандомного роута
// который постоянно не дожидаясь воркеров, создает маршрут и готов его отдать на исследование

// несколько потоков воркеров, которые через канал принимают маршруты
// и проделывают нужные операции для получения лучших маршрутов

// по окончанию времени/итераций, среди лучших маршрутов воркеров
// выбирается лучший из лучших (или просто получаем все)

func worker(routesChan <-chan []byte) {
	for {
		fmt.Println(<-routesChan)
		time.Sleep(time.Second * 1)
	}
}

func main() {

	rand.Seed(time.Now().UnixNano())

	routesChan := make(chan []byte)

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

	go func() {
		for {
			routesChan <- dm.GetRandomRoute()
		}
	}()

	for i := 0; i < 3; i++ {
		go worker(routesChan)
	}

	time.Sleep(time.Second * 10)

}
