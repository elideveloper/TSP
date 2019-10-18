package ga

import (
	"fmt"

	"github.com/elideveloper/TSP/eval"
	"github.com/elideveloper/TSP/tsp"
)

type GA struct {
	generationSize int
}

func NewGA(genSize int) *GA {
	return &GA{
		generationSize: genSize,
	}
}

func (g GA) Worker(evalFunc eval.EvalFunc, dm *tsp.DataManager, routesChan <-chan tsp.Route) {
	routes := make([]tsp.Route, g.generationSize)
	scores := make([]float64, g.generationSize)
	for i := 0; i < g.generationSize; i++ {
		routes[i] = <-routesChan
		scores[i] = evalFunc(routes[i], dm)
	}

	fmt.Println(scores)
}
