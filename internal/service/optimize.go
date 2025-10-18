package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	hungarianAlgorithm "github.com/oddg/hungarian-algorithm"
)

const PLOT_COUNT = 53

type Assignment struct {
	Battletag string `json:"player"`
	Plot      int    `json:"plot"`
	Score     int    `json:"score"`
}

func Optimize(community *Community) []Assignment {
	matrix := buildCostMatrix(community)
	mapping, _ := hungarianAlgorithm.Solve(matrix)

	return buildAssignments(mapping, community)
}

func GetAssignments(ctx context.Context, db *pgxpool.Pool, communityId string) ([]Assignment, error) {
	rows, err := db.Query(ctx, `
		SELECT battletag, plot_id, plot_score
		FROM assignments
		WHERE community_id = $1
	`, communityId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := []Assignment{}

	for rows.Next() {
		var a Assignment
		if err := rows.Scan(&a.Battletag, &a.Plot, &a.Score); err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assignments, nil
}

func UnlockCommunity(ctx context.Context, db *pgxpool.Pool, communityId string) error {
	_, err := db.Exec(ctx,
		`UPDATE communities
		SET locked = false
		WHERE id = $1`,
		communityId,
	)
	return err
}

func SaveAssignmentsAndLock(ctx context.Context, db *pgxpool.Pool, assignments []Assignment, communityId string) error {
	sqlStr := `INSERT INTO assignments (battletag, community_id, plot_id, plot_score) VALUES `
	args := []any{}

	for i, a := range assignments {
		idx := i * 2
		sqlStr += fmt.Sprintf("($%d, $%d, $%d, $%d),", idx+1, idx+2, idx+3, idx+4)
		args = append(args, a.Battletag, communityId, a.Plot, a.Score)
	}

	sqlStr = strings.TrimSuffix(sqlStr, ",")
	sqlStr += ` ON CONFLICT (battletag)
		        DO UPDATE SET
                  plot_id = EXCLUDED.plot_id,
                  plot_score = EXCLUDED.plot_score`

	_, err := db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx,
		`UPDATE communities
		SET locked = true
		WHERE id = $1`,
		communityId,
	)

	return err
}

func buildAssignments(mapping []int, community *Community) []Assignment {
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

func buildCostMatrix(community *Community) [][]int {

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
