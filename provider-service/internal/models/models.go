package models

import (
	"strings"

	"github.com/google/uuid"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CreateChatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Provider string        `json:"provider"`
}

func (c *CreateChatCompletionRequest) Normalize() {
	c.Model = strings.ToLower(strings.TrimSpace(c.Model))
	c.Provider = strings.ToLower(strings.TrimSpace(c.Provider))
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
	ID       uuid.UUID              `json:"id"`
	Object   string                 `json:"object"`
	Model    string                 `json:"model"`
	Provider string                 `json:"provider"`
	Choices  []ChatCompletionChoice `json:"choices"`
}
