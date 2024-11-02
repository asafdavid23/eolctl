package scanner

import (
	// "fmt"

	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var languageExtensions = map[string]string{
	".go":   "Go",
	".py":   "Python",
	".js":   "JavaScript",
	".ts":   "TypeScript",
	".java": "Java",
}

func DetectLanguage(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))

	if lang, exists := languageExtensions[ext]; exists {
		return lang
	}

	return ""
}

func IdentifyLanguages(dir string, recurse bool) (map[string]string, error) {
	projectLanguages := make(map[string]string)

	// Helper function to identify language based on file extensions.
	findLanguage := func(path string) string {
		ext := filepath.Ext(path)
		if lang, exists := languageExtensions[ext]; exists {
			return lang
		}
		return ""
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If it's a directory and we are not recursing, skip further scanning
		if info.IsDir() && path != dir && !recurse {
			return filepath.SkipDir
		}

		// Check if file extension matches a known language
		if !info.IsDir() {
			lang := findLanguage(path)
			if lang != "" {
				projectDir := filepath.Dir(path)
				if _, exists := projectLanguages[projectDir]; !exists {
					projectLanguages[projectDir] = lang
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return projectLanguages, nil
}

// IdentifyPythonVersion identifies the Python version in a Python project.
func IdentifyPythonVersion(projectPath string) (string, error) {
	// Check for `pyproject.toml` or `requirements.txt`
	pyprojectPath := projectPath + "/pyproject.toml"
	pipfilePath := projectPath + "/Pipfile"
	setupPath := projectPath + "/setup.py"

	// Check `pyproject.toml` first
	if _, err := os.Stat(pyprojectPath); err == nil {
		content, err := os.ReadFile(pyprojectPath)
		if err != nil {
			return "", err
		}
		re := regexp.MustCompile(`(?i)requires-python\s*=\s*['"][><=~]*\s*([\d\.]+)['"]`)
		match := re.FindStringSubmatch(string(content))
		if len(match) > 1 {
			return match[1], nil
		}
	}

	// Fallback to `Pipfile`
	if _, err := os.Stat(pipfilePath); err == nil {
		content, err := os.ReadFile(pipfilePath)
		if err != nil {
			return "", err
		}
		re := regexp.MustCompile(`(?i)python_version\s*=\s*['"][><=~]*\s*([\d\.]+)['"]`)
		match := re.FindStringSubmatch(string(content))
		if len(match) > 1 {
			return match[1], nil
		}
	}

	// Fallback to `setup.py`

	if _, err := os.Stat(setupPath); err == nil {
		content, err := os.ReadFile(setupPath)

		if err != nil {
			return "", fmt.Errorf("cant read file %v", err)
		}

		re := regexp.MustCompile(`(?i)python_requires\s*=\s*['"][><=~]*\s*([\d\.]+)['"]`)
		match := re.FindStringSubmatch(string(content))

		if len(match) > 1 {
			return match[1], nil
		}
	}

	return "", fmt.Errorf("python version not found")
}

// IdentifyGoVersion identifies the Go version in a Go project.
func IdentifyGoVersion(projectPath string) (string, error) {
	goModPath := projectPath + "/go.mod"
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`go\s+(\d+\.\d+)`)
	match := re.FindStringSubmatch(string(content))
	if len(match) > 1 {
		return match[1], nil
	}

	return "", fmt.Errorf("go version not found")
}

// IdentifyNodeVersion identifies the Node.js version in a Node.js project.
func IdentifyNodeVersion(projectPath string) (string, error) {
	packageJsonPath := projectPath + "/package.json"
	content, err := os.ReadFile(packageJsonPath)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return "", err
	}

	if engines, ok := data["engines"].(map[string]interface{}); ok {
		if nodeVersion, ok := engines["node"].(string); ok {
			return nodeVersion, nil
		}
	}

	return "", fmt.Errorf("Node version not found")
}
