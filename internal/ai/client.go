package ai

import (
	"context"

	"github.com/go-goll/aigit/internal/config"
)

type Client interface {
	GenerateCommitMessage(ctx context.Context, diff, language string) (string, error)
	ReviewCode(ctx context.Context, diff, language string) (string, error)
}

func NewClient(cfg *config.Config) (Client, error) {
	switch cfg.Provider {
	case config.ProviderOpenAI:
		return NewOpenAIClient(cfg)
	case config.ProviderClaude:
		return NewClaudeClient(cfg)
	case config.ProviderGoogle:
		return NewGoogleClient(cfg)
	case config.ProviderOpenRouter:
		return NewOpenRouterClient(cfg)
	default:
		return NewOpenAIClient(cfg)
	}
}
