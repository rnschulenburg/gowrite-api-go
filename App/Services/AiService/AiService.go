package AiService

import (
	"encoding/json"
	"github.com/rnschulenburg/gowrite-api-go/App/Entities/AiChatEntity"
	"io"
	"log"
	"net/http"
	"os"
)

var AiSecret = ""

func InitAi() {
	secret := os.Getenv("AiSecret")

	if secret == "" {
		log.Fatal("JwtSecret not set")
	}

	AiSecret = secret
}

func GoChat(r *http.Request) (AiChatEntity.AiChatResponse, error) {

	req, err := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/chat/completions",
		r.Body,
	)
	if err != nil {
		return AiChatEntity.AiChatResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+AiSecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AiChatEntity.AiChatResponse{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return AiChatEntity.AiChatResponse{}, err
	}

	var aiResp AiChatEntity.AiChatResponse
	err = json.Unmarshal(bodyBytes, &aiResp)
	if err != nil {
		return AiChatEntity.AiChatResponse{}, err
	}

	return aiResp, nil
}
