package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ListPlayers(db *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query(context.Background(), `SELECT name, uuid FROM players ORDER BY name`)
	if err != nil {
		http.Error(w, "failed to fetch players: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Player struct {
		Name string `json:"name"`
		UUID string `json:"uuid"`
	}

	var players []Player
	for rows.Next() {
		var p Player
		if err := rows.Scan(&p.Name, &p.UUID); err != nil {
			http.Error(w, "failed to read player: "+err.Error(), http.StatusInternalServerError)
			return
		}
		players = append(players, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}
