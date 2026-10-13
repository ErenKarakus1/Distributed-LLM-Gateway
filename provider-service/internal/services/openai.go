package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/provider-service/internal/config"
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/provider-service/internal/models"
	"github.com/google/uuid"
)

type OpenAIChatRequest struct {
	Model    string               `json:"model"`
	Messages []models.ChatMessage `json:"messages"`
}

type OpenAIChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int                `json:"index"`
		Message      models.ChatMessage `json:"message"`
		FinishReason string             `json:"finish_reason"`
	} `json:"choices"`
}

func callOpenAI(req models.CreateChatCompletionRequest, cfg config.Config) (models.CreateChatCompletionResponse, error) {
	if cfg.OpenAIAPIKey == "" {
		return models.CreateChatCompletionResponse{}, errors.New("openai api key is required")
	}

	openAIReq := OpenAIChatRequest{
		Model:    cfg.OpenAIModel,
		Messages: req.Messages,
	}

	body, err := json.Marshal(openAIReq)
	if err != nil {
		return models.CreateChatCompletionResponse{}, errors.New("internal server error")
	}

	url := strings.TrimRight(cfg.OpenAIBaseURL, "/") + "/chat/completions"

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return models.CreateChatCompletionResponse{}, errors.New("internal server error")
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+cfg.OpenAIAPIKey)

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return models.CreateChatCompletionResponse{}, errors.New("internal server error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.CreateChatCompletionResponse{}, errors.New("openai request failed")
	}

	var openAIResp OpenAIChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return models.CreateChatCompletionResponse{}, errors.New("internal server error")
	}

	if len(openAIResp.Choices) == 0 {
		return models.CreateChatCompletionResponse{}, errors.New("openai returned no choices")
	}

	var choices []models.ChatCompletionChoice
	for _, choice := range openAIResp.Choices {
		choices = append(choices, models.ChatCompletionChoice{
			Index: choice.Index,
			Message: models.ChatMessage{
				Role:    choice.Message.Role,
				Content: choice.Message.Content,
			},
			FinishReason: choice.FinishReason,
		})
	}

	return models.CreateChatCompletionResponse{
		ID:        uuid.New(),
		RequestID: req.RequestID,
		Object:    "chat.completion",
		Model:     req.Model,
		Provider:  req.Provider,
		Choices:   choices,
	}, nil
}
