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
		prompt_tokens,
		completion_tokens,
		total_tokens,
		input_cost_usd,
		output_cost_usd,
		total_cost_usd,
		duration_ms
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
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
		log.PromptTokens,
		log.CompletionTokens,
		log.TotalTokens,
		log.InputCostUSD,
		log.OutputCostUSD,
		log.TotalCostUSD,
		log.DurationMS,
	)
	if err != nil {
		return errors.New("internal server error")
	}
	return nil
}
