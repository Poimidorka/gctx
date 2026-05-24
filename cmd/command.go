package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	ConfigGlobalEnv = `GIT_CONFIG_GLOBAL`
)

func runProfileCommand(out io.Writer, profileStore *ProfileStore, profile string, command string, interactive bool) error {
	exists, profiles := profileStore.Contains(profile)
	if !exists {
		return fmt.Errorf("%s", MissingContextMessage(profile, profiles))
	}

	content, err := profileStore.Get(profile)
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp("", "gctx-command-*.config")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if _, err := tmp.Write(*content); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := setGtxProfileInFile(tmpPath, profile); err != nil {
		return err
	}

	args := strings.Fields(command)
	if len(args) == 0 {
		return fmt.Errorf("%s", CommandRequiredMessage())
	}

	fmt.Fprintln(out, RunningCommandMessage(profile, command))

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = append(os.Environ(), ConfigGlobalEnv+"="+tmpPath)
	if interactive {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	} else {
		output, err := cmd.CombinedOutput()
		if len(output) > 0 {
			fmt.Fprint(out, string(output))
		}
		if err != nil {
			return err
		}
	}

	fmt.Fprintln(out, ExitedCommandMessage(profile))
	return nil
}
