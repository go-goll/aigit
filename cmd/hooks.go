package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var hooksCmd = &cobra.Command{
	Use:   "hooks",
	Short: "Manage git hooks",
	Long:  `Install or uninstall git hooks for automatic code review on commit.`,
}

var installHooksCmd = &cobra.Command{
	Use:   "install",
	Short: "Install pre-commit hook for auto review",
	RunE:  runInstallHooks,
}

var uninstallHooksCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall pre-commit hook",
	RunE:  runUninstallHooks,
}

func init() {
	hooksCmd.AddCommand(installHooksCmd)
	hooksCmd.AddCommand(uninstallHooksCmd)
	rootCmd.AddCommand(hooksCmd)
}

const preCommitHook = `#!/bin/sh
# aigit pre-commit hook - auto review code before commit

echo "Running aigit code review..."
aigit review --staged --hook

if [ $? -ne 0 ]; then
    echo ""
    echo "Code review found issues. Commit aborted."
    echo "Use 'git commit --no-verify' to skip this check."
    exit 1
fi
`

func getHooksDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository")
	}
	gitDir := string(out[:len(out)-1])
	return filepath.Join(gitDir, "hooks"), nil
}

func runInstallHooks(cmd *cobra.Command, args []string) error {
	hooksDir, err := getHooksDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	hookPath := filepath.Join(hooksDir, "pre-commit")

	if _, err := os.Stat(hookPath); err == nil {
		backupPath := hookPath + ".backup"
		if err := os.Rename(hookPath, backupPath); err != nil {
			return fmt.Errorf("failed to backup existing hook: %w", err)
		}
		fmt.Printf("Existing pre-commit hook backed up to: %s\n", backupPath)
	}

	if err := os.WriteFile(hookPath, []byte(preCommitHook), 0755); err != nil {
		return fmt.Errorf("failed to write hook: %w", err)
	}

	fmt.Println("✓ Pre-commit hook installed successfully!")
	fmt.Println("  Code will be reviewed automatically before each commit.")
	fmt.Println("  Use 'git commit --no-verify' to skip the review.")
	return nil
}

func runUninstallHooks(cmd *cobra.Command, args []string) error {
	hooksDir, err := getHooksDir()
	if err != nil {
		return err
	}

	hookPath := filepath.Join(hooksDir, "pre-commit")

	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		fmt.Println("No pre-commit hook found.")
		return nil
	}

	content, err := os.ReadFile(hookPath)
	if err != nil {
		return fmt.Errorf("failed to read hook: %w", err)
	}

	if string(content) != preCommitHook {
		return fmt.Errorf("pre-commit hook was not installed by aigit, refusing to remove")
	}

	if err := os.Remove(hookPath); err != nil {
		return fmt.Errorf("failed to remove hook: %w", err)
	}

	backupPath := hookPath + ".backup"
	if _, err := os.Stat(backupPath); err == nil {
		if err := os.Rename(backupPath, hookPath); err != nil {
			fmt.Printf("Warning: failed to restore backup hook: %v\n", err)
		} else {
			fmt.Println("Restored previous pre-commit hook from backup.")
		}
	}

	fmt.Println("✓ Pre-commit hook uninstalled successfully!")
	return nil
}
