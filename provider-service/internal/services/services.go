package services

import (
	"errors"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/provider-service/internal/models"
	"github.com/google/uuid"
)

var ErrUnsupportedProvider = errors.New("unsupported provider")

func generateMockContent(provider string) (string, error) {
	switch provider {
	case "mock-openai":
		return "mock response from openai", nil
	case "mock-anthropic":
		return "mock response from anthropic", nil
	case "mock-gemini":
		return "mock response from gemini", nil
	default:
		return "", ErrUnsupportedProvider
	}
}

func ProvideChat(req models.CreateChatCompletionRequest) (models.CreateChatCompletionResponse, error) {
	content, err := generateMockContent(req.Provider)
	if err != nil {
		return models.CreateChatCompletionResponse{}, err
	}
	return models.CreateChatCompletionResponse{
		ID:       uuid.New(),
		Object:   "chat.completion",
		Model:    req.Model,
		Provider: req.Provider,
		Choices: []models.ChatCompletionChoice{
			{
				Index: 0,
				Message: models.ChatMessage{
					Role:    "assistant",
					Content: content,
				},
				FinishReason: "stop",
			},
		},
	}, nil
}
