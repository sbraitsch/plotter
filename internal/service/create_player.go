package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CreatePlayerRequest struct {
	Name string `json:"name"`
}

func CreatePlayer(db *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {

	var req CreatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "missing name", http.StatusBadRequest)
		return
	}

	var uuid string
	err := db.QueryRow(context.Background(),
		`INSERT INTO players (name) VALUES ($1) RETURNING uuid`, req.Name).Scan(&uuid)
	if err != nil {
		http.Error(w, "failed to create player: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"name": req.Name,
		"uuid": uuid,
	})
}
