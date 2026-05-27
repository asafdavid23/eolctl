package helm

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type HelmRelease struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Chart      string `json:"chart"`
	AppVersion string `json:"app_version"`
	Status     string `json:"status"`
}

func ListReleases() ([]HelmRelease, error) {
	cmd := exec.Command("helm", "list", "--all-namespaces", "-o", "json")
	output, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf("failed to execute helm command: %w", err)
	}

	var releases []HelmRelease
	err = json.Unmarshal(output, &releases)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal helm output: %w", err)
	}

	return releases, nil
}
