package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PlayerResponse struct {
	Name    string `json:"name"`
	UUID    string `json:"uuid"`
	IsAdmin bool   `json:"idAdmin"`
}

type PlayerInit struct {
	Name    string `json:"name"`
	IsAdmin bool   `json:"isAdmin"`
}

func AddPlayer(ctx context.Context, db *pgxpool.Pool, name string) (*PlayerResponse, error) {
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

	return &PlayerResponse{Name: name, UUID: uuid}, nil
}

func InitializePlayerData(ctx context.Context, db *pgxpool.Pool, players []PlayerInit) ([]PlayerResponse, error) {

	if len(players) == 0 {
		return nil, errors.New("no players provided")
	}

	result := make([]PlayerResponse, 0, len(players))

	for _, p := range players {
		var uuid string

		err := db.QueryRow(ctx, `
			INSERT INTO players (name, is_admin)
			VALUES ($1, $2)
			ON CONFLICT (name)
			DO UPDATE SET is_admin = EXCLUDED.is_admin
			RETURNING uuid
		`, p.Name, p.IsAdmin).Scan(&uuid)

		if err != nil {
			return nil, err
		}

		result = append(result, PlayerResponse{
			Name:    p.Name,
			UUID:    uuid,
			IsAdmin: p.IsAdmin,
		})
	}

	return result, nil
}
