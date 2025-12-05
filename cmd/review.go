package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/go-goll/aigit/internal/ai"
	"github.com/go-goll/aigit/internal/config"
	"github.com/go-goll/aigit/internal/git"
	"github.com/spf13/cobra"
)

var reviewStaged bool

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review code changes for potential bugs",
	Long:  `Analyze code changes using AI to identify potential bugs, security issues, and code quality problems.`,
	RunE:  runReview,
}

func init() {
	reviewCmd.Flags().BoolVarP(&reviewStaged, "staged", "s", false, "Review only staged changes (default: all changes)")
}

func runReview(cmd *cobra.Command, args []string) error {
	if !git.IsGitRepo() {
		return fmt.Errorf("not a git repository")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	var diff string
	if reviewStaged {
		diff, err = git.GetStagedDiff()
		if err != nil {
			return fmt.Errorf("failed to get staged diff: %w", err)
		}
		if diff == "" {
			return fmt.Errorf("no staged changes to review")
		}
		fmt.Println("Reviewing staged changes...")
	} else {
		diff, err = git.GetAllDiff()
		if err != nil {
			return fmt.Errorf("failed to get diff: %w", err)
		}
		fmt.Println("Reviewing all changes...")
	}

	client, err := ai.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create AI client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	result, err := client.ReviewCode(ctx, diff, cfg.Language)
	if err != nil {
		return fmt.Errorf("failed to review code: %w", err)
	}

	fmt.Println("\n=== Code Review Results ===")
	fmt.Println(result)
	fmt.Println("===========================")

	return nil
}
