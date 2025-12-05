package git

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	return cmd.Run() == nil
}

func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func GetUnstagedDiff() (string, error) {
	cmd := exec.Command("git", "diff")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func GetAllDiff() (string, error) {
	staged, err := GetStagedDiff()
	if err != nil {
		return "", err
	}

	unstaged, err := GetUnstagedDiff()
	if err != nil {
		return "", err
	}

	if staged == "" && unstaged == "" {
		return "", errors.New("no changes to commit")
	}

	var result strings.Builder
	if staged != "" {
		result.WriteString("=== Staged Changes ===\n")
		result.WriteString(staged)
	}
	if unstaged != "" {
		if staged != "" {
			result.WriteString("\n")
		}
		result.WriteString("=== Unstaged Changes ===\n")
		result.WriteString(unstaged)
	}

	return result.String(), nil
}

func GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(files) == 1 && files[0] == "" {
		return nil, nil
	}
	return files, nil
}

func Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	return cmd.Run()
}

func StageAll() error {
	cmd := exec.Command("git", "add", "-A")
	return cmd.Run()
}
