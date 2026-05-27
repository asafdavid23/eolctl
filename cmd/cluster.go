package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/asafdavid23/eolctl/internal/logging"
	ai "github.com/asafdavid23/eolctl/pkg/ai"
	"github.com/asafdavid23/eolctl/pkg/artifacthub"
	"github.com/asafdavid23/eolctl/pkg/helm"
	helpers "github.com/asafdavid23/eolctl/pkg/helpers"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type ClusterReleaseInfo struct {
	Release      string `json:"release"`
	Namespace    string `json:"namespace"`
	Chart        string `json:"chart"`
	Product      string `json:"product"`
	Version      string `json:"version"`
	Eol          string `json:"eol"`
	Risk         string `json:"risk"`
	DaysUntilEOL int    `json:"days_until_eol,omitempty"`
}

// chartToSlug strips the trailing version suffix from a chart name.
// e.g., "cert-manager-v1.12.2" → "cert-manager", "nginx-ingress-4.7.1" → "nginx-ingress"
func chartToSlug(chart string) string {
	parts := strings.Split(chart, "-")
	for i, part := range parts {
		if len(part) > 0 && (part[0] == 'v' || (part[0] >= '0' && part[0] <= '9')) {
			if i > 0 {
				return strings.Join(parts[:i], "-")
			}
		}
	}
	return chart
}

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Scan installed Helm chart releases on a Kubernetes cluster for EOL information.",
	Long: `The 'cluster' command lists all Helm releases across all namespaces and checks
each chart's app version against the endoflife.date API to report EOL status and risk level.`,
	Run: func(cmd *cobra.Command, args []string) {
		logLevel, _ := cmd.Flags().GetString("log-level")
		logger := logging.NewLogger(logLevel)
		output, _ := cmd.Flags().GetString("output")

		logger.Debug("Listing Helm releases from cluster")
		releases, err := helm.ListReleases()
		if err != nil {
			logger.Fatalf("failed to list Helm releases: %v", err)
		}
		logger.Debugf("Found %d Helm release(s)", len(releases))

		if len(releases) == 0 {
			fmt.Println("No Helm releases found.")
			return
		}

		logger.Debug("Using Claude to map Helm charts to endoflife.date product slugs")
		stacks, err := ai.DetectHelmEOL(releases)
		if err != nil {
			logger.Fatalf("failed to detect Helm EOL info: %v", err)
		}
		logger.Debugf("Mapped %d chart(s) to product slugs", len(stacks))

		// Build a lookup from release name to the original HelmRelease for namespace/chart fields
		releaseByName := make(map[string]helm.HelmRelease, len(releases))
		for _, r := range releases {
			releaseByName[r.Name] = r
		}

		// Track which releases Claude already mapped; fill in the rest as fallbacks
		mappedNames := make(map[string]bool, len(stacks))
		for _, s := range stacks {
			mappedNames[s.ReleaseName] = true
		}
		for _, r := range releases {
			if !mappedNames[r.Name] {
				logger.Debugf("Claude did not map release %q — using chart name fallback", r.Name)
				stacks = append(stacks, ai.HelmStackInfo{
					ReleaseName: r.Name,
					Language:    chartToSlug(r.Chart),
					Version:     r.AppVersion,
				})
			}
		}

		var results []ClusterReleaseInfo

		for _, stack := range stacks {
			orig := releaseByName[stack.ReleaseName]

			var resultMap map[string]interface{}
			productData, err := helpers.GetProduct(stack.Language, stack.Version)
			if err != nil {
				// endoflife.date doesn't know this product — try ArtifactHub
				logger.Debugf("endoflife.date lookup failed for %s, trying ArtifactHub fallback", stack.Language)
				pkg, ahErr := artifacthub.SearchPackage(stack.Language)
				if ahErr != nil {
					logger.Debugf("ArtifactHub fallback also failed for %s: %v", stack.Language, ahErr)
					results = append(results, ClusterReleaseInfo{
						Release:   stack.ReleaseName,
						Namespace: orig.Namespace,
						Chart:     orig.Chart,
						Product:   stack.Language,
						Version:   stack.Version,
						Eol:       "unknown",
						Risk:      "UNKNOWN",
					})
					continue
				}
				riskLevel, eolLabel := artifacthub.RiskFromStaleness(orig.AppVersion, pkg.AppVersion, pkg.Deprecated)
				results = append(results, ClusterReleaseInfo{
					Release:   stack.ReleaseName,
					Namespace: orig.Namespace,
					Chart:     orig.Chart,
					Product:   stack.Language,
					Version:   stack.Version,
					Eol:       eolLabel,
					Risk:      riskLevel,
				})
				continue
			}

			if err := json.Unmarshal(productData, &resultMap); err != nil {
				logger.Errorf("failed to parse JSON for %s: %v", stack.Language, err)
				continue
			}

			riskInfo := helpers.CalculateRisk(resultMap["eol"])

			results = append(results, ClusterReleaseInfo{
				Release:      stack.ReleaseName,
				Namespace:    orig.Namespace,
				Chart:        orig.Chart,
				Product:      stack.Language,
				Version:      stack.Version,
				Eol:          helpers.GetStringValue(resultMap["eol"]),
				Risk:         string(riskInfo.Level),
				DaysUntilEOL: riskInfo.DaysUntilEOL,
			})
		}

		if output == "json" {
			jsonOutput, err := json.MarshalIndent(results, "", "  ")
			if err != nil {
				logger.Fatalf("failed to marshal results to JSON: %v", err)
			}
			fmt.Println(string(jsonOutput))
		} else if output == "table" {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Release", "Namespace", "Product", "Version", "EOL", "Risk"})
			table.SetAutoWrapText(false)

			for _, r := range results {
				renderRichRow(table, []string{r.Release, r.Namespace, r.Product, r.Version, r.Eol, r.Risk})
			}
			table.Render()
		} else {
			logger.Fatal("Invalid output type. Use 'table' or 'json'.")
		}

		riskReport, _ := cmd.Flags().GetBool("risk-report")
		suggestVersion, _ := cmd.Flags().GetBool("suggest-version")

		if riskReport || suggestVersion {
			var riskItems []ai.RiskItem
			var upgradeItems []ai.UpgradeItem

			for _, item := range results {
				riskItems = append(riskItems, ai.RiskItem{
					Product:      item.Product,
					Version:      item.Version,
					EOL:          item.Eol,
					RiskLevel:    item.Risk,
					DaysUntilEOL: item.DaysUntilEOL,
				})
				upgradeItems = append(upgradeItems, ai.UpgradeItem{
					Language:  item.Product,
					Version:   item.Version,
					EOL:       item.Eol,
					RiskLevel: item.Risk,
				})
			}

			if riskReport {
				printRiskNarrative(riskItems, logger)
			}
			if suggestVersion {
				printUpgradeSuggestions(upgradeItems, logger)
			}
		}
	},
}

func init() {
	// registered in scan.go
}
