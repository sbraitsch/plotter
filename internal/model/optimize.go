package model

import (
	hungarianAlgorithm "github.com/oddg/hungarian-algorithm"
)

const PLOT_COUNT = 53

func (community *CommunityData) Optimize() []Assignment {
	matrix := buildCostMatrix(community)
	mapping, _ := hungarianAlgorithm.Solve(matrix)

	return buildAssignments(mapping, community)
}

func buildAssignments(mapping []int, community *CommunityData) []Assignment {
	n := len(community.Members)
	m := PLOT_COUNT

	assignments := []Assignment{}

	for i, plotIndex := range mapping {
		plotId := plotIndex + 1

		if i < n && plotIndex < m {
			member := community.Members[i]

			weight := m
			if w, ok := member.PlotData[plotId]; ok {
				weight = w
			}

			assignments = append(assignments, Assignment{
				Battletag: member.BattleTag,
				Plot:      plotId,
				Score:     weight,
			})
		}
	}
	return assignments
}

func buildCostMatrix(community *CommunityData) [][]int {

	n := len(community.Members)
	m := PLOT_COUNT

	matrix := make([][]int, m)

	for i, member := range community.Members {
		row := make([]int, m)

		for plot := 1; plot <= PLOT_COUNT; plot++ {
			if w, ok := member.PlotData[plot]; ok {
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
