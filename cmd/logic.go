package cmd

import (
	"fmt"
	"os"
)

// Primary logic unit

// saveCurrentGitProfile puts the current git config into the profile.
func saveCurrentGitProfile(profileStore *ProfileStore, name string, global bool) error {
	content, err := readCurrentGitConfig(global)
	if err != nil {
		return err
	}
	err = profileStore.Set(name, content)
	if err != nil {
		return err
	}
	return nil
}

func applyGitProfile(profileStore *ProfileStore, profile string, global bool) error {
	content, err := profileStore.Get(profile)
	if err != nil {
		return err
	}
	path, err := getGitConfigPath(global)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, *content, 0o600)
	if err != nil {
		return fmt.Errorf("could not write the config %s: %w", path, err)
	}
	err = setGtxProfile(profile, global)
	if err != nil {
		return err
	}
	return nil
}
