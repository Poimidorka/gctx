package cmd

import (
	"fmt"
	"os"
)

// Primary logic unit

// appendGctxMetadata adds the section for identifying the profile
func appendGctxMetadata(content FileContent, profile string) FileContent {
	metadata := fmt.Sprintf(`
[gctx]
	profile = %s
`, profile)

	out := make([]byte, 0, len(content)+len(metadata)+1)
	out = append(out, content...)

	if len(out) > 0 && out[len(out)-1] != '\n' {
		out = append(out, '\n')
	}

	out = append(out, []byte(metadata)...)
	return out
}

// saveCurrentGitProfile puts current gitconfig into the profile
func saveCurrentGitProfile(profileStore *ProfileStore, name string) error {
	content, err := readCurrentGitConfig()
	if err != nil {
		return err
	}
	err = profileStore.Set(name, content)
	if err != nil {
		return err
	}
	return nil
}

// TODO: add support for global flag
func applyGitProfile(profileStore *ProfileStore, profile string, global bool) error {
	content, err := profileStore.Get(profile)
	if err != nil {
		return err
	}
	path, err := getGitConfigPath()
	if err != nil {
		return err
	}
	err = os.WriteFile(path, *content, 0o600)
	if err != nil {
		return fmt.Errorf("could not write the config %s: %w", path, err)
	}
	err = setGtxProfile(profile)
	if err != nil {
		return err
	}
	return nil
}
