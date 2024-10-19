package scanner

import (
	// "fmt"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
			if filepath.Base(path) == "package.json" || filepath.Base(path) == "go.mod" {
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
		return "", fmt.Errorf("package.json file not found")
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
func DetectVersionFromRequirementsTxt(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "python==") {
			version := strings.Split(line, "==")[1] // Extract the version after "python=="
			return version, nil
		}
	}
	return "Unknown", nil
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
