package format

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAmount tests the Amount formatting function
func TestAmount(t *testing.T) {
	tests := []struct {
		name     string
		amount   *big.Int
		decimals int
		expected string
	}{
		{
			name:     "one ZNN",
			amount:   big.NewInt(100000000),
			decimals: 8,
			expected: "1.00000000",
		},
		{
			name:     "half ZNN",
			amount:   big.NewInt(50000000),
			decimals: 8,
			expected: "0.50000000",
		},
		{
			name:     "zero",
			amount:   big.NewInt(0),
			decimals: 8,
			expected: "0.00000000",
		},
		{
			name:     "nil amount",
			amount:   nil,
			decimals: 8,
			expected: "0.00000000",
		},
		{
			name:     "large amount",
			amount:   big.NewInt(123456789012345678),
			decimals: 8,
			expected: "1234567890.12345678",
		},
		{
			name:     "one unit",
			amount:   big.NewInt(1),
			decimals: 8,
			expected: "0.00000001",
		},
		{
			name:     "negative amount",
			amount:   big.NewInt(-100000000),
			decimals: 8,
			expected: "-1.00000000",
		},
		{
			name:     "two decimals",
			amount:   big.NewInt(150),
			decimals: 2,
			expected: "1.50",
		},
		{
			name:     "zero decimals",
			amount:   big.NewInt(100),
			decimals: 0,
			expected: "100.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Amount(tt.amount, tt.decimals)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseAmount tests the ParseAmount parsing function
func TestParseAmount(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		decimals  int
		expected  *big.Int
		expectErr bool
	}{
		{
			name:     "one ZNN",
			input:    "1",
			decimals: 8,
			expected: big.NewInt(100000000),
		},
		{
			name:     "one ZNN with decimals",
			input:    "1.0",
			decimals: 8,
			expected: big.NewInt(100000000),
		},
		{
			name:     "half ZNN",
			input:    "0.5",
			decimals: 8,
			expected: big.NewInt(50000000),
		},
		{
			name:     "fractional",
			input:    "1.23456789",
			decimals: 8,
			expected: big.NewInt(123456789),
		},
		{
			name:     "large amount",
			input:    "1234567890.12345678",
			decimals: 8,
			expected: big.NewInt(123456789012345678),
		},
		{
			name:     "zero",
			input:    "0",
			decimals: 8,
			expected: big.NewInt(0),
		},
		{
			name:     "zero with decimals",
			input:    "0.0",
			decimals: 8,
			expected: big.NewInt(0),
		},
		{
			name:     "leading decimal",
			input:    ".5",
			decimals: 8,
			expected: big.NewInt(50000000),
		},
		{
			name:     "with spaces",
			input:    " 10.5 ",
			decimals: 8,
			expected: big.NewInt(1050000000),
		},
		{
			name:      "empty string",
			input:     "",
			decimals:  8,
			expectErr: true,
		},
		{
			name:      "only spaces",
			input:     "   ",
			decimals:  8,
			expectErr: true,
		},
		{
			name:      "too many decimal places",
			input:     "1.123456789",
			decimals:  8,
			expectErr: true,
		},
		{
			name:      "multiple dots",
			input:     "1.2.3",
			decimals:  8,
			expectErr: true,
		},
		{
			name:      "invalid characters",
			input:     "abc",
			decimals:  8,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseAmount(tt.input, tt.decimals)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestDuration tests the Duration formatting function
func TestDuration(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int64
		expected string
	}{
		{
			name:     "zero",
			seconds:  0,
			expected: "00:00:00",
		},
		{
			name:     "one second",
			seconds:  1,
			expected: "00:00:01",
		},
		{
			name:     "one minute",
			seconds:  60,
			expected: "00:01:00",
		},
		{
			name:     "one hour",
			seconds:  3600,
			expected: "01:00:00",
		},
		{
			name:     "one day",
			seconds:  86400,
			expected: "24:00:00",
		},
		{
			name:     "complex duration",
			seconds:  3661,
			expected: "01:01:01",
		},
		{
			name:     "30 days (stake duration)",
			seconds:  30 * 24 * 60 * 60,
			expected: "720:00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Duration(tt.seconds)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseTokenStandard tests the ParseTokenStandard function
func TestParseTokenStandard(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		expectErr bool
	}{
		{
			name:     "ZNN uppercase",
			input:    "ZNN",
			expected: "zts1znnxxxxxxxxxxxxx9z4ulx",
		},
		{
			name:     "ZNN lowercase",
			input:    "znn",
			expected: "zts1znnxxxxxxxxxxxxx9z4ulx",
		},
		{
			name:     "QSR uppercase",
			input:    "QSR",
			expected: "zts1qsrxxxxxxxxxxxxxmrhjll",
		},
		{
			name:     "QSR lowercase",
			input:    "qsr",
			expected: "zts1qsrxxxxxxxxxxxxxmrhjll",
		},
		{
			name:     "valid custom ZTS (40 chars)",
			input:    "zts1customtokenxxxxxxxxxxxxxxxxx12345678",
			expected: "zts1customtokenxxxxxxxxxxxxxxxxx12345678",
		},
		{
			name:     "ZTS with spaces (40 chars)",
			input:    " zts1customtokenxxxxxxxxxxxxxxxxx12345678 ",
			expected: "zts1customtokenxxxxxxxxxxxxxxxxx12345678",
		},
		{
			name:      "invalid prefix",
			input:     "zzz1qsrxxxxxxxxxxxxxmrhjll",
			expectErr: true,
		},
		{
			name:      "too short",
			input:     "zts1qsr",
			expectErr: true,
		},
		{
			name:      "too long (41 chars)",
			input:     "zts1qsrxxxxxxxxxxxxxmrhjllzzzzzzzzzzzz",
			expectErr: true,
		},
		{
			name:      "wrong length (26 chars)",
			input:     "zts1qsrxxxxxxxxxxxxmrhj",
			expectErr: true,
		},
		{
			name:      "empty string",
			input:     "",
			expectErr: true,
		},
		{
			name:      "invalid token name",
			input:     "INVALID",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTokenStandard(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestValidateAddress tests the ValidateAddress function
func TestValidateAddress(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "valid address",
			input:     "z1qzal6c5s9rjnnxd2z672tx3apscy5s5qqhslq5",
			expectErr: false,
		},
		{
			name:      "valid address with spaces",
			input:     " z1qzal6c5s9rjnnxd2z672tx3apscy5s5qqhslq5 ",
			expectErr: false,
		},
		{
			name:      "another valid address",
			input:     "z1qqjnwjjpnue8xmmpanz6csze6tcmtzzdtfsww7",
			expectErr: false,
		},
		{
			name:      "empty string",
			input:     "",
			expectErr: true,
		},
		{
			name:      "wrong prefix",
			input:     "x1qzal6c5s9rjnnxd2z672tx3apscy5s5qqhslq5",
			expectErr: true,
		},
		{
			name:      "too short",
			input:     "z1qzal6c5s9rjnnxd2z672tx3apscy5s5qqhslq",
			expectErr: true,
		},
		{
			name:      "too long",
			input:     "z1qzal6c5s9rjnnxd2z672tx3apscy5s5qqhslq5z",
			expectErr: true,
		},
		{
			name:      "invalid format",
			input:     "invalid-address",
			expectErr: true,
		},
		{
			name:      "uppercase (invalid)",
			input:     "Z1QZAL6C5S9RJNNXD2Z672TX3APSCY5S5QQHSLQ5",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAddress(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestSetVerbose tests the verbose mode setter and getter
func TestSetVerbose(t *testing.T) {
	// Save original state
	originalVerbose := GetVerbose()
	defer SetVerbose(originalVerbose)

	// Test setting to true
	SetVerbose(true)
	assert.True(t, GetVerbose())

	// Test setting to false
	SetVerbose(false)
	assert.False(t, GetVerbose())
}

// TestFormatZNN tests ZNN formatting (mainly checks it doesn't panic)
func TestFormatZNN(t *testing.T) {
	amount := big.NewInt(100000000)
	result := FormatZNN(amount)
	assert.Contains(t, result, "1.00000000")
	assert.Contains(t, result, "ZNN")
}

// TestFormatQSR tests QSR formatting (mainly checks it doesn't panic)
func TestFormatQSR(t *testing.T) {
	amount := big.NewInt(100000000)
	result := FormatQSR(amount)
	assert.Contains(t, result, "1.00000000")
	assert.Contains(t, result, "QSR")
}

// TestFormatToken tests token formatting with different symbols
func TestFormatToken(t *testing.T) {
	tests := []struct {
		name     string
		amount   *big.Int
		decimals int
		symbol   string
	}{
		{
			name:     "ZNN token",
			amount:   big.NewInt(100000000),
			decimals: 8,
			symbol:   "ZNN",
		},
		{
			name:     "QSR token",
			amount:   big.NewInt(100000000),
			decimals: 8,
			symbol:   "QSR",
		},
		{
			name:     "custom token",
			amount:   big.NewInt(1000000),
			decimals: 6,
			symbol:   "TEST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatToken(tt.amount, tt.decimals, tt.symbol)
			assert.Contains(t, result, tt.symbol)
		})
	}
}

// TestAmountRoundTrip tests that parsing and formatting are inverse operations
func TestAmountRoundTrip(t *testing.T) {
	tests := []string{
		"1.0",
		"0.5",
		"123.456789",
		"1000000.0",
		"0.00000001",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			// Parse to big.Int
			amount, err := ParseAmount(input, 8)
			require.NoError(t, err)

			// Format back to string
			output := Amount(amount, 8)

			// Parse again
			amount2, err := ParseAmount(output, 8)
			require.NoError(t, err)

			// Should be equal
			assert.Equal(t, amount, amount2)
		})
	}
}

// BenchmarkAmount benchmarks the Amount formatting function
func BenchmarkAmount(b *testing.B) {
	amount := big.NewInt(123456789012345678)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Amount(amount, 8)
	}
}

// BenchmarkParseAmount benchmarks the ParseAmount parsing function
func BenchmarkParseAmount(b *testing.B) {
	input := "1234567890.12345678"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseAmount(input, 8)
	}
}
