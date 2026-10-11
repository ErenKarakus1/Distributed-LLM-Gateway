package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/models"
)

const routerURL = "http://localhost:8081/route/chat"

func CreateChatCompletion(req models.CreateChatCompletionRequest) (models.CreateChatCompletionResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return models.CreateChatCompletionResponse{}, errors.New("internal server error")
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	httpReq, err := http.NewRequest(http.MethodPost, routerURL, bytes.NewBuffer(body))
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
