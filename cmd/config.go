package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-goll/aigit/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config [key] [value]",
	Short: "Configure aigit settings",
	Long: `Configure AI provider, API key, model, and language settings.

Usage:
  aigit config              # Interactive setup (full configuration)
  aigit config --show       # Show current configuration
  aigit config <key> <value> # Set a specific config value

Available keys:
  provider   - AI provider (openai, claude, google)
  api_key    - API key for the provider
  model      - Model name
  language   - Output language (en, zh)
  base_url   - Custom API base URL`,
	RunE: runConfig,
}

var showConfig bool

func init() {
	configCmd.Flags().BoolVar(&showConfig, "show", false, "Show current configuration")
}

func runConfig(cmd *cobra.Command, args []string) error {
	if showConfig {
		return showCurrentConfig()
	}

	if len(args) >= 2 {
		return setConfigValue(args[0], args[1])
	}

	if len(args) == 1 {
		return fmt.Errorf("missing value for key '%s'", args[0])
	}

	return runInteractiveConfig()
}

func showCurrentConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Println("=== Current Configuration ===")
	fmt.Printf("provider:  %s\n", cfg.Provider)
	fmt.Printf("api_key:   %s\n", maskAPIKey(cfg.APIKey))
	fmt.Printf("model:     %s\n", cfg.Model)
	fmt.Printf("language:  %s\n", cfg.Language)
	if cfg.BaseURL != "" {
		fmt.Printf("base_url:  %s\n", cfg.BaseURL)
	}
	return nil
}

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

func setConfigValue(key, value string) error {
	cfg, err := config.Load()
	if err != nil {
		cfg = config.DefaultConfig()
	}

	switch key {
	case "provider":
		switch value {
		case "openai":
			cfg.Provider = config.ProviderOpenAI
		case "claude":
			cfg.Provider = config.ProviderClaude
		case "google":
			cfg.Provider = config.ProviderGoogle
		default:
			return fmt.Errorf("invalid provider: %s (use: openai, claude, google)", value)
		}
	case "api_key":
		cfg.APIKey = value
	case "model":
		cfg.Model = value
	case "language":
		if value != "en" && value != "zh" {
			return fmt.Errorf("invalid language: %s (use: en, zh)", value)
		}
		cfg.Language = value
	case "base_url":
		cfg.BaseURL = value
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✓ Set %s = %s\n", key, value)
	return nil
}

func runInteractiveConfig() error {
	reader := bufio.NewReader(os.Stdin)

	existingCfg, _ := config.Load()
	cfg := config.DefaultConfig()
	if existingCfg != nil {
		cfg = existingCfg
	}

	fmt.Println("=== aigit Configuration ===")
	fmt.Println()

	fmt.Println("Select AI provider:")
	fmt.Println("  1. OpenAI (GPT-4)")
	fmt.Println("  2. Claude (Anthropic)")
	fmt.Println("  3. Google (Gemini)")

	currentProvider := "1"
	switch cfg.Provider {
	case config.ProviderClaude:
		currentProvider = "2"
	case config.ProviderGoogle:
		currentProvider = "3"
	}
	fmt.Printf("Enter choice [%s]: ", currentProvider)

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	if choice == "" {
		choice = currentProvider
	}

	switch choice {
	case "2":
		cfg.Provider = config.ProviderClaude
	case "3":
		cfg.Provider = config.ProviderGoogle
	default:
		cfg.Provider = config.ProviderOpenAI
	}

	currentKey := maskAPIKey(cfg.APIKey)
	fmt.Printf("\nEnter API key for %s [%s]: ", cfg.Provider, currentKey)
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)
	if apiKey != "" {
		cfg.APIKey = apiKey
	}

	if cfg.APIKey == "" {
		return fmt.Errorf("API key is required")
	}

	defaultModel := cfg.Model
	if defaultModel == "" {
		defaultModel = config.GetDefaultModel(cfg.Provider)
	}
	fmt.Printf("\nEnter model name [%s]: ", defaultModel)
	model, _ := reader.ReadString('\n')
	model = strings.TrimSpace(model)
	if model != "" {
		cfg.Model = model
	} else if cfg.Model == "" {
		cfg.Model = defaultModel
	}

	fmt.Println("\nSelect language for commit messages:")
	fmt.Println("  1. English")
	fmt.Println("  2. Chinese (中文)")

	currentLang := "1"
	if cfg.Language == "zh" {
		currentLang = "2"
	}
	fmt.Printf("Enter choice [%s]: ", currentLang)

	langChoice, _ := reader.ReadString('\n')
	langChoice = strings.TrimSpace(langChoice)
	if langChoice == "" {
		langChoice = currentLang
	}

	if langChoice == "2" {
		cfg.Language = "zh"
	} else {
		cfg.Language = "en"
	}

	currentURL := cfg.BaseURL
	if currentURL == "" {
		currentURL = "default"
	}
	fmt.Printf("\nEnter custom base URL [%s]: ", currentURL)
	baseURL, _ := reader.ReadString('\n')
	baseURL = strings.TrimSpace(baseURL)
	if baseURL != "" && baseURL != "default" {
		cfg.BaseURL = baseURL
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("\n✓ Configuration saved successfully!")
	return nil
}
