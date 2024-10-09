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
	"strings"
)

type Release struct {
	Cycle             string `json:"cycle"`
	ReleaseDate       string `json:"releaseDate"`
	Latest            string `json:"latest"`
	LatestReleaseDate string `json:"latestReleaseDate"`
	LTS               bool   `json:"lts"`
}

var languageDetailesFile = map[string]string{
	"go.mod":           "go",
	"package.json":     "node.js",
	"requirements.txt": "python",
	"java":             ".java",
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

func GetProduct(product string, version string, outputFolder string, minVersion string, maxVersion string) []byte {

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

// Helper function to check if a cycle is within a given range
func IsWithinRange(cycle, minVersion, maxVersion string) bool {
	return strings.Compare(cycle, minVersion) >= 0 && strings.Compare(cycle, maxVersion) <= 0
}

func FilterVersions(outputData []byte, minVersion, maxVersion string) ([]byte, error) {

	var releases []Release
	var filteredReleases []map[string]string

	if err := json.Unmarshal([]byte(outputData), &releases); err != nil {
		log.Fatal(err)
	}

	for _, release := range releases {

		if IsWithinRange(release.Cycle, minVersion, maxVersion) {

			releaseMap := map[string]string{
				"cycle":             release.Cycle,
				"releaseDate":       release.ReleaseDate,
				"latest":            release.Latest,
				"latestReleaseDate": release.LatestReleaseDate,
			}

			filteredReleases = append(filteredReleases, releaseMap)
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

	_, err = fmt.Fprint(file, string(body))

	if err != nil {
		log.Fatalf("Error writing to the file: %v", err)
	}

	log.Printf("Content written to file successfully, output file located in: %s", filePath)
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

	log.Printf("Identifed Product is: %s", product)
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

	log.Printf("Identified product versoin is: %s", version)

	return version
}

func CompareTwoVersions(version1, version2 []byte) ([]byte, []byte) {

	// if len(version1) != len(version2) {
	// 	return
	// }

	for i := range version1 {
		if version1[i] != version2[i] {
			return version1, version2
		}
	}

	return nil, nil
}
