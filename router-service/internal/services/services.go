package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/router-service/internal/models"
)

var ErrUnsupportedModel = errors.New("unsupported model")

func selectProvider(model string) (string, error) {
	switch model {
	case "smart":
		return "mock-openai", nil
	case "cheap":
		return "mock-anthropic", nil
	case "fast":
		return "mock-gemini", nil
	default:
		return "", ErrUnsupportedModel
	}
}

func RouteChat(req models.CreateChatCompletionRequest, providerURL string) (models.CreateChatCompletionResponse, error) {
	provider, err := selectProvider(req.Model)
	if err != nil {
		return models.CreateChatCompletionResponse{}, err
	}
	reqWithProvider := models.ProviderChatRequest{Provider: provider, CreateChatCompletionRequest: req}
	log.Printf("routing chat request request_id=%v model=%v provider=%v", reqWithProvider.RequestID, reqWithProvider.Model, reqWithProvider.Provider)
	body, err := json.Marshal(reqWithProvider)
	if err != nil {
		return models.CreateChatCompletionResponse{}, errors.New("internal server error")
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	httpReq, err := http.NewRequest(http.MethodPost, providerURL, bytes.NewBuffer(body))
	if err != nil {
		return models.CreateChatCompletionResponse{}, errors.New("internal server error")
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return models.CreateChatCompletionResponse{}, errors.New("internal server error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.CreateChatCompletionResponse{}, fmt.Errorf("router service returned status %v", resp.StatusCode)
	}

	var result models.CreateChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.CreateChatCompletionResponse{}, errors.New("internal server error")
	}

	return result, nil
}
