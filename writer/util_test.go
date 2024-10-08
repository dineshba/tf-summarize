package writer

import (
	"testing"
)

// Test for hasOutputChanges function
func TestHasOutputChanges(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string][]string
		expected bool
	}{
		{
			name: "No changes",
			input: map[string][]string{
				"change1": {},
				"change2": {},
			},
			expected: false,
		},
		{
			name: "Has changes",
			input: map[string][]string{
				"change1": {"added"},
				"change2": {},
			},
			expected: true,
		},
		{
			name: "All empty",
			input: map[string][]string{
				"change1": {},
				"change2": {},
				"change3": {},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasOutputChanges(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test for removeANSI function
func TestRemoveANSI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "String without ANSI escape",
			input:    "This is a normal string",
			expected: "This is a normal string",
		},
		{
			name:     "String with ANSI escape",
			input:    "\x1b[31mThis is red\x1b[0m",
			expected: "This is red",
		},
		{
			name:     "String with multiple ANSI escapes",
			input:    "\x1b[31mThis is red\x1b[0m and \x1b[32mthis is green\x1b[0m",
			expected: "This is red and this is green",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeANSI(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
