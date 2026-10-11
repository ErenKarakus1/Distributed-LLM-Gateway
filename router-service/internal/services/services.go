package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/router-service/internal/models"
)

const providerURL = "http://localhost:8082/providers/chat"

func RouteChat(req models.CreateChatCompletionRequest) (models.CreateChatCompletionResponse, error) {
	body, err := json.Marshal(req)
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
