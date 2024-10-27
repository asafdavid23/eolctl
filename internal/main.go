package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io"
	"net/http"
	"os"
	"strings"
)

func GetAvailableProducts(output string) ([]byte, error) {
	url := "https://endoflife.date/api/all.json" // Update this with the correct URL

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var products []interface{}
	if err := json.Unmarshal(body, &products); err != nil {
		return nil, fmt.Errorf("faild to parse JSON response: %w", err)
	}

	if output == "table" {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Product"})

		for _, product := range products {
			if str, ok := product.(string); ok {
				table.Append([]string{str})
			}
		}

		table.Render()
	} else if output == "json" {
		fmt.Print(string(body))
	}

	return body, nil
}

func GetProduct(product string, version string, output string) ([]byte, error) {

	url := fmt.Sprintf("https://endoflife.date/api/%s.json", product)

	if version != "" {
		url = fmt.Sprintf("https://endoflife.date/api/%s/%s.json", product, version)
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from the API: %w", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("faild to parse JSON response: %w", err)
	}

	if output == "table" {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Latest", "LatestReleaseDate", "ReleaseDate", "LTS", "EOL", "SUPPORT"})

		row := []string{
			result["latest"],
		}
	}

	return nil, err
}

// Helper function to check if a cycle is within a given range
func IsWithinRange(cycle, minVersion, maxVersion string) bool {
	return strings.Compare(cycle, minVersion) >= 0 && strings.Compare(cycle, maxVersion) <= 0
}

func FilterVersions(outputData []byte, minVersion, maxVersion string) ([]byte, error) {

	var result []map[string]interface{}
	var filteredReleases []map[string]string

	if err := json.Unmarshal(outputData, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %w", err)
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

func ExportToFile(outputData []byte, outputFolder string) error {

	err := os.MkdirAll(outputFolder, os.ModePerm)

	if err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}

	filePath := outputFolder + "/output.json"

	body := outputData

	file, err := os.Create(filePath)

	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer file.Close()

	_, err = fmt.Fprint(file, string(body))

	if err != nil {
		return fmt.Errorf("error writing to the file: %w", err)
	}

	fmt.Printf("Content written to file successfully, output file located in: %s", filePath)
	return nil
}

func ConvertOutput(outputData []byte, outputType string) error {
	// Define a variable to hold the unmarshalled data
	var result interface{}

	// Unmarshal the JSON into the map
	if err := json.Unmarshal(outputData, &result); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}

	// Switch case to handle different output types
	switch outputType {
	case "table":
		PrintTable(result)
	default:
		return fmt.Errorf("invalid output type: %s", outputType)
	}

	return nil
}

func PrintTable(data interface{}) {
	table := tablewriter.NewWriter(os.Stdout)

	switch v := data.(type) {
	case []interface{}:
		if _, ok := v[0].(string); ok {
			// Handle []string case
			table.SetHeader([]string{"Product"})
			for _, item := range v {
				if str, ok := item.(string); ok {
					table.Append([]string{str})
				}
			}
		} else {
			table.SetHeader([]string{"Cycle", "Latest", "LatestReleaseDate", "ReleaseDate", "LTS", "EOL", "SUPPORT"})
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
		}
	case map[string]interface{}:
		table.SetHeader([]string{"Latest", "LatestReleaseDate", "ReleaseDate", "LTS", "EOL", "SUPPORT"})

		row := []string{
			getStringValue(v["latest"]),
			getStringValue(v["latestReleaseDate"]),
			getStringValue(v["releaseDate"]),
			getStringValue(v["lts"]),
			getStringValue(v["eol"]),
			getStringValue(v["support"]),
		}
		table.Append(row)
	}

	table.Render()
}
