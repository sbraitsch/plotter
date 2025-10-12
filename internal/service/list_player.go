package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ListPlayerIds(ctx context.Context, db *pgxpool.Pool) ([]Player, error) {

	rows, err := db.Query(ctx, `SELECT name, uuid FROM players ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []Player
	for rows.Next() {
		var p Player
		if err := rows.Scan(&p.Name, &p.UUID); err != nil {
			return nil, err
		}
		players = append(players, p)
	}

	return players, nil
}

func ListPlayerData(ctx context.Context, db *pgxpool.Pool) ([]PlayerData, error) {
	rows, err := db.Query(ctx, `
        SELECT p.name, pm.from_num, pm.to_num
        FROM players p
        LEFT JOIN player_mappings pm ON pm.player_id = p.id
        ORDER BY p.name, pm.from_num
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	playerMap := make(map[string]*PlayerData)

	for rows.Next() {
		var name string
		var fromNum, toNum *int
		if err := rows.Scan(&name, &fromNum, &toNum); err != nil {
			return nil, err
		}

		if _, exists := playerMap[name]; !exists {
			playerMap[name] = &PlayerData{
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
	result := make([]PlayerData, 0, len(playerMap))
	for _, pd := range playerMap {
		result = append(result, *pd)
	}

	return result, nil
}
