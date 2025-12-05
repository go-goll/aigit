package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-goll/aigit/internal/ai"
	"github.com/go-goll/aigit/internal/config"
	"github.com/go-goll/aigit/internal/git"
	"github.com/spf13/cobra"
)

var (
	autoCommit bool
	stageAll   bool
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate a commit message using AI",
	Long:  `Analyze staged changes and generate a meaningful commit message using AI.`,
	RunE:  runCommit,
}

func init() {
	commitCmd.Flags().BoolVarP(&autoCommit, "yes", "y", false, "Auto commit without confirmation")
	commitCmd.Flags().BoolVarP(&stageAll, "all", "a", false, "Stage all changes before commit")
}

func runCommit(cmd *cobra.Command, args []string) error {
	if !git.IsGitRepo() {
		return fmt.Errorf("not a git repository")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if stageAll {
		if err := git.StageAll(); err != nil {
			return fmt.Errorf("failed to stage changes: %w", err)
		}
		fmt.Println("✓ Staged all changes")
	}

	diff, err := git.GetStagedDiff()
	if err != nil {
		return fmt.Errorf("failed to get diff: %w", err)
	}

	if diff == "" {
		return fmt.Errorf("no staged changes to commit")
	}

	files, _ := git.GetStagedFiles()
	if len(files) > 0 {
		fmt.Println("Staged files:")
		for _, f := range files {
			fmt.Printf("  • %s\n", f)
		}
		fmt.Println()
	}

	fmt.Println("Generating commit message...")

	client, err := ai.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create AI client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	message, err := client.GenerateCommitMessage(ctx, diff, cfg.Language)
	if err != nil {
		return fmt.Errorf("failed to generate commit message: %w", err)
	}

	message = strings.TrimSpace(message)

	fmt.Println("\n--- Generated Commit Message ---")
	fmt.Println(message)
	fmt.Println("--------------------------------")

	if autoCommit {
		return doCommit(message)
	}

	fmt.Print("\nCommit with this message? [Y/n/e(dit)]: ")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.ToLower(strings.TrimSpace(answer))

	switch answer {
	case "", "y", "yes":
		return doCommit(message)
	case "e", "edit":
		fmt.Print("Enter new message: ")
		newMsg, _ := reader.ReadString('\n')
		newMsg = strings.TrimSpace(newMsg)
		if newMsg != "" {
			return doCommit(newMsg)
		}
		fmt.Println("Empty message, commit aborted.")
		return nil
	default:
		fmt.Println("Commit aborted.")
		return nil
	}
}

func doCommit(message string) error {
	if err := git.Commit(message); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	fmt.Println("✓ Committed successfully!")
	return nil
}
