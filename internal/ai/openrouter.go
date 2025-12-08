package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-goll/aigit/internal/config"
)

type OpenRouterClient struct {
	apiKey  string
	model   string
	baseURL string
}

func NewOpenRouterClient(cfg *config.Config) (*OpenRouterClient, error) {
	model := cfg.Model
	if model == "" {
		model = config.GetDefaultModel(config.ProviderOpenRouter)
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = config.GetDefaultBaseURL(config.ProviderOpenRouter)
	}

	return &OpenRouterClient{
		apiKey:  cfg.APIKey,
		model:   model,
		baseURL: baseURL,
	}, nil
}

type openRouterRequest struct {
	Model    string              `json:"model"`
	Messages []openRouterMessage `json:"messages"`
}

type openRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *OpenRouterClient) call(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	reqBody := openRouterRequest{
		Model: c.model,
		Messages: []openRouterMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("HTTP-Referer", "https://github.com/go-goll/aigit")
	req.Header.Set("X-Title", "aigit")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result openRouterResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.Error != nil {
		return "", fmt.Errorf("OpenRouter API error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenRouter")
	}

	return result.Choices[0].Message.Content, nil
}

func (c *OpenRouterClient) GenerateCommitMessage(ctx context.Context, diff, language string) (string, error) {
	return c.call(ctx, getCommitPrompt(language), diff)
}

func (c *OpenRouterClient) ReviewCode(ctx context.Context, diff, language string) (string, error) {
	return c.call(ctx, getReviewPrompt(language), diff)
}
