package cmd

import (
	"fmt"
	"strings"
)

const NoActiveContextMessage = "No active profile in current git config."

func CurrentContextMessage(profile string) string {
	return fmt.Sprintf("Current git profile: %q.", profile)
}

func SwitchedContextMessage(profile string) string {
	return fmt.Sprintf("✔ Switched to profile %q.", profile)
}

func SavedContextMessage(profile string) string {
	return fmt.Sprintf("✔ Saved context to profile: %q.", profile)
}

func RemovedContextMessage(profile string) string {
	return fmt.Sprintf("✔ Removed profile %q.", profile)
}

func MissingContextMessage(profile string, available []string) string {
	if len(available) == 0 {
		return fmt.Sprintf("Profile %q not found. No saved profiles.", profile)
	}
	return fmt.Sprintf("Profile %q not found. Available profiles: %s.", profile, strings.Join(available, " "))
}

func ProfileNameRequiredMessage() string {
	return "context name is required with --save or --remove"
}

func ConflictingActionMessage() string {
	return "use either --save or --remove, not both"
}
