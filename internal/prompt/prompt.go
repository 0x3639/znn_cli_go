// Package prompt provides utilities for secure password input and user confirmations.
package prompt

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"golang.org/x/term"
)

// Password prompts the user for a password with echo disabled.
// The prompt message is displayed, and user input is hidden.
func Password(message string) (string, error) {
	fmt.Print(message)

	// Get file descriptor for stdin
	fd := int(os.Stdin.Fd())

	// Save terminal state and disable echo
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", fmt.Errorf("failed to set terminal mode: %w", err)
	}
	defer func() { _ = term.Restore(fd, oldState) }()

	// Read password
	password, err := term.ReadPassword(fd)
	fmt.Println() // Print newline after password input

	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	return string(password), nil
}

// PasswordWithConfirm prompts for a password twice and ensures they match.
func PasswordWithConfirm(message string) (string, error) {
	pass1, err := Password(message)
	if err != nil {
		return "", err
	}

	pass2, err := Password("Confirm passphrase: ")
	if err != nil {
		return "", err
	}

	if pass1 != pass2 {
		return "", fmt.Errorf("passphrases do not match")
	}

	return pass1, nil
}

// Confirm prompts the user for a yes/no confirmation.
// Returns true if the user confirms, false otherwise.
func Confirm(message string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		// User declined or error occurred
		if err == promptui.ErrAbort {
			return false, nil
		}
		return false, err
	}

	// Check for affirmative response
	result = strings.ToLower(strings.TrimSpace(result))
	return result == "y" || result == "yes", nil
}

// Select prompts the user to select from a list of options.
// Returns the index of the selected item.
func Select(message string, items []string) (int, string, error) {
	prompt := promptui.Select{
		Label: message,
		Items: items,
	}

	index, result, err := prompt.Run()
	if err != nil {
		return -1, "", err
	}

	return index, result, nil
}

// Input prompts the user for text input with validation.
func Input(message string, validate promptui.ValidateFunc) (string, error) {
	prompt := promptui.Prompt{
		Label:    message,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}
