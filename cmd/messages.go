package cmd

import (
	"fmt"
	"strings"
)

const (
	green = "\x1b[32m"
	red   = "\x1b[31m"
	reset = "\x1b[0m"

	LongDescription = `gctx saves named git profiles and applies them to a repository or to
your global Git config.

Profiles are snapshots of Git config files, commonly used for switching
user.name, user.email, signing keys, and other identity settings.`

	Examples = `  gctx
  gctx personal
  git config --local --add "user.name" "Alice"
  gctx personal --save
  gctx work --remove
  gctx personal --global`

	NoActiveContextMessage = "No active profile in current git config."
)

func CurrentContextMessage(profile string) string {
	return fmt.Sprintf("Current git profile: %q.", profile)
}

func SwitchedContextMessage(profile string) string {
	return fmt.Sprintf("%s Switched to profile %q.", successMark(), profile)
}

func SavedContextMessage(profile string) string {
	return fmt.Sprintf("%s Saved context to profile: %q.", successMark(), profile)
}

func RemovedContextMessage(profile string) string {
	return fmt.Sprintf("%s Removed profile %q.", successMark(), profile)
}

func MissingContextMessage(profile string, available []string) string {
	if len(available) == 0 {
		return fmt.Sprintf("%s Profile %q not found. No saved profiles.", errorMark(), profile)
	}
	return fmt.Sprintf("%s Profile %q not found. Available profiles: %s.", errorMark(), profile, strings.Join(available, " "))
}

func ProfileNameRequiredMessage() string {
	return fmt.Sprintf("%s Profile name is required with --save or --remove.", errorMark())
}

func ConflictingActionMessage() string {
	return fmt.Sprintf("%s Use either --save or --remove, not both.", errorMark())
}

func successMark() string {
	return green + "✔" + reset
}

func errorMark() string {
	return red + "✘" + reset
}
