package optimizer

import (
	"fmt"
	"math/rand"
	"os"

	"gopkg.in/yaml.v3"
)

func Generate(numPlayers, numPlots, numWeights int) {

	cfg := Config{
		Plots:   numPlots,
		Players: make([]PlayerConfig, numPlayers),
	}

	desirability := make([]int, numPlots)
	for i := range numPlots {
		switch {
		case i < numPlots/10:
			// top 10% of plots are extremely desirable
			desirability[i] = 30
		case i < numPlots/5:
			// next 10% moderately desirable
			desirability[i] = 10
		case i < numPlots/2:
			// next 30% slightly desirable
			desirability[i] = 4
		default:
			// remaining 50% very low desirability
			desirability[i] = 1
		}
	}

	// Simple adjacency for demo (circular)

	// for i := 1; i <= numPlots; i++ {
	// 	neighbors := []int{}
	// 	if i > 1 {
	// 		neighbors = append(neighbors, i-1)
	// 	}
	// 	if i < numPlots {
	// 		neighbors = append(neighbors, i+1)
	// 	}
	// 	cfg.Adjacency[i] = neighbors
	// }

	// Assign 3 weighted plots per player with overlap
	for p := range numPlayers {
		player := PlayerConfig{
			Name:    fmt.Sprintf("Player%d", p+1),
			Weights: make(map[int]int),
		}

		choices := map[int]struct{}{}
		for len(choices) < numWeights {
			plot := pickWeightedPlot(desirability, numPlots)
			choices[plot] = struct{}{}
		}

		// assign incremental weights
		weight := 1
		for plot := range choices {
			player.Weights[plot] = weight
			weight++
		}

		cfg.Players[p] = player
	}

	// Write YAML
	f, err := os.Create("generated.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)
	if err := enc.Encode(cfg); err != nil {
		panic(err)
	}

	fmt.Println("Wrote generated.yaml")
}

func pickWeightedPlot(desirability []int, fallback int) int {
	total := 0
	for _, w := range desirability {
		total += w
	}
	r := rand.Intn(total)
	for i, w := range desirability {
		if r < w {
			return i + 1
		}
		r -= w
	}
	return fallback
}
