package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Provider string

const (
	ProviderOpenAI  Provider = "openai"
	ProviderClaude  Provider = "claude"
	ProviderGoogle  Provider = "google"
)

type Config struct {
	Provider Provider `json:"provider"`
	APIKey   string   `json:"api_key"`
	Model    string   `json:"model"`
	Language string   `json:"language"` // "en" or "zh"
	BaseURL  string   `json:"base_url,omitempty"`
}

func DefaultConfig() *Config {
	return &Config{
		Provider: ProviderOpenAI,
		Language: "en",
	}
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".aigit", "config.json"), nil
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("config not found, please run 'aigit config' first")
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.APIKey == "" {
		return nil, errors.New("api_key is required")
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func GetDefaultModel(provider Provider) string {
	switch provider {
	case ProviderOpenAI:
		return "gpt-4o"
	case ProviderClaude:
		return "claude-sonnet-4-20250514"
	case ProviderGoogle:
		return "gemini-1.5-pro"
	default:
		return ""
	}
}
