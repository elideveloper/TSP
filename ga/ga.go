package ga

import (
	"fmt"
	"sync"

	"github.com/elideveloper/TSP/eval"
	"github.com/elideveloper/TSP/tsp"
)

type GA struct {
	generationSize int
	mut            *sync.Mutex
	parents        []tsp.Route
}

func NewGA(genSize int) *GA {
	return &GA{
		generationSize: genSize,
		mut:            &sync.Mutex{},
		parents:        make([]tsp.Route, 0),
	}
}

func (g *GA) Worker(evalFunc eval.EvalFunc, dm *tsp.DataManager, routesChan <-chan tsp.Route) {
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

	g.mut.Lock()
	g.parents = append(g.parents, routes[bestIndex])
	g.mut.Unlock()

}

func (g GA) PrintParents(evalFunc eval.EvalFunc, dm *tsp.DataManager) {
	for _, p := range g.parents {
		fmt.Println(p, evalFunc(p, dm))
	}
}
