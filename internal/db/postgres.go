package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectWithRetry(ctx context.Context, dbURL string, maxRetries int, retryDelay time.Duration) *pgxpool.Pool {
	var pool *pgxpool.Pool
	var err error

	for i := 1; i <= maxRetries; i++ {
		pool, err = pgxpool.New(ctx, dbURL)
		if err == nil {
			pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
			err = pool.Ping(pingCtx)
			cancel()

			if err == nil {
				log.Println("✅ Connected to Postgres!")
				return pool
			}
			pool.Close()
		}

		log.Printf("⏳ Postgres not ready (attempt %d/%d): %v", i, maxRetries, err)
		time.Sleep(retryDelay)
	}

	log.Fatalf("❌ Could not connect to Postgres after %d attempts: %v", maxRetries, err)
	return nil
}
