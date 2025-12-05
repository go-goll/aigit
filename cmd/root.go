package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aigit",
	Short: "AI-powered git commit message generator",
	Long: `aigit is a CLI tool that uses AI to generate meaningful git commit messages
and review code changes for potential bugs.

Supported AI providers: OpenAI, Claude, Google Gemini`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(reviewCmd)
}
