package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type OpenAIClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		apiKey: apiKey,
		// タイムアウトを設定し、OpenAI が応答しないときにゴルーチンと接続が
		// 無限に滞留するのを防ぐ。
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type embeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type embeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}

func (c *OpenAIClient) GenerateEmbedding(ctx context.Context, text string) ([]float64, error) {
	reqBody := embeddingRequest{
		Model: "text-embedding-3-small",
		Input: text,
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 429 || resp.StatusCode == 402 {
			return nil, fmt.Errorf("OpenAI APIの利用上限に達しています。管理者に連絡してください")
		}
		return nil, fmt.Errorf("AI処理でエラーが発生しました。しばらく経ってから再度お試しください")
	}

	var result embeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no embedding data returned")
	}

	return result.Data[0].Embedding, nil
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type ChatMessage = chatMessage

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

func (c *OpenAIClient) ChatCompletion(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	return c.ChatCompletionWithHistory(ctx, systemPrompt, []chatMessage{{Role: "user", Content: userMessage}})
}

func (c *OpenAIClient) ChatCompletionWithHistory(ctx context.Context, systemPrompt string, history []chatMessage) (string, error) {
	messages := []chatMessage{{Role: "system", Content: systemPrompt}}
	messages = append(messages, history...)
	reqBody := chatRequest{
		Model:    "gpt-4o-mini",
		Messages: messages,
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 429 || resp.StatusCode == 402 {
			return "", fmt.Errorf("OpenAI APIの利用上限に達しています。管理者に連絡してください")
		}
		return "", fmt.Errorf("AI処理でエラーが発生しました。しばらく経ってから再度お試しください")
	}

	var result chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no chat completion returned")
	}

	return result.Choices[0].Message.Content, nil
}
