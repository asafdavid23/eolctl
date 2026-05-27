package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/asafdavid23/eolctl/pkg/helm"
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

// HelmStackInfo maps a Helm release to an endoflife.date product slug and version.
type HelmStackInfo struct {
	ReleaseName string `json:"release_name"`
	Language    string `json:"language"`
	Version     string `json:"version"`
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

	prompt := fmt.Sprintf("You are a software upgrade advisor. For each component listed below, recommend the specific version to upgrade to and briefly explain why. Use plain text only — no markdown, no bullet points, no tables, no bold or italic syntax:\n%s", string(data))

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

func DetectHelmEOL(releases []helm.HelmRelease) ([]HelmStackInfo, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable is not set")
	}

	data, err := json.MarshalIndent(releases, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to serialize helm releases: %w", err)
	}

	prompt := fmt.Sprintf(
		"You are given a list of Helm chart releases from a Kubernetes cluster.\n"+
			"For EVERY release, map the chart to its corresponding product slug on endoflife.date and extract the major.minor version from app_version.\n"+
			"Return ONLY a valid JSON array in this exact format, no extra text:\n"+
			"[{\"release_name\": \"my-nginx\", \"language\": \"nginx\", \"version\": \"1.25\"}]\n"+
			"Rules:\n"+
			"- Include ALL releases — do not omit any.\n"+
			"- Use known endoflife.date slugs when possible (e.g. nginx, kubernetes, redis, postgresql, cert-manager, prometheus, grafana).\n"+
			"- If unsure, use the chart name (without the version suffix) as the slug.\n"+
			"- Extract major.minor from app_version (e.g. \"1.25.3\" → \"1.25\").\n\n%s",
		string(data),
	)

	client := anthropic.NewClient()

	resp, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeHaiku4_5,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call Claude API: %w", err)
	}

	for _, block := range resp.Content {
		if b, ok := block.AsAny().(anthropic.TextBlock); ok {
			text := strings.TrimSpace(b.Text)
			if start := strings.IndexByte(text, '['); start != -1 {
				if end := strings.LastIndexByte(text, ']'); end > start {
					text = text[start : end+1]
				}
			}
			var result []HelmStackInfo
			if err := json.Unmarshal([]byte(text), &result); err != nil {
				return nil, fmt.Errorf("failed to parse Claude response: %w", err)
			}
			return result, nil
		}
	}

	return nil, fmt.Errorf("no text response received from Claude")
}
