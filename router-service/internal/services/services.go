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

const maxProviderAttempts = 2

var ErrUnsupportedModel = errors.New("unsupported model")

func selectProviders(model string) ([]string, error) {
	switch model {
	case "smart":
		return []string{"mock-openai", "mock-anthropic"}, nil
	case "cheap":
		return []string{"mock-anthropic", "mock-gemini"}, nil
	case "fast":
		return []string{"mock-gemini", "mock-openai"}, nil
	case "fail-first":
		return []string{"mock-openai-fail", "mock-gemini"}, nil
	default:
		return nil, ErrUnsupportedModel
	}
}

func callProvider(req models.CreateChatCompletionRequest, provider string, providerURL string) (models.CreateChatCompletionResponse, error) {
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
		return models.CreateChatCompletionResponse{}, fmt.Errorf("provider service returned status %v", resp.StatusCode)
	}

	var result models.CreateChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.CreateChatCompletionResponse{}, errors.New("internal server error")
	}

	return result, nil
}

func RouteChat(req models.CreateChatCompletionRequest, providerURL string) (models.CreateChatCompletionResponse, error) {
	providers, err := selectProviders(req.Model)
	if err != nil {
		return models.CreateChatCompletionResponse{}, err
	}
	for _, provider := range providers {
		for attempt := 1; attempt <= maxProviderAttempts; attempt++ {
			resp, err := callProvider(req, provider, providerURL)
			if err == nil {
				return resp, nil
			}
			log.Printf("provider attempt failed request_id=%v provider=%v attempt=%v error=%v", req.RequestID, provider, attempt, err.Error())
		}
	}
	return models.CreateChatCompletionResponse{}, errors.New("all providers failed")
}
