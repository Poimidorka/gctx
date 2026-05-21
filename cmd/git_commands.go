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

//TODO: add support for the global flag

// getGitConfigPath returns the config path for the active .git directory.
func getGitConfigPath() (string, error) {
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

func readCurrentGitConfig() (*FileContent, error) {
	path, err := getGitConfigPath()
	if err != nil {
		return nil, err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", path, err)
	}
	return (*FileContent)(&content), nil
}

func getCurrentGtxProfile() (string, error) {
	out, err := exec.Command("git", "config", "--local", "--get", "gctx.profile").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("couldn't get current gctx profile: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

func setGtxProfile(profile string) error {
	out, err := exec.Command("git", "config", "--local", "--replace-all", "gctx.profile", profile).CombinedOutput()
	if err != nil {
		return fmt.Errorf("couldn't set current gctx profile: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}
