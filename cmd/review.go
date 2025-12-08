package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/go-goll/aigit/internal/ai"
	"github.com/go-goll/aigit/internal/config"
	"github.com/go-goll/aigit/internal/git"
)

var (
	colorHigh   = color.New(color.FgRed, color.Bold)
	colorMedium = color.New(color.FgYellow, color.Bold)
	colorLow    = color.New(color.FgCyan)
	colorOK     = color.New(color.FgGreen, color.Bold)
)

var (
	reviewStaged bool
	hookMode     bool
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review code changes for potential bugs",
	Long:  `Analyze code changes using AI to identify potential bugs, security issues, and code quality problems.`,
	RunE:  runReview,
}

func init() {
	reviewCmd.Flags().BoolVarP(&reviewStaged, "staged", "s", false, "Review only staged changes (default: all changes)")
	reviewCmd.Flags().BoolVar(&hookMode, "hook", false, "Run in hook mode (exit with error if issues found)")
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
	printColoredResult(result)
	fmt.Println("===========================")

	return nil
}

func printColoredResult(result string) {
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		coloredLine := colorLine(line)
		fmt.Println(coloredLine)
	}
}

func colorLine(line string) string {
	upperLine := strings.ToUpper(line)

	if containsAny(upperLine, []string{"HIGH", "高严重性", "CRITICAL"}) {
		return colorHigh.Sprint(line)
	}
	if containsAny(upperLine, []string{"MEDIUM", "中严重性", "中", "WARN"}) {
		return colorMedium.Sprint(line)
	}
	if containsAny(upperLine, []string{"LOW", "低", "非严重问题", "INFO"}) {
		return colorLow.Sprint(line)
	}
	if containsAny(upperLine, []string{"NO CRITICAL", "NO ISSUE", "未发现", "NO SIGNIFICANT", "LOOKS GOOD", "NO PROBLEM"}) {
		return colorOK.Sprint(line)
	}

	return line
}

func containsAny(s string, substrs []string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func hasHighSeverityIssues(result string) bool {
	upperResult := strings.ToUpper(result)
	keywords := []string{
		"HIGH", "高",
		"CRITICAL", "严重",
		"SECURITY", "安全漏洞",
		"VULNERABILITY", "漏洞",
	}
	for _, kw := range keywords {
		if strings.Contains(upperResult, strings.ToUpper(kw)) {
			return true
		}
	}
	return false
}
