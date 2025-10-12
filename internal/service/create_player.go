package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InitializePlayersRequest struct {
	Names []string `json:"names"`
}

type AddPlayerRequest struct {
	Name string `json:"name"`
}

func AddPlayer(ctx context.Context, db *pgxpool.Pool, name string) (*Player, error) {
	if name == "" {
		return nil, errors.New("no player name provided")
	}

	var uuid string
	err := db.QueryRow(ctx,
		`INSERT INTO players (name) VALUES ($1)
             ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
             RETURNING uuid`,
		name,
	).Scan(&uuid)
	if err != nil {
		return nil, err
	}

	return &Player{Name: name, UUID: uuid}, nil
}

func InitializePlayerData(ctx context.Context, db *pgxpool.Pool, names []string) ([]Player, error) {

	if len(names) == 0 {
		return nil, errors.New("no player names provided")
	}

	result := make([]Player, 0, len(names))

	for _, name := range names {
		var uuid string

		err := db.QueryRow(ctx,
			`INSERT INTO players (name) VALUES ($1)
             ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
             RETURNING uuid`,
			name,
		).Scan(&uuid)
		if err != nil {
			return nil, err
		}

		result = append(result, Player{
			Name: name,
			UUID: uuid,
		})
	}

	return result, nil
}
