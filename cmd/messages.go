package cmd

import "fmt"

const NoActiveContextMessage = "No active context."

func CurrentContextMessage(profile string) string {
	return fmt.Sprintf("Current context: %q.", profile)
}

func SwitchedContextMessage(profile string) string {
	return fmt.Sprintf("✔ Switched to context %q.", profile)
}

func SavedContextMessage(profile string) string {
	return fmt.Sprintf("✔ Saved context %q.", profile)
}

func RemovedContextMessage(profile string) string {
	return fmt.Sprintf("✔ Removed context %q.", profile)
}

func ProfileNameRequiredMessage() string {
	return "context name is required with --save or --remove"
}

func ConflictingActionMessage() string {
	return "use either --save or --remove, not both"
}
