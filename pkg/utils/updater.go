package utils

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Quote struct {
	Q string `json:"q"`
	A string `json:"a"`
}

// UpdateREADME generates README.md from template with quote + timestamp
func UpdateREADME(q Quote) error {
	templateContent, err := os.ReadFile("README.template.md")
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	quoteText := fmt.Sprintf("💬 *%s*\n\n— **%s**", q.Q, q.A)
	currentTime := time.Now().UTC().Format("January 2, 2006")

	finalContent := strings.ReplaceAll(string(templateContent), "{{QUOTE}}", quoteText)
	finalContent = strings.ReplaceAll(finalContent, "{{DATE}}", currentTime)

	err = os.WriteFile("README.md", []byte(finalContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	return nil
}
