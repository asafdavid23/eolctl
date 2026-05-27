package artifacthub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 15 * time.Second}

type Package struct {
	Name       string     `json:"name"`
	Version    string     `json:"version"`
	AppVersion string     `json:"app_version"`
	Deprecated bool       `json:"deprecated"`
	Repository Repository `json:"repository"`
}

type Repository struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type searchResponse struct {
	Packages []Package `json:"packages"`
}

// SearchPackage queries ArtifactHub for the most relevant stable Helm package matching name.
// It fetches up to 5 candidates and returns the first whose version is not a pre-release.
func SearchPackage(name string) (*Package, error) {
	endpoint := fmt.Sprintf(
		"https://artifacthub.io/api/v1/packages/search?kind=0&ts_query_web=%s&limit=5",
		url.QueryEscape(name),
	)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query ArtifactHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ArtifactHub returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result searchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Packages) == 0 {
		return nil, fmt.Errorf("no packages found for %q", name)
	}

	// Prefer the first stable (non-pre-release) result
	for i := range result.Packages {
		pkg := &result.Packages[i]
		if !isPreRelease(pkg.AppVersion) && !isPreRelease(pkg.Version) {
			return pkg, nil
		}
	}

	// All results are pre-releases — fall back to the top result
	return &result.Packages[0], nil
}

// isPreRelease returns true if the version string looks like a pre-release
// (contains -rc, -alpha, -beta, -dev, -pre, -snapshot, or a git hash suffix).
func isPreRelease(version string) bool {
	v := strings.ToLower(strings.TrimPrefix(version, "v"))
	for _, marker := range []string{"-rc", "-alpha", "-beta", "-dev", "-pre", "-snapshot", ".g"} {
		if strings.Contains(v, marker) {
			return true
		}
	}
	return false
}

// RiskFromStaleness derives a risk level by comparing the installed app version
// against the latest version reported by ArtifactHub.
// Returns a risk label (CRITICAL/HIGH/MEDIUM/LOW/UNKNOWN) and a short EOL-column label.
func RiskFromStaleness(installed, latest string, deprecated bool) (riskLevel, eolLabel string) {
	if deprecated {
		return "CRITICAL", "deprecated"
	}

	installed = strings.TrimPrefix(strings.TrimSpace(installed), "v")
	latest = strings.TrimPrefix(strings.TrimSpace(latest), "v")

	if installed == "" || latest == "" {
		return "UNKNOWN", "unknown"
	}

	if installed == latest {
		return "LOW", fmt.Sprintf("latest: %s", latest)
	}

	installedMajor := majorOf(installed)
	latestMajor := majorOf(latest)
	majorDiff := latestMajor - installedMajor

	if majorDiff >= 2 {
		return "CRITICAL", fmt.Sprintf("latest: %s", latest)
	}
	if majorDiff == 1 {
		return "HIGH", fmt.Sprintf("latest: %s", latest)
	}

	// Same major — check minor
	installedMinor := minorOf(installed)
	latestMinor := minorOf(latest)
	minorDiff := latestMinor - installedMinor

	if minorDiff >= 3 {
		return "HIGH", fmt.Sprintf("latest: %s", latest)
	}
	if minorDiff >= 1 {
		return "MEDIUM", fmt.Sprintf("latest: %s", latest)
	}

	return "LOW", fmt.Sprintf("latest: %s", latest)
}

func majorOf(version string) int {
	parts := strings.SplitN(version, ".", 3)
	n, _ := strconv.Atoi(parts[0])
	return n
}

func minorOf(version string) int {
	parts := strings.SplitN(version, ".", 3)
	if len(parts) < 2 {
		return 0
	}
	n, _ := strconv.Atoi(parts[1])
	return n
}
