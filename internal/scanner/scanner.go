package scanner

import (
	// "fmt"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var languageMap = map[string]string{
	".go": "go",
	".js": "nodejs",
	".py": "python",
	".tf": "terraform",
}

// DetectLanguage scans the directory and identifies the programming language based on file extensions.
func DetectLanguages(projectDir string) ([]string, []string, error) {
	languageCount := make(map[string]int)
	projectDirs := make(map[string]struct{})

	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if it's a regular file (not a directory)
		if !info.IsDir() {
			ext := filepath.Ext(path)
			if lang, ok := languageMap[ext]; ok {
				languageCount[lang]++

				dir := filepath.Dir(path)
				projectDirs[dir] = struct{}{}
			}

			// Additional check for project-defining files
			projectIndicators := []string{"package.json", "requirements.txt", "go.mod", "pom.xml"}
			if contains(projectIndicators, info.Name()) {
				dir := filepath.Dir(path)
				projectDirs[dir] = struct{}{}
			}
		}
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	// Find the most frequent language
	var detectedLanguages []string
	var detectedProjects []string

	for dir := range projectDirs {
		detectedProjects = append(detectedProjects, dir)
	}

	for lang := range languageCount {
		detectedLanguages = append(detectedLanguages, lang)
	}

	if len(detectedLanguages) == 0 {
		return []string{"Unknown"}, detectedProjects, nil
	}

	return detectedLanguages, detectedProjects, nil
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}

func DetectVersion(projectDir string) (string, error) {
	var detectedVersion string
	// Iterate through all files and subdirectories
	err := filepath.Walk(projectDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories, only process files
		if info.IsDir() {
			return nil
		}

		// Define file patterns to check (you can extend this)
		filesToCheck := map[string]*regexp.Regexp{
			// "package.json":   regexp.MustCompile(`"node":\s*"([0-9\.]+)`),
			"Pipfile":        regexp.MustCompile(`python_version\s*=\s*['"]([0-9\.]+)['"]`),
			"pyproject.toml": regexp.MustCompile(`python\s*=\s*['"]([0-9\.]+)['"]`),
			"go.mod":         regexp.MustCompile(`go\s([0-9\.]+)`),
			"setup.py":       regexp.MustCompile(`python_requires\s*=\s*['"]([0-9\.]+)['"]`),
		}

		// Get the file name and decide how to process based on the file type
		for file, regex := range filesToCheck {
			if strings.HasSuffix(filePath, file) {
				f, err := os.Open(filePath)
				if err != nil {
					fmt.Printf("Could not open %s: %v\n", filePath, err)
					return nil
				}
				defer f.Close()

				// // Special handling for package.json (Node.js)
				// if file == "package.json" {
				// 	var jsonContent map[string]interface{}
				// 	if err := json.NewDecoder(f).Decode(&jsonContent); err != nil {
				// 		fmt.Printf("Failed to parse JSON in %s: %v\n", filePath, err)
				// 		return nil
				// 	}

				// 	// Look for the "engines.node" field
				// 	if engines, ok := jsonContent["engines"].(map[string]interface{}); ok {
				// 		if nodeVersion, ok := engines["node"].(string); ok {
				// 			if matches := regex.FindStringSubmatch(nodeVersion); matches != nil {
				// 				detectedVersion = matches[1]
				// 				return nil
				// 			}
				// 		}
				// 	}
				// 	continue
				// }

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

					// Match Go version in go.mod
					if file == "go.mod" {
						if matches := regex.FindStringSubmatch(line); matches != nil {
							detectedVersion = matches[1]
							return nil
						}
					}

					// Match Python version in Python-related files
					if inRelevantSection || file == "setup.py" {
						if matches := regex.FindStringSubmatch(line); matches != nil {
							detectedVersion = matches[1]
							return nil // Return the version requirement found
						}
					}
				}

				// If no match was found in the current file, show a debug message
				fmt.Printf("No version requirement found in %s\n", filePath)
			}
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("error walking the path %v: %v", projectDir, err)
	}

	if detectedVersion == "" {
		return "", fmt.Errorf("version requirement not found in any checked files")
	}

	return detectedVersion, nil
}
