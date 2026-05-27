package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
)

type RiskItem struct {
	Product      string `json:"product"`
	Version      string `json:"version"`
	EOL          string `json:"eol"`
	RiskLevel    string `json:"risklevel"`
	DaysUntilEOL int    `json:"daysuntileol"`
}

type UpgradeItem struct {
	Language  string `json:"language"`
	Version   string `json:"version"`
	EOL       string `json:"eol"`
	RiskLevel string `json:"risklevel"`
}

func GenerateRiskNarrative(items []RiskItem) (string, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")

	if apiKey == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY environment variable is not set")
	}

	client := anthropic.NewClient()

	data, err := json.MarshalIndent(items, "", " ")

	if err != nil {
		return "", fmt.Errorf("failed to serialize risk data: %w", err)
	}

	prompt := fmt.Sprintf("You are a security advisor. Here is an EOL risk assessment. Provide a concise 3-5 sentence summary prioritizing the most critical issues: \n%s", string(data))

	resp, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeOpus4_6,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})

	if err != nil {
		return "", fmt.Errorf("failed to call Claude API: %w", err)
	}

	for _, block := range resp.Content {
		switch b := block.AsAny().(type) {
		case anthropic.TextBlock:
			return b.Text, nil
		}
	}

	return "", fmt.Errorf("no text response received from Claude")
}

func SuggestUpgradePath(items []UpgradeItem) (string, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")

	if apiKey == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY environment variable is not set")
	}

	client := anthropic.NewClient()

	data, err := json.MarshalIndent(items, "", " ")
	if err != nil {
		return "", fmt.Errorf("failed to serialize upgrade data: %w", err)
	}

	prompt := fmt.Sprintf("You are a software upgrade advisor. For each component listed below, recommend the specific version to upgrade to and briefly explain why. Be concise and direct:\n%s", string(data))

	resp, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeOpus4_6,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})

	if err != nil {
		return "", fmt.Errorf("failed to call Claude API: %w", err)
	}

	for _, block := range resp.Content {
		switch b := block.AsAny().(type) {
		case anthropic.TextBlock:
			return b.Text, nil
		}
	}

	return "", fmt.Errorf("no text response received from Claude")
}
