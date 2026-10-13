package repository

import (
	"context"
	"errors"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

const insertLogSQL = `
	INSERT INTO request_logs (
		id,
		request_id,
		model,
		provider,
		status_code,
		duration_ms
	)
	VALUES ($1,$2,$3,$4,$5,$6)
`

func InsertLog(ctx context.Context, pool *pgxpool.Pool, log models.RequestLog) error {
	_, err := pool.Exec(
		ctx,
		insertLogSQL,
		log.ID,
		log.RequestID,
		log.Model,
		log.Provider,
		log.StatusCode,
		log.DurationMS,
	)
	if err != nil {
		return errors.New("internal server error")
	}
	return nil
}
