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

type ClaudeClient struct {
	apiKey  string
	model   string
	baseURL string
}

func NewClaudeClient(cfg *config.Config) (*ClaudeClient, error) {
	model := cfg.Model
	if model == "" {
		model = config.GetDefaultModel(config.ProviderClaude)
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1"
	}

	return &ClaudeClient{
		apiKey:  cfg.APIKey,
		model:   model,
		baseURL: baseURL,
	}, nil
}

type claudeRequest struct {
	Model     string           `json:"model"`
	MaxTokens int              `json:"max_tokens"`
	System    string           `json:"system,omitempty"`
	Messages  []claudeMessage  `json:"messages"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *ClaudeClient) call(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	reqBody := claudeRequest{
		Model:     c.model,
		MaxTokens: 4096,
		System:    systemPrompt,
		Messages: []claudeMessage{
			{Role: "user", Content: userPrompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result claudeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.Error != nil {
		return "", fmt.Errorf("Claude API error: %s", result.Error.Message)
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("no response from Claude")
	}

	return result.Content[0].Text, nil
}

func (c *ClaudeClient) GenerateCommitMessage(ctx context.Context, diff, language string) (string, error) {
	return c.call(ctx, getCommitPrompt(language), diff)
}

func (c *ClaudeClient) ReviewCode(ctx context.Context, diff, language string) (string, error) {
	return c.call(ctx, getReviewPrompt(language), diff)
}
