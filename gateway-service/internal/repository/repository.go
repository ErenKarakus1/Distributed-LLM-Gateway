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

const getLogsSQL = `
	SELECT
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
		duration_ms,
		created_at
	FROM request_logs
	ORDER BY created_at DESC
	LIMIT 50
`

const usageSummaryTotalsSQL = `
	SELECT
		COUNT(*),
		COUNT(*) FILTER (WHERE status_code >= 200 AND status_code < 300),
		COUNT(*) FILTER (WHERE status_code < 200 OR status_code >= 300),
		COALESCE(SUM(total_tokens), 0),
		COALESCE(SUM(total_cost_usd), 0),
		COALESCE(AVG(duration_ms), 0)
	FROM request_logs
`

const usageSummaryModelsSQL = `
	SELECT 
		model, 
		COUNT(*)
	FROM request_logs
	GROUP BY model
`

const usageSummaryProvidersSQL = `
	SELECT 
		provider, 
		COUNT(*)
	FROM request_logs
	WHERE provider IS NOT NULL
	GROUP BY provider
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

func GetLogs(ctx context.Context, pool *pgxpool.Pool) ([]models.RequestLog, error) {
	rows, err := pool.Query(
		ctx,
		getLogsSQL,
	)
	if err != nil {
		return []models.RequestLog{}, errors.New("internal server error")
	}
	defer rows.Close()

	var logs []models.RequestLog
	for rows.Next() {
		var log models.RequestLog
		err := rows.Scan(
			&log.ID,
			&log.RequestID,
			&log.Model,
			&log.Provider,
			&log.StatusCode,
			&log.PromptTokens,
			&log.CompletionTokens,
			&log.TotalTokens,
			&log.InputCostUSD,
			&log.OutputCostUSD,
			&log.TotalCostUSD,
			&log.DurationMS,
			&log.CreatedAt,
		)
		if err != nil {
			return []models.RequestLog{}, errors.New("internal server error")
		}
		logs = append(logs, log)
	}
	if err := rows.Err(); err != nil {
		return []models.RequestLog{}, errors.New("internal server error")
	}
	return logs, nil
}

func GetUsageSummary(ctx context.Context, pool *pgxpool.Pool) (models.UsageSummary, error) {
	var summary models.UsageSummary
	summary.RequestsByModel = make(map[string]int64)
	summary.RequestsByProvider = make(map[string]int64)

	err := pool.QueryRow(
		ctx,
		usageSummaryTotalsSQL,
	).Scan(
		&summary.TotalRequests,
		&summary.SuccessfulRequests,
		&summary.FailedRequests,
		&summary.TotalTokens,
		&summary.TotalCostUSD,
		&summary.AverageDurationMS,
	)
	if err != nil {
		return models.UsageSummary{}, errors.New("internal server error")
	}

	modelRows, err := pool.Query(
		ctx,
		usageSummaryModelsSQL,
	)
	if err != nil {
		return models.UsageSummary{}, errors.New("internal server error")
	}
	defer modelRows.Close()
	for modelRows.Next() {
		var model string
		var count int64
		if err := modelRows.Scan(&model, &count); err != nil {
			return models.UsageSummary{}, errors.New("internal server error")
		}
		summary.RequestsByModel[model] = count
	}
	if err := modelRows.Err(); err != nil {
		return models.UsageSummary{}, errors.New("internal server error")
	}

	providerRows, err := pool.Query(
		ctx,
		usageSummaryProvidersSQL,
	)
	if err != nil {
		return models.UsageSummary{}, errors.New("internal server error")
	}
	defer providerRows.Close()

	for providerRows.Next() {
		var provider string
		var count int64
		if err := providerRows.Scan(&provider, &count); err != nil {
			return models.UsageSummary{}, errors.New("internal server error")
		}
		summary.RequestsByProvider[provider] = count
	}
	if err := providerRows.Err(); err != nil {
		return models.UsageSummary{}, errors.New("internal server error")
	}

	return summary, nil
}
