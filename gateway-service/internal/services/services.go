package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/models"
	"github.com/google/uuid"
)

type providerPrice struct {
	inputPerMillion  float64
	outputPerMillion float64
}

var providerPrices = map[string]providerPrice{
	"mock-openai":    {inputPerMillion: 0.15, outputPerMillion: 0.60},
	"mock-anthropic": {inputPerMillion: 0.05, outputPerMillion: 0.20},
	"mock-gemini":    {inputPerMillion: 0.10, outputPerMillion: 0.40},
	"openai":         {inputPerMillion: 0.10, outputPerMillion: 0.40},
}

func estimateCost(provider string, promptTokens int, completionTokens int) (float64, float64, float64) {
	prices, ok := providerPrices[provider]
	if !ok {
		return 0, 0, 0
	}

	inputCost := float64(promptTokens) / 1000000 * prices.inputPerMillion
	outputCost := float64(completionTokens) / 1000000 * prices.outputPerMillion
	totalCost := inputCost + outputCost

	return inputCost, outputCost, totalCost
}

func CreateChatCompletion(routerReq models.RouterRequest, routerURL string) (models.CreateChatCompletionResponse, error) {
	log.Printf("received chat completion request request_id=%v model=%v", routerReq.RequestID, routerReq.Model)
	body, err := json.Marshal(routerReq)
	if err != nil {
		return models.CreateChatCompletionResponse{}, models.UpstreamError{StatusCode: 500, Message: "internal server error"}
	}
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	httpReq, err := http.NewRequest(http.MethodPost, routerURL, bytes.NewBuffer(body))
	if err != nil {
		return models.CreateChatCompletionResponse{}, models.UpstreamError{StatusCode: 500, Message: "internal server error"}
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return models.CreateChatCompletionResponse{}, models.UpstreamError{StatusCode: 500, Message: "internal server error"}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorBody struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorBody); err != nil {
			return models.CreateChatCompletionResponse{}, models.UpstreamError{StatusCode: resp.StatusCode, Message: "upstream service error"}
		}
		if errorBody.Error == "" {
			errorBody.Error = "upstream service error"
		}
		return models.CreateChatCompletionResponse{}, models.UpstreamError{StatusCode: resp.StatusCode, Message: errorBody.Error}
	}

	var result models.CreateChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.CreateChatCompletionResponse{}, models.UpstreamError{StatusCode: 500, Message: "internal server error"}
	}

	return result, nil
}

func estimateTokens(charsLength int) int {
	if charsLength == 0 {
		return 0
	}
	return (charsLength + 3) / 4
}

func EstimatePromptTokens(messages []models.ChatMessage) int {
	totalChars := 0
	for _, msg := range messages {
		totalChars += len(msg.Content)
	}
	return estimateTokens(totalChars)
}

func estimateCompletionTokens(resp models.CreateChatCompletionResponse) int {
	totalChars := 0
	for _, choice := range resp.Choices {
		totalChars += len(choice.Message.Content)
	}
	return estimateTokens(totalChars)
}

func CreateLogSuccess(resp models.CreateChatCompletionResponse, durationMS int, promptTokens int) models.RequestLog {
	completionTokens := estimateCompletionTokens(resp)
	totalTokens := completionTokens + promptTokens
	inputCost, outputCost, totalCost := estimateCost(resp.Provider, promptTokens, completionTokens)
	return models.RequestLog{
		ID:               uuid.New(),
		RequestID:        resp.RequestID,
		Model:            resp.Model,
		Provider:         &resp.Provider,
		StatusCode:       200,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      totalTokens,
		InputCostUSD:     inputCost,
		OutputCostUSD:    outputCost,
		TotalCostUSD:     totalCost,
		DurationMS:       durationMS,
	}
}

func CreateLogFail(req models.RouterRequest, statusCode int, durationMS int, promptTokens int) models.RequestLog {
	return models.RequestLog{
		ID:           uuid.New(),
		RequestID:    req.RequestID,
		Model:        req.Model,
		StatusCode:   statusCode,
		PromptTokens: promptTokens,
		TotalTokens:  promptTokens,
		DurationMS:   durationMS,
	}
}
