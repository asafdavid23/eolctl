package helpers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var logger *log.Logger

var languageDetailesFile = map[string]string{
	"go.mod":           "go",
	"package.json":     "node.js",
	"requirements.txt": "python",
	"java":             ".java",
}

func GetAvailableProducts() ([]byte, error) {
	url := "https://endoflife.date/api/all.json"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		logger.Fatal(err)
	}

	req.Header.Add("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		logger.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	// if err != nil {
	// 	logger.Fatal(err)
	// }

	return body, err
}

func GetProduct(product string, version string) ([]byte, error) {
	logger = NewLogger()

	url := fmt.Sprintf("https://endoflife.date/api/%s.json", product)

	if version != "" {
		url = fmt.Sprintf("https://endoflife.date/api/%s/%s.json", product, version)
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		logger.Fatalf("Failed to send request to the API: %v", err)
	}

	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		logger.Fatalf("Failed to fetch data from the API: %v", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	return body, err
}

// Helper function to check if a cycle is within a given range
func IsWithinRange(cycle, minVersion, maxVersion string) bool {
	return strings.Compare(cycle, minVersion) >= 0 && strings.Compare(cycle, maxVersion) <= 0
}

func FilterVersions(outputData []byte, minVersion, maxVersion string) ([]byte, error) {

	var result []map[string]interface{}
	var filteredReleases []map[string]string

	if err := json.Unmarshal(outputData, &result); err != nil {
		logger.Fatal(err)
	}

	for _, release := range result {
		// Check if Cycle is present and within range
		if cycleVal, ok := release["cycle"]; ok && cycleVal != nil {
			if cycle, ok := cycleVal.(string); ok && IsWithinRange(cycle, minVersion, maxVersion) {
				// Create a new map for filtered release data
				releaseMap := map[string]string{
					"cycle":             cycle,
					"releaseDate":       getStringValue(release["releaseDate"]),
					"latest":            getStringValue(release["latest"]),
					"latestReleaseDate": getStringValue(release["latestReleaseDate"]),
					"lts":               getStringValue(release["lts"]),
					"eol":               getStringValue(release["eol"]),
					"support":           getStringValue(release["support"]),
				}

				filteredReleases = append(filteredReleases, releaseMap)
			}
		}
	}

	filteredReleasesJSON, err := json.Marshal(filteredReleases)

	if err != nil {
		logger.Fatal(err)
	}

	return filteredReleasesJSON, nil
}

func getStringValue(value interface{}) string {
	if value == nil {
		return "" // Return an empty string if the value is nil
	}
	if str, ok := value.(string); ok {
		return str // Return the string value
	}
	return "" // Return an empty string if type assertion fails
}

func ExportToFile(outputData []byte, outputFolder string) {

	err := os.MkdirAll(outputFolder, os.ModePerm)

	if err != nil {
		logger.Fatalf("Failed to create a folder: %v", err)
	}

	filePath := outputFolder + "/output.json"

	body := outputData

	file, err := os.Create(filePath)

	if err != nil {
		logger.Fatalf("Failed to create file: %v", err)
	}

	defer file.Close()

	_, err = fmt.Fprint(file, string(body))

	if err != nil {
		logger.Fatalf("Error writing to the file: %v", err)
	}

	logger.Printf("Content written to file successfully, output file located in: %s", filePath)
}

func IdentifyProduct(project string) (string, string) {
	var product string
	var productFile string

	filepath.Walk(project, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			logger.Fatalf("Failed to access files or directories in project path: %v", walkErr)
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

	logger.Printf("Identifed Product is: %s", product)
	return product, productFile
}

func IdentifyProductVersion(product string, project string, productFile string) string {
	var version string

	s := []string{project, "/", productFile}
	path := strings.Join(s, "")

	file, err := os.Open(path)

	if err != nil {
		logger.Fatal(err)
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

	logger.Printf("Identified product versoin is: %s", version)

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

func ConvertOutput(outputData []byte, outputType string) error {
	// Define a variable to hold the unmarshalled data
	var result interface{}

	// Unmarshal the JSON into the map
	if err := json.Unmarshal(outputData, &result); err != nil {
		logger.Fatal(err)
	}

	// Switch case to handle different output types
	switch outputType {
	case "table":
		PrintTable(result)
	case "yaml":
		PrintYaml(result)
	default:
		return fmt.Errorf("invalid output type: %s", outputType)
	}

	return nil
}

func PrintTable(data interface{}) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Cycle", "Latest", "LatestReleaseDate", "ReleaseDate", "LTS", "EOL", "SUPPORT"})

	switch v := data.(type) {
	case []interface{}:
		for _, item := range v {
			if release, ok := item.(map[string]interface{}); ok {
				row := []string{
					getStringValue(release["cycle"]),
					getStringValue(release["latest"]),
					getStringValue(release["latestReleaseDate"]),
					getStringValue(release["releaseDate"]),
					getStringValue(release["lts"]),
					getStringValue(release["eol"]),
					getStringValue(release["support"]),
				}
				table.Append(row)
			}
		}
	case map[string]interface{}:
		row := []string{
			getStringValue(v["releaseDate"]),
			getStringValue(v["latest"]),
			getStringValue(v["latestReleaseDate"]),
			getStringValue(v["lts"]),
			getStringValue(v["eol"]),
			getStringValue(v["support"]),
		}
		table.Append(row)
	}

	table.Render()
}

func PrintYaml(data interface{}) {
	var yamlData []byte
	var err error

	switch v := data.(type) {
	case []interface{}:
		yamlData, err = yaml.Marshal(v)
	case map[string]interface{}:
		yamlData, err = yaml.Marshal(v)
	default:
		logger.Fatalf("unsupported type for yaml marshaling: %T", data)
	}

	if err != nil {
		logger.Fatal(err)
	}

	fmt.Print(string(yamlData))
}
