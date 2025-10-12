package optimizer

import (
	hungarianAlgorithm "github.com/oddg/hungarian-algorithm"
)

func RunHungarian(cfg Config) []Assignment {
	matrix := buildCostMatrix(cfg)
	mapping, _ := hungarianAlgorithm.Solve(matrix)

	maxCost := 0
	for _, player := range cfg.Players {
		for _, w := range player.Weights {
			if w > maxCost {
				maxCost = w
			}
		}
	}

	return buildAssignments(mapping, cfg)
}

func buildAssignments(mapping []int, cfg Config) []Assignment {
	n := len(cfg.Players)
	m := cfg.Plots
	assignments := []Assignment{}

	for i, plotIndex := range mapping {
		plotId := plotIndex + 1

		if i < n && plotIndex < m {
			player := cfg.Players[i]

			weight := m
			if w, ok := player.Weights[plotId]; ok {
				weight = w
			}

			assignments = append(assignments, Assignment{
				Player:  player.Name,
				Plot:    plotId,
				Score:   weight,
				Cheater: player.Cheater,
			})
		}
	}
	return assignments
}

func buildCostMatrix(cfg Config) [][]int {
	m := cfg.Plots
	n := len(cfg.Players)

	matrix := make([][]int, m)

	for i, player := range cfg.Players {
		row := make([]int, m)

		for plot := 1; plot <= cfg.Plots; plot++ {
			if w, ok := player.Weights[plot]; ok {
				row[plot-1] = w
			} else {
				row[plot-1] = m
			}
		}

		matrix[i] = row
	}

	// padding matrix if less players than plots
	for i := n; i < m; i++ {
		row := make([]int, m)
		for j := range row {
			row[j] = 1000
		}
		matrix[i] = row
	}

	return matrix
}
