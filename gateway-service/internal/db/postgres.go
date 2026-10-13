package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(pgURL string) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), pgURL)
}
