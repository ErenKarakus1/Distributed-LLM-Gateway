package services

import (
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/router-service/internal/models"
	"github.com/google/uuid"
)

func RouteChat(req models.CreateChatCompletionRequest) (models.CreateChatCompletionResponse, error) {
	return models.CreateChatCompletionResponse{
		ID:       uuid.New(),
		Object:   "chat.completion",
		Model:    req.Model,
		Provider: "mock-provider",
		Choices: []models.ChatCompletionChoice{
			{
				Index: 0,
				Message: models.ChatMessage{
					Role:    "assistant",
					Content: "mock content",
				},
				FinishReason: "stop",
			},
		},
	}, nil
}
