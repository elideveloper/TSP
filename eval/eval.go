package eval

import "github.com/elideveloper/TSP/tsp"

type EvalFunc func(route tsp.Route, dm *tsp.DataManager) float64

func Evaluate(route tsp.Route, dm *tsp.DataManager) float64 {
	var score float64
	for i := 0; i < len(route)-1; i++ {
		score += dm.GetDistance(route[i], route[i+1])
	}
	return score
}
