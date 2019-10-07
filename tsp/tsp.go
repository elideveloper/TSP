package tsp

import (
	"math/rand"
)

// ExploredRoutes keeps all explored routes
var ExploredRoutes = map[string]bool{}

// GetRandomRoute returns random string of all destinations
func GetRandomRoute(destinations []byte) []byte {
	for {
		for i := len(destinations) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			destinations[i], destinations[j] = destinations[j], destinations[i]
		}

		if !ExploredRoutes[string(destinations)] {
			ExploredRoutes[string(destinations)] = true
			return destinations
		}
	}
}
