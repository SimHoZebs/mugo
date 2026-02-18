package agents

import (
	"log"
	"strings"
)

func StripMarkdownFences(text string, agentName string) string {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "```") {
		log.Printf("Warning: Gemini wrapped response in markdown code blocks for %s", agentName)
		lines := strings.Split(text, "\n")
		if len(lines) > 2 {
			text = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}
	return strings.TrimSpace(text)
}
