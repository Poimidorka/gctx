package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//TODO: consider go-git as an alternative of running commands manually
// dependency weights 10MB

// getGitConfigPath returns the git config path for the requested scope.
func getGitConfigPath(global bool) (string, error) {
	if global {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".gitconfig"), nil
	}

	output, err := exec.Command(`git`, `rev-parse`, `--absolute-git-dir`).CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(output))
		if strings.Contains(message, `not a git repository`) {
			return "", fmt.Errorf(`couldn't find the git repo`)
		}
		return "", fmt.Errorf("%w: %s", err, message)
	}
	return filepath.Join(strings.TrimSpace(string(output)), "config"), nil
}

func readCurrentGitConfig(global bool) (*FileContent, error) {
	path, err := getGitConfigPath(global)
	if err != nil {
		return nil, err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", path, err)
	}
	return (*FileContent)(&content), nil
}

func getCurrentGtxProfile(global bool) (string, error) {
	path, err := getGitConfigPath(global)
	if err != nil {
		return "", err
	}
	out, err := exec.Command("git", "config", "-f", path, "--get", "gctx.profile").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("couldn't get current gctx profile: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

func setGtxProfile(profile string, global bool) error {
	path, err := getGitConfigPath(global)
	if err != nil {
		return err
	}
	out, err := exec.Command("git", "config", "-f", path, "--replace-all", "gctx.profile", profile).CombinedOutput()
	if err != nil {
		return fmt.Errorf("couldn't set current gctx profile: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}
