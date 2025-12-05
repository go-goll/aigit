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

type GoogleClient struct {
	apiKey  string
	model   string
	baseURL string
}

func NewGoogleClient(cfg *config.Config) (*GoogleClient, error) {
	model := cfg.Model
	if model == "" {
		model = config.GetDefaultModel(config.ProviderGoogle)
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
	}

	return &GoogleClient{
		apiKey:  cfg.APIKey,
		model:   model,
		baseURL: baseURL,
	}, nil
}

type googleRequest struct {
	Contents         []googleContent       `json:"contents"`
	SystemInstruction *googleContent       `json:"systemInstruction,omitempty"`
}

type googleContent struct {
	Parts []googlePart `json:"parts"`
	Role  string       `json:"role,omitempty"`
}

type googlePart struct {
	Text string `json:"text"`
}

type googleResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *GoogleClient) call(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	reqBody := googleRequest{
		SystemInstruction: &googleContent{
			Parts: []googlePart{{Text: systemPrompt}},
		},
		Contents: []googleContent{
			{
				Role:  "user",
				Parts: []googlePart{{Text: userPrompt}},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", c.baseURL, c.model, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result googleResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.Error != nil {
		return "", fmt.Errorf("Google API error: %s", result.Error.Message)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Google")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}

func (c *GoogleClient) GenerateCommitMessage(ctx context.Context, diff, language string) (string, error) {
	return c.call(ctx, getCommitPrompt(language), diff)
}

func (c *GoogleClient) ReviewCode(ctx context.Context, diff, language string) (string, error) {
	return c.call(ctx, getReviewPrompt(language), diff)
}
