package gpt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	url = "https://api.openai.com/v1/chat/completions"
)

type Client interface {
	SendRequest(ctx context.Context, apiKey, prompt string) (string, error)
}

type client struct {
}

func NewClient() Client {
	return &client{}
}

func (c *client) SendRequest(ctx context.Context, apiKey, prompt string) (string, error) {
	gptReq := NewRequest(prompt)
	jsonData, err := json.Marshal(gptReq)
	if err != nil {
		return "", err
	}
	// Создаем новый POST запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка создания запроса:", err)
		return "", err
	}

	// Устанавливаем необходимые заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка выполнения запроса:", err)
		return "", err
	}
	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}
	return response.Choices[0].Message.Content, nil
}
