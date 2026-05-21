package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ProfileStore struct {
	// stores paths to the saved profiles for git
	configDir string
}

type FileContent []byte

func hasExt(filename string) bool {
	return strings.HasSuffix(filename, ".config")
}

func addExt(profile string) string {
	if hasExt(profile) {
		return profile
	}
	return profile + `.config`
}

func removeExt(fileName string) string {
	return strings.TrimSuffix(fileName, `.config`)
}

func (s *ProfileStore) Set(name string, content *FileContent) error {
	err := os.WriteFile(filepath.Join(s.configDir, addExt(name)), *content, 0o600)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProfileStore) Remove(name string) error {
	err := os.Remove(filepath.Join(s.configDir, addExt(name)))
	if err != nil {
		return err
	}
	return nil
}

// List returns array of profiles
func (s *ProfileStore) List() []string {
	res, err := os.ReadDir(s.configDir)
	if err != nil {
		return nil
	}
	profiles := make([]string, 0, len(res))

	for _, entry := range res {
		name := entry.Name()
		if !entry.IsDir() && hasExt(name) {
			profiles = append(profiles, removeExt(name))
		}
	}

	return profiles
}

// Get returns FileContent for a given profile
func (s *ProfileStore) Get(profile string) (*FileContent, error) {
	content, err := os.ReadFile(filepath.Join(s.configDir, addExt(profile)))
	if err != nil {
		return nil, fmt.Errorf("failed to read profile %s: %w", profile, err)
	}
	return (*FileContent)(&content), nil
}

func NewProfileStore(configDir string) *ProfileStore {
	return &ProfileStore{configDir: configDir}
}
