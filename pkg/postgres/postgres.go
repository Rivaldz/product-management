package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(url string, poolMax int) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(poolMax)

	var pool *pgxpool.Pool
	for i := 0; i < 5; i++ {
		pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			err = pool.Ping(context.Background())
			if err == nil {
				return pool, nil
			}
		}
		log.Printf("Postgres is trying to connect, attempt %d", i+1)
		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("postgres - New - pool.Ping: %w", err)
}
