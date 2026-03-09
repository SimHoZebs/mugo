package agents

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripMarkdownFences(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "removes json code block",
			input:    "```json\n{\"name\":\"test\"}\n```",
			expected: "{\"name\":\"test\"}",
		},
		{
			name:     "removes plain code block",
			input:    "```\n{\"name\":\"test\"}\n```",
			expected: "{\"name\":\"test\"}",
		},
		{
			name:     "no code blocks",
			input:    "{\"name\":\"test\"}",
			expected: "{\"name\":\"test\"}",
		},
		{
			name:     "trims whitespace",
			input:    "  \n{\"name\":\"test\"}\n  ",
			expected: "{\"name\":\"test\"}",
		},
		{
			name:     "code block with whitespace",
			input:    "  ```json\n{\"name\":\"test\"}\n```  ",
			expected: "{\"name\":\"test\"}",
		},
		{
			name:     "incomplete code block - only opening",
			input:    "```\nonly opening fence",
			expected: "```\nonly opening fence",
		},
		{
			name:     "incomplete code block - only one line",
			input:    "```",
			expected: "```",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "multiline content in code block",
			input:    "```json\n{\"name\":\"test\",\n\"value\":123}\n```",
			expected: "{\"name\":\"test\",\n\"value\":123}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripMarkdownFences(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
