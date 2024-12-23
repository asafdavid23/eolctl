package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type EOL struct {
	Value interface{} `json:"eol"`
}

type ApiResponse struct {
	Cycle             string `json:"cycle"`
	ReleaseDate       string `json:"releaseDate"`
	EOL               EOL    `json:"eol"`
	Latest            string `json:"latest"`
	LatestReleaseDate string `json:"latestReleaseDate"`
	LTS               bool   `json:"lts"`
}

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

func (e *EOL) UnmarshalJSON(data []byte) error {
	var boolVal bool
	if err := json.Unmarshal(data, &boolVal); err == nil {
		e.Value = boolVal
		return nil
	}

	var stringVal string
	if err := json.Unmarshal(data, &stringVal); err == nil {
		e.Value = stringVal
		return nil
	}

	return fmt.Errorf("invalid EOL value: %s", string(data))
}

func CheckProductEOL(product string, version string) (bool, string, error) {
	var response ApiResponse

	productData, err := GetProduct(product, version)

	if err != nil {
		return false, "", fmt.Errorf("failed to fetch product data: %w", err)
	}

	// parse the product data into the response struct
	if err := json.Unmarshal(productData, &response); err != nil {
		return false, "", fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}

	// parse the date
	var eolDate time.Time
	switch value := response.EOL.Value.(type) {
	case string:
		var err error
		eolDate, err = time.Parse("2006-01-02", value)

		if err != nil {
			return false, "", fmt.Errorf("failed to parse EOL date: %w", err)
		}
	case bool:
		if value {
			return true, fmt.Sprintf("Product %s version %s is EOL", product, version), nil
		}
	}

	if err != nil {
		return false, "", fmt.Errorf("failed to parse EOL date: %w", err)
	}

	// Check if the current date is after the EOL date
	currentDate := time.Now()

	if currentDate.After(eolDate) {
		return true, fmt.Sprintf("Product %s version %s is EOL", product, version), nil
	}

	return false, fmt.Sprintf("Product %s version %s is not EOL", product, version), nil
}
