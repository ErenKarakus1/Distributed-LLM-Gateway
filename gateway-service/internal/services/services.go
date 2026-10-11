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

func CreateChatCompletion(req models.CreateChatCompletionRequest, routerURL string) (models.CreateChatCompletionResponse, error) {
	routerReq := models.RouterRequest{RequestID: uuid.New(), CreateChatCompletionRequest: req}
	log.Printf("received chat completion request request_id=%v model=%v", routerReq.RequestID, routerReq.Model)
	body, err := json.Marshal(routerReq)
	if err != nil {
		return models.CreateChatCompletionResponse{}, models.UpstreamError{StatusCode: 500, Message: "internal server error"}
	}
	client := http.Client{
		Timeout: 5 * time.Second,
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
