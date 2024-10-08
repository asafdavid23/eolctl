package helpers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var languageDetailesFile = map[string]string{
	"go.mod":           "go",
	"package.json":     "node.js",
	"requirements.txt": "python",
	"java":             ".java",
}

type Release struct {
	Cycle             string `json:"cycle"`
	ReleaseDate       string `json:"releaseDate"`
	EOL               string `json:"eol"`
	Latest            string `json:"latest"`
	LatestReleaseDate string `json:"latestReleaseDate"`
	LTS               bool   `json:"lts"`
}

func GetAvailableProducts() {
	url := "https://endoflife.date/api/all.json"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}

func GetProduct(product string, version string) []byte {

	url := fmt.Sprintf("https://endoflife.date/api/%s.json", product)

	if version != "" {
		url = fmt.Sprintf("https://endoflife.date/api/%s/%s.json", product, version)
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatalf("Failed to send request to the API: %v", err)
	}

	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalf("Failed to fetch data from the API: %v", err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return body
}

func ParseVersion(version string) float64 {
	v, err := strconv.ParseFloat(version, 64)

	if err != nil {
		log.Fatalf("Failed to parse version %s: %v", version, err)
	}
	return v
}

// Helper function to check if a cycle is within a given range
func IsWithinRange(cycle, min, max string) bool {
	// Convert the cycle versions to integers for comparison
	cycleInt, _ := strconv.ParseFloat(cycle, 64)
	minInt, _ := strconv.ParseFloat(min, 64)
	maxInt, _ := strconv.ParseFloat(max, 64)

	return cycleInt >= minInt && cycleInt <= maxInt
}

func FilterVersions(outputData []byte, minVersion, maxVersion string) ([]byte, error) {

	var releases []Release
	var filteredReleases []Release

	if err := json.Unmarshal([]byte(outputData), &releases); err != nil {
		log.Fatal(err)
	}

	for _, release := range releases {
		if IsWithinRange(release.Cycle, minVersion, maxVersion) {
			filteredReleases = append(filteredReleases, release)
		}
	}

	filteredReleasesJSON, err := json.Marshal(filteredReleases)

	if err != nil {
		log.Fatal(err)
	}

	return filteredReleasesJSON, nil
}

func ExportToFile(outputData []byte, outputFolder string) {

	err := os.MkdirAll(outputFolder, os.ModePerm)

	if err != nil {
		log.Fatalf("Failed to create a folder: %v", err)
	}

	filePath := outputFolder + "/output.json"

	body := outputData

	file, err := os.Create(filePath)

	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}

	defer file.Close()

	_, err = fmt.Fprintln(file, string(body))

	if err != nil {
		log.Fatalf("Error writing to the file: %v", err)
	}

	fmt.Println("Content written to file successfully.")
}

func IdentifyProduct(project string) (string, string) {
	var product string
	var productFile string

	filepath.Walk(project, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			log.Fatalf("Failed to access files or directories in project path: %v", walkErr)
		}

		if !info.IsDir() {
			for key := range languageDetailesFile {
				if key == info.Name() {
					product = languageDetailesFile[info.Name()]
					productFile = info.Name()
				} else {
					continue
				}
			}
		}
		return walkErr
	})

	return product, productFile
}

func IdentifyProductVersion(product string, project string, productFile string) string {
	var version string

	s := []string{project, "/", productFile}
	path := strings.Join(s, "")

	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "go") {
			version = strings.TrimSpace(strings.TrimPrefix(line, "go"))
			parts := strings.Split(version, ".")

			if len(parts) >= 2 {
				version = strings.Join(parts[:len(parts)-1], ".")
				return version
			}

		} else if strings.HasPrefix(line, "python==") {
			version = strings.TrimSpace(strings.TrimPrefix(line, "python=="))
			parts := strings.Split(version, ".")

			if len(parts) > 2 {
				version = strings.Join(parts[:len(parts)-1], ".")
				return version
			}

		} else if strings.HasPrefix(line, "<java.version>") && strings.HasSuffix(line, "<java.version>") {
			re := regexp.MustCompile(`<java\.version>(.*?)<\/java\.version>`)
			matches := re.FindStringSubmatch(line)

			if len(matches) > 1 {
				version := matches[1]
				return version
			}

		}
	}

	return version
}
