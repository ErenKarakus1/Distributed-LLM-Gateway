package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CreateChatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

func (c *CreateChatCompletionRequest) Normalize() {
	c.Model = strings.ToLower(strings.TrimSpace(c.Model))
	for i, message := range c.Messages {
		c.Messages[i].Role = strings.ToLower(strings.TrimSpace(message.Role))
		c.Messages[i].Content = strings.TrimSpace(message.Content)
	}
}

type ChatCompletionChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type CreateChatCompletionResponse struct {
	ID        uuid.UUID              `json:"id"`
	RequestID uuid.UUID              `json:"request_id"`
	Object    string                 `json:"object"`
	Model     string                 `json:"model"`
	Provider  string                 `json:"provider"`
	Choices   []ChatCompletionChoice `json:"choices"`
}

type UpstreamError struct {
	StatusCode int
	Message    string
}

func (e UpstreamError) Error() string {
	if e.Message == "" {
		return "upstream service error"
	}
	return e.Message
}

type RouterRequest struct {
	RequestID uuid.UUID `json:"request_id"`
	CreateChatCompletionRequest
}

type RequestLog struct {
	ID               uuid.UUID `json:"id"`
	RequestID        uuid.UUID `json:"request_id"`
	Model            string    `json:"model"`
	Provider         *string   `json:"provider,omitempty"`
	StatusCode       int       `json:"status_code"`
	DurationMS       int       `json:"duration_ms"`
	PromptTokens     int       `json:"prompt_tokens"`
	CompletionTokens int       `json:"completion_tokens"`
	TotalTokens      int       `json:"total_tokens"`
	InputCostUSD     float64   `json:"input_cost_usd"`
	OutputCostUSD    float64   `json:"output_cost_usd"`
	TotalCostUSD     float64   `json:"total_cost_usd"`
	CreatedAt        time.Time `json:"created_at"`
}

type UsageSummary struct {
	TotalRequests      int64            `json:"total_requests"`
	SuccessfulRequests int64            `json:"successful_requests"`
	FailedRequests     int64            `json:"failed_requests"`
	TotalTokens        int64            `json:"total_tokens"`
	TotalCostUSD       float64          `json:"total_cost_usd"`
	AverageDurationMS  float64          `json:"average_duration_ms"`
	RequestsByModel    map[string]int64 `json:"requests_by_model"`
	RequestsByProvider map[string]int64 `json:"requests_by_provider"`
}
