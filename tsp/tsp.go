package tsp

import (
	"math/rand"
	"strings"
)

// ExploredRoutes keeps all explored routes
var ExploredRoutes = map[string]bool{}

// GetRandomRoute returns random string of all destinations
func GetRandomRoute(destinations []string) []string {
	for {
		rand.Shuffle(len(destinations), func(i, j int) {
			destinations[i], destinations[j] = destinations[j], destinations[i]
		})

		if !ExploredRoutes[strings.Join(destinations, "")] {
			ExploredRoutes[strings.Join(destinations, "")] = true
			return destinations
		}
	}
}
