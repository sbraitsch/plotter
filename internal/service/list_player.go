package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PlayerDataResponse struct {
	Name     string      `json:"name"`
	PlotData map[int]int `json:"plotData"`
}

func ListPlayerIds(ctx context.Context, db *pgxpool.Pool) ([]PlayerResponse, error) {

	rows, err := db.Query(ctx, `SELECT name, uuid FROM players ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []PlayerResponse
	for rows.Next() {
		var p PlayerResponse
		if err := rows.Scan(&p.Name, &p.UUID); err != nil {
			return nil, err
		}
		players = append(players, p)
	}

	return players, nil
}

func ListPlayerData(ctx context.Context, db *pgxpool.Pool) ([]PlayerDataResponse, error) {
	rows, err := db.Query(ctx, `
        SELECT p.name, pm.from_num, pm.to_num
        FROM players p
        LEFT JOIN player_mappings pm ON pm.player_id = p.id
		WHERE p.name <> 'admin'
        ORDER BY p.name, pm.from_num
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	playerMap := make(map[string]*PlayerDataResponse)

	for rows.Next() {
		var name string
		var fromNum, toNum *int
		if err := rows.Scan(&name, &fromNum, &toNum); err != nil {
			return nil, err
		}

		if _, exists := playerMap[name]; !exists {
			playerMap[name] = &PlayerDataResponse{
				Name:     name,
				PlotData: make(map[int]int),
			}
		}

		if fromNum != nil && toNum != nil {
			playerMap[name].PlotData[*fromNum] = *toNum
		}
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	// Convert map to slice
	result := make([]PlayerDataResponse, 0, len(playerMap))
	for _, pd := range playerMap {
		result = append(result, *pd)
	}

	return result, nil
}
