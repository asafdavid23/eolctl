package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 15 * time.Second}

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

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func GetProduct(product string, version string) ([]byte, error) {
	endpoint := fmt.Sprintf("https://endoflife.date/api/%s.json", url.PathEscape(product))

	if version != "" {
		endpoint = fmt.Sprintf("https://endoflife.date/api/%s/%s.json", url.PathEscape(product), url.PathEscape(version))
	}

	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Add("Accept", "application/json")

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from the API: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	return body, err
}

func compareVersions(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}

	for i := 0; i < maxLen; i++ {
		var aNum, bNum int
		if i < len(aParts) {
			aNum, _ = strconv.Atoi(aParts[i])
		}
		if i < len(bParts) {
			bNum, _ = strconv.Atoi(bParts[i])
		}
		if aNum != bNum {
			if aNum > bNum {
				return 1
			}
			return -1
		}
	}
	return 0
}

func IsWithinRange(cycle, minVersion, maxVersion string) bool {
	return compareVersions(cycle, minVersion) >= 0 && compareVersions(cycle, maxVersion) <= 0
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
		return nil, fmt.Errorf("failed to marshal filtered releases: %w", err)
	}

	return filteredReleasesJSON, nil
}

func GetStringValue(value interface{}) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case bool:
		if v {
			return "true"
		}
		return "false"
	}
	return ""
}

func ExportToFile(outputData []byte, outputFolder string) error {

	err := os.MkdirAll(outputFolder, os.ModePerm)

	if err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}

	filePath := filepath.Join(outputFolder, "output.json")

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
		return false, fmt.Sprintf("Product %s version %s is not EOL", product, version), nil
	}

	// Check if the current date is after the EOL date
	currentDate := time.Now()

	if currentDate.After(eolDate) {
		return true, fmt.Sprintf("Product %s version %s is EOL", product, version), nil
	}

	return false, fmt.Sprintf("Product %s version %s is not EOL", product, version), nil
}
