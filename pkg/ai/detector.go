package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

type StackInfo struct {
	Language string `json:"language"`
	Version  string `json:"version"`
}

func DetectStack(projectDir string) ([]StackInfo, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")

	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable is not set")
	}

	filesToCheck := []string{"go.mod", "requirements.txt", "Pipfile", "pyproject.toml", "package.json"}

	var sb strings.Builder

	skipDirs := map[string]bool{
		"node_modules": true,
		"vendor":       true,
		".git":         true,
		".venv":        true,
		"__pycache__":  true,
	}

	filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if skipDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		for _, f := range filesToCheck {
			if info.Name() == f {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				fmt.Fprintf(&sb, "--- %s ---\n%s\n", info.Name(), string(content))
				break
			}
		}
		return nil
	})

	if sb.Len() == 0 {
		return nil, fmt.Errorf("no recognizable project files found in %s", projectDir)
	}

	prompt := fmt.Sprintf(
		"Based on these project files, identify all programming languages and their versions.\n"+
			"Reply with ONLY a valid JSON array in this exact format, no extra text: [{\"language\": \"go\", \"version\": \"1.23\"}, {\"language\": \"nodejs\", \"version\": \"20\"}]\n\n%s",
		sb.String(),
	)

	client := anthropic.NewClient()

	resp, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeHaiku4_5,
		MaxTokens: 512,
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
			// Extract the JSON array, ignoring any surrounding markdown or text
			if start := strings.IndexByte(text, '['); start != -1 {
				if end := strings.LastIndexByte(text, ']'); end > start {
					text = text[start : end+1]
				}
			}
			var result []StackInfo
			if err := json.Unmarshal([]byte(text), &result); err != nil {
				return nil, fmt.Errorf("failed to parse Claude response: %w", err)
			}
			return result, nil
		}
	}

	return nil, fmt.Errorf("no text response received from Claude")
}
