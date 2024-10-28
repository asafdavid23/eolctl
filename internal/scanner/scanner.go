package scanner

import (
	// "fmt"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var languageMap = map[string]string{
	".go": "Go",
	".js": "JavaScript",
	".py": "Python",
}

// For reading from package.json
type PackageJSON struct {
	Engines struct {
		Node string `json:"node"`
	} `json:"engines"`
}

// DetectLanguage scans the directory and identifies the programming language based on file extensions.
func DetectLanguage(projectDir string) (string, error) {
	languageCount := make(map[string]int)

	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if it's a regular file (not a directory)
		if !info.IsDir() {
			ext := filepath.Ext(path)
			if lang, ok := languageMap[ext]; ok {
				languageCount[lang]++
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	// Find the most frequent language
	var detectedLang string
	var maxCount int
	for lang, count := range languageCount {
		if count > maxCount {
			maxCount = count
			detectedLang = lang
		}
	}

	if detectedLang == "" {
		return "Unknown", nil
	}
	return detectedLang, nil
}

func DetectPackgesFile(projectDir string) (string, error) {
	var packagesFile string

	err := filepath.WalkDir(projectDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			if filepath.Base(path) == "package.json" || filepath.Base(path) == "go.mod" || filepath.Base(path) == "setup.py" || filepath.Base(path) == "pyproject.toml" || filepath.Base(path) == "Pipfile" {
				packagesFile = path
				return filepath.SkipDir // Stop once we've found the first package.json or go.mod
			}

		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("error walking the directory: %v", err)
	}

	if packagesFile == "" {
		return "", fmt.Errorf("package file not found")
	}

	return packagesFile, nil
}

// DetectVersionFromPackageJSON reads the version from package.json
func DetectVersionFromPackageJSON(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)
	var pkg PackageJSON
	err = json.Unmarshal(byteValue, &pkg)
	if err != nil {
		return "", err
	}

	if pkg.Engines.Node != "" {
		return pkg.Engines.Node, nil
	}
	return "Unknown", nil
}

// DetectVersionFromRequirementsTxt reads versions from requirements.txt
func DetectPythonVersion(path string) (string, error) {
	// Define the files and regex patterns for each file
	filesToCheck := map[string]*regexp.Regexp{
		"setup.py":       regexp.MustCompile(`(?i)python_requires\s*=\s*['"][><=~]*\s*([\d\.]+)['"]`),
		"pyproject.toml": regexp.MustCompile(`(?i)requires-python\s*=\s*['"][><=~]*\s*([\d\.]+)['"]`),
		"Pipfile":        regexp.MustCompile(`(?i)python_version\s*=\s*['"][><=~]*\s*([\d\.]+)['"]`),
	}

	for file, regex := range filesToCheck {
		f, err := os.Open(path + "/" + file)
		if err != nil {
			fmt.Printf("Could not open %s: %v\n", file, err) // Debug: show file open errors
			continue                                         // skip if file does not exist
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		inRelevantSection := false // Tracks relevant sections for Pipfile and pyproject.toml

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			// Handle specific section requirements for Pipfile and pyproject.toml
			if file == "Pipfile" {
				if strings.HasPrefix(line, "[requires]") {
					inRelevantSection = true
					continue
				} else if strings.HasPrefix(line, "[") { // End of [requires] section in Pipfile
					inRelevantSection = false
				}
			} else if file == "pyproject.toml" {
				// Look for relevant sections in pyproject.toml: [project] or [tool.poetry]
				if strings.HasPrefix(line, "[project]") || strings.HasPrefix(line, "[tool.poetry]") {
					inRelevantSection = true
					continue
				} else if strings.HasPrefix(line, "[") { // End of relevant sections in pyproject.toml
					inRelevantSection = false
				}
			}

			// Only search for version if in the relevant section (for Pipfile and pyproject.toml)
			if inRelevantSection || file == "setup.py" {
				if matches := regex.FindStringSubmatch(line); matches != nil {
					return matches[1], nil // Return the version requirement found
				}
			}
		}

		// If no match was found in the current file, show a debug message
		fmt.Printf("No Python version requirement found in %s\n", file)
	}

	return "", fmt.Errorf("Python version requirement not found in any checked files")
}

// DetectVersionFromGoMod reads the Go version from go.mod
func DetectVersionFromGoMod(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "go ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "go ")), nil
		}
	}
	return "Unknown", nil
}
