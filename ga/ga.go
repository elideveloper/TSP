package ga

import (
	"fmt"
	"math"
	"sync"

	"github.com/elideveloper/TSP/eval"
	"github.com/elideveloper/TSP/tsp"
)

type GA struct {
	generationSize int
	numWorkers     int
	mut            *sync.Mutex
	parents        []tsp.Route
}

func NewGA(genSize, numWorkers int) *GA {
	return &GA{
		generationSize: genSize,
		numWorkers:     numWorkers,
		mut:            &sync.Mutex{},
		parents:        make([]tsp.Route, numWorkers),
	}
}

func (g *GA) Worker(evalFunc eval.EvalFunc, dm *tsp.DataManager, routesChan <-chan tsp.Route, parentIndex int) {
	routes := make([]tsp.Route, g.generationSize)
	bestIndex := 0
	var bestScore, buffScore float64
	for i := 0; i < g.generationSize; i++ {
		routes[i] = <-routesChan
		buffScore = evalFunc(routes[i], dm)
		if bestScore == 0 {
			bestScore = buffScore
		} else {
			if buffScore < bestScore {
				bestScore = buffScore
				bestIndex = i
			}
		}
	}

	g.parents[parentIndex] = routes[bestIndex]
}

func (g *GA) RunSearch(evalFunc eval.EvalFunc, dm *tsp.DataManager, routesChan <-chan tsp.Route) {
	wg := sync.WaitGroup{}
	for i := 0; i < g.numWorkers; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, parIndex int) {
			defer wg.Done()
			g.Worker(evalFunc, dm, routesChan, parIndex)
		}(&wg, i)
	}
	wg.Wait()

	// GA operations on parents
	// and make a new generation
}

func (g GA) PrintParents(evalFunc eval.EvalFunc, dm *tsp.DataManager) {
	for _, p := range g.parents {
		fmt.Println(p, evalFunc(p, dm))
	}
}

func (g GA) GetBestFoundRoute(evalFunc eval.EvalFunc, dm *tsp.DataManager) tsp.Route {
	var bestRoute tsp.Route
	bestScore := math.MaxFloat64
	for _, p := range g.parents {
		score := evalFunc(p, dm)
		if score < bestScore {
			bestRoute = p
			bestScore = score
		}
	}
	return bestRoute
}
