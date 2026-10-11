package validations

import (
	"errors"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/models"
)

func ValidateCreateChatCompletionRequest(c models.CreateChatCompletionRequest) error {
	if c.Model == "" {
		return errors.New("model is required")
	}
	if len(c.Messages) == 0 {
		return errors.New("messages are required")
	}
	for _, message := range c.Messages {
		if message.Role == "" {
			return errors.New("message role is required")
		}
		if message.Content == "" {
			return errors.New("message content is required")
		}
	}
	return nil
}
