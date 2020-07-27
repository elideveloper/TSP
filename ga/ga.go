package ga

import (
	"fmt"
	"math"
	"math/rand"
	"sync"

	"github.com/elideveloper/TSP/eval"
	"github.com/elideveloper/TSP/tsp"
)

// crossover by 2 points
// mutation with given probability

type GA struct {
	generationSize int
	numWorkers     int
	parents        []tsp.Route
}

func NewGA(genSize, numWorkers int) *GA {
	return &GA{
		generationSize: genSize,
		numWorkers:     numWorkers,
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

	// TODO reuse found parents from last generation and not only generate a new randoms

	g.parents[parentIndex] = routes[bestIndex]
}

func (g *GA) RunSearch(evalFunc eval.EvalFunc, dm *tsp.DataManager, routesChan <-chan tsp.Route, numGenerations int) {

	for j := 0; j < numGenerations; j++ {
		wg := sync.WaitGroup{}
		for i := 0; i < g.numWorkers; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, parIndex int) {
				defer wg.Done()
				g.Worker(evalFunc, dm, routesChan, parIndex)
			}(&wg, i)
		}
		wg.Wait()

		// get recombinations from parents
		g.parents = buildNewGeneration(g.parents)
	}

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

func buildNewGeneration(generation []tsp.Route) []tsp.Route {
	newGeneration := make([]tsp.Route, 0)
	for j := 0; j < len(generation); j++ {
		if j == len(generation)-1 {
			// recombination of last and first parents
			newGeneration = append(newGeneration, getRecombination(generation[j], generation[0]))
		} else {
			newGeneration = append(newGeneration, getRecombination(generation[j], generation[j+1]))
		}
	}
	return newGeneration
}

func getRecombination(left, right tsp.Route) tsp.Route {

	x := rand.Intn(len(left))

	child := make(tsp.Route, x)
	copy(child, left[:x])

	existingMap := make(map[byte]struct{})

	for i := 0; i < x; i++ {
		existingMap[child[i]] = struct{}{}
	}
	for i := 0; i < len(right); i++ {
		if _, ok := existingMap[right[i]]; !ok {
			child = append(child, right[i])
		}
	}

	// just append initial 'home' point
	return append(child, child[0])
}
