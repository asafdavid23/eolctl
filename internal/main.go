package helpers

import (
	"encoding/json"
	"fmt"
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

	return body, nil
}

func GetProduct(product string, version string) ([]byte, error) {
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
		return nil, fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}

	for _, release := range result {
		// Check if Cycle is present and within range
		if cycleVal, ok := release["cycle"]; ok && cycleVal != nil {
			if cycle, ok := cycleVal.(string); ok && IsWithinRange(cycle, minVersion, maxVersion) {
				// Create a new map for filtered release data
				releaseMap := map[string]string{
					"cycle":             cycle,
					"releaseDate":       GetStringValue(release["releaseDate"]),
					"latest":            GetStringValue(release["latest"]),
					"latestReleaseDate": GetStringValue(release["latestReleaseDate"]),
					"lts":               GetStringValue(release["lts"]),
					"eol":               GetStringValue(release["eol"]),
					"support":           GetStringValue(release["support"]),
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

func GetStringValue(value interface{}) string {
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

	filePath := outputFolder + "output.json"

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

// func PrintOutput(outputData []json.RawMessage, output string) error {
// 	if output == "table" {
// 		var data []map[string]interface{}

// 		for _, rawMessage := range outputData {
// 			var item map[string]interface{}
// 			if err := json.Unmarshal(rawMessage, &item); err != nil {
// 				return fmt.Errorf("failed to unmarshal JSON raw message: %w", err)
// 			}
// 			data = append(data, item)
// 		}

// 		headers := []string{}

// 		if len(data) > 0 {
// 			for key := range data[0] {
// 				headers = append(headers, key)
// 			}
// 		}

// 		table := tablewriter.NewWriter(os.Stdout)
// 		table.SetHeader(headers)

// 		for _, row := range data {
// 			rowData := []string{}

// 			for _, header := range headers {
// 				if value, exists := row[header]; exists {
// 					rowData = append(rowData, fmt.Sprintf("%v", value))
// 				} else {
// 					rowData = append(rowData, "")
// 				}
// 			}
// 			table.Append(rowData)
// 		}

// 		table.Render()

// 	} else if output == "json" {
// 		productsJSON, err := json.Marshal(outputData)

// 		if err != nil {
// 			return fmt.Errorf("failed to marshal JSON data: %w", err)
// 		}

// 		fmt.Print(string(productsJSON))
// 	} else {
// 		return fmt.Errorf("output type is not valid")
// 	}

// 	return nil
// }
