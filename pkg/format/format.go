// Package format provides formatting utilities for amounts, durations, and output display.
// It handles conversion between base units and human-readable formats, with proper
// color coding for different token types.
package format

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	// CoinDecimals is the number of decimal places for ZNN and QSR
	CoinDecimals = 8

	// OneZnn represents 1 ZNN in base units
	OneZnn = 100000000

	// OneQsr represents 1 QSR in base units
	OneQsr = 100000000
)

var (
	// Color functions for terminal output
	Red     = color.New(color.FgRed).SprintFunc()
	Green   = color.New(color.FgGreen).SprintFunc()
	Blue    = color.New(color.FgBlue).SprintFunc()
	Magenta = color.New(color.FgMagenta).SprintFunc()
	Yellow  = color.New(color.FgYellow).SprintFunc()
	Cyan    = color.New(color.FgCyan).SprintFunc()

	// verboseMode controls whether verbose messages are displayed
	verboseMode = false
)

// SetVerbose sets the verbose mode
func SetVerbose(v bool) {
	verboseMode = v
}

// GetVerbose returns the current verbose mode setting
func GetVerbose() bool {
	return verboseMode
}

// Amount formats a token amount from base units to human-readable format.
// Example: 100000000 with decimals=8 returns "1.00000000"
func Amount(amount *big.Int, decimals int) string {
	if amount == nil {
		return "0." + strings.Repeat("0", decimals)
	}

	// Convert to string with proper decimal places
	str := amount.String()
	negative := strings.HasPrefix(str, "-")
	if negative {
		str = str[1:]
	}

	// Pad with zeros if necessary
	for len(str) <= decimals {
		str = "0" + str
	}

	// Insert decimal point
	dotPos := len(str) - decimals
	result := str[:dotPos] + "." + str[dotPos:]

	if negative {
		result = "-" + result
	}

	return result
}

// ParseAmount parses a human-readable amount string to base units.
// Example: "1.5" with decimals=8 returns 150000000
func ParseAmount(amountStr string, decimals int) (*big.Int, error) {
	// Handle empty string
	amountStr = strings.TrimSpace(amountStr)
	if amountStr == "" {
		return nil, fmt.Errorf("amount cannot be empty")
	}

	// Parse the float string
	parts := strings.Split(amountStr, ".")
	if len(parts) > 2 {
		return nil, fmt.Errorf("invalid amount format: %s", amountStr)
	}

	integerPart := parts[0]
	decimalPart := ""
	if len(parts) == 2 {
		decimalPart = parts[1]
	}

	// Validate integer part
	if integerPart == "" {
		integerPart = "0"
	}

	// Pad or truncate decimal part to match decimals
	if len(decimalPart) > decimals {
		return nil, fmt.Errorf("too many decimal places: max %d", decimals)
	}
	decimalPart = decimalPart + strings.Repeat("0", decimals-len(decimalPart))

	// Combine and parse as big.Int
	combined := integerPart + decimalPart
	amount := new(big.Int)
	if _, ok := amount.SetString(combined, 10); !ok {
		return nil, fmt.Errorf("invalid amount: %s", amountStr)
	}

	return amount, nil
}

// Duration formats a duration in seconds to HH:MM:SS format
func Duration(seconds int64) string {
	duration := time.Duration(seconds) * time.Second
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	secs := int(duration.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}

// FormatZNN formats ZNN amount with green color
func FormatZNN(amount *big.Int) string {
	return Green(Amount(amount, CoinDecimals) + " ZNN")
}

// FormatQSR formats QSR amount with blue color
func FormatQSR(amount *big.Int) string {
	return Blue(Amount(amount, CoinDecimals) + " QSR")
}

// FormatToken formats a token amount with appropriate color and symbol
func FormatToken(amount *big.Int, decimals int, symbol string) string {
	formatted := Amount(amount, decimals) + " " + symbol
	switch strings.ToUpper(symbol) {
	case "ZNN":
		return Green(formatted)
	case "QSR":
		return Blue(formatted)
	default:
		return Magenta(formatted)
	}
}

// ParseTokenStandard parses a token identifier (ZNN, QSR, or ZTS address)
// Returns the normalized token identifier
func ParseTokenStandard(token string) (string, error) {
	token = strings.TrimSpace(token)
	upper := strings.ToUpper(token)

	switch upper {
	case "ZNN":
		return "zts1znnxxxxxxxxxxxxx9z4ulx", nil // ZNN token standard
	case "QSR":
		return "zts1qsrxxxxxxxxxxxxxmrhjll", nil // QSR token standard
	default:
		// Validate ZTS format: starts with "zts1", 40 characters
		if !strings.HasPrefix(token, "zts1") {
			return "", fmt.Errorf("invalid token standard: must be ZNN, QSR, or start with zts1")
		}
		if len(token) != 40 {
			return "", fmt.Errorf("invalid token standard: must be 40 characters")
		}
		return token, nil
	}
}

// ValidateAddress validates a Zenon address format
func ValidateAddress(addr string) error {
	addr = strings.TrimSpace(addr)

	if !strings.HasPrefix(addr, "z1") {
		return fmt.Errorf("invalid address: must start with z1")
	}

	if len(addr) != 40 {
		return fmt.Errorf("invalid address: must be 40 characters")
	}

	return nil
}

// Success prints a success message in green
func Success(message string) {
	fmt.Println(Green("✓ " + message))
}

// Error prints an error message in red
func Error(message string) {
	fmt.Println(Red("✗ Error! " + message))
}

// Warning prints a warning message in yellow
func Warning(message string) {
	fmt.Println(Yellow("⚠ " + message))
}

// Info prints an info message in cyan
func Info(message string) {
	fmt.Println(Cyan("ℹ " + message))
}
