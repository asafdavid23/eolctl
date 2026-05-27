/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/asafdavid23/eolctl/internal/logging"
	ai "github.com/asafdavid23/eolctl/pkg/ai"
	helpers "github.com/asafdavid23/eolctl/pkg/helpers"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// productCmd represents the product command
var productCmd = &cobra.Command{
	Use:   "product",
	Short: "Retrieve End-of-Life (EOL) information for a specific product.",
	Long: `The 'product' command allows you to query the API for detailed End-of-Life (EOL) information about a specific product. 
By specifying the product name or ID, you can retrieve its EOL status, version information, and other relevant details.`,
	Run: func(cmd *cobra.Command, args []string) {

		initConfig()

		var enableCustomRange bool
		var outputData []byte

		name, _ := cmd.Flags().GetString("name")
		version, _ := cmd.Flags().GetString("version")
		outputFolder := viper.GetString("output.path")
		minVersion, _ := cmd.Flags().GetString("min")
		maxVersion, _ := cmd.Flags().GetString("max")
		output, _ := cmd.Flags().GetString("output")
		logLevel, _ := cmd.Flags().GetString("log-level")

		logger := logging.NewLogger(logLevel)

		if name != "" {
			logger.Debug("Fetching available products list from the API")
			availableProducts, err := helpers.GetAvailableProducts(output)

			if err != nil {
				logger.Fatalf("Failed to fetch available products from the API: %v", err)
			}

			logger.Debug("Verifying product does exist on the API")
			var products []string
			if err := json.Unmarshal(availableProducts, &products); err != nil {
				logger.Fatalf("failed to parse available products list: %v", err)
			}
			found := false
			for _, p := range products {
				if p == name {
					found = true
					break
				}
			}
			if !found {
				logger.Fatalf("%s doesn't exist on the API", name)
			}
		} else {
			logger.Fatal("Product name is required.")
		}

		logger.Debug("Fetching product data from the API")
		outputData, err := helpers.GetProduct(name, version)

		if err != nil {
			logger.Fatalf("Failed to fetch data for product %s: %v", name, err)
		}

		if minVersion != "" && maxVersion != "" {
			enableCustomRange = true
		}

		if enableCustomRange && version != "" {
			logger.Fatal("Custom range can't be run alongside with specific version")
		} else if enableCustomRange {
			logger.Debug("Custom range mode is enabled, fetching data from the API for product version from min to max")
			filtered, err := helpers.FilterVersions(outputData, minVersion, maxVersion)
			if err != nil {
				logger.Fatalf("Failed to filter versions: %v", err)
			}
			outputData = filtered
		}

		if outputFolder != "" {
			helpers.ExportToFile(outputData, outputFolder)
		}

		var result interface{}
		if err := json.Unmarshal(outputData, &result); err != nil {
			logger.Fatalf("failed to parse JSON response: %v", err)
		}

		if output == "table" {
			table := tablewriter.NewWriter(os.Stdout)
			switch v := result.(type) {
			case []interface{}:
				table.SetHeader([]string{"Cycle", "Latest", "LatestReleaseDate", "ReleaseDate", "LTS", "EOL", "Support", "Risk"})
				for _, item := range v {
					if release, ok := item.(map[string]interface{}); ok {
						renderRichRow(table, []string{
							helpers.GetStringValue(release["cycle"]),
							helpers.GetStringValue(release["latest"]),
							helpers.GetStringValue(release["latestReleaseDate"]),
							helpers.GetStringValue(release["releaseDate"]),
							helpers.GetStringValue(release["lts"]),
							helpers.GetStringValue(release["eol"]),
							helpers.GetStringValue(release["support"]),
							string(helpers.CalculateRisk(release["eol"]).Level),
						})
					}
				}
			case map[string]interface{}:
				table.SetHeader([]string{"Latest", "LatestReleaseDate", "ReleaseDate", "LTS", "EOL", "Support", "Risk"})
				renderRichRow(table, []string{
					helpers.GetStringValue(v["latest"]),
					helpers.GetStringValue(v["latestReleaseDate"]),
					helpers.GetStringValue(v["releaseDate"]),
					helpers.GetStringValue(v["lts"]),
					helpers.GetStringValue(v["eol"]),
					helpers.GetStringValue(v["support"]),
					string(helpers.CalculateRisk(v["eol"]).Level),
				})
			}
			table.Render()
		} else if output == "json" {
			fmt.Print(string(outputData))
		} else {
			logger.Fatal("output type is not valid.")
		}

		riskReport, _ := cmd.Flags().GetBool("risk-report")
		suggestVersion, _ := cmd.Flags().GetBool("suggest-version")

		if riskReport || suggestVersion {
			var riskItems []ai.RiskItem
			var upgradeItems []ai.UpgradeItem

			switch v := result.(type) {
			case map[string]interface{}:
				riskInfo := helpers.CalculateRisk(v["eol"])
				riskItems = append(riskItems, ai.RiskItem{
					Product:      name,
					Version:      version,
					EOL:          helpers.GetStringValue(v["eol"]),
					RiskLevel:    string(riskInfo.Level),
					DaysUntilEOL: riskInfo.DaysUntilEOL,
				})
				upgradeItems = append(upgradeItems, ai.UpgradeItem{
					Language:  name,
					Version:   version,
					EOL:       helpers.GetStringValue(v["eol"]),
					RiskLevel: string(riskInfo.Level),
				})
			case []interface{}:
				for _, item := range v {
					if cycle, ok := item.(map[string]interface{}); ok {
						riskInfo := helpers.CalculateRisk(cycle["eol"])
						riskItems = append(riskItems, ai.RiskItem{
							Product:      name,
							Version:      helpers.GetStringValue(cycle["cycle"]),
							EOL:          helpers.GetStringValue(cycle["eol"]),
							RiskLevel:    string(riskInfo.Level),
							DaysUntilEOL: riskInfo.DaysUntilEOL,
						})
						upgradeItems = append(upgradeItems, ai.UpgradeItem{
							Language:  name,
							Version:   helpers.GetStringValue(cycle["cycle"]),
							EOL:       helpers.GetStringValue(cycle["eol"]),
							RiskLevel: string(riskInfo.Level),
						})
					}
				}
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
	// rootCmd.AddCommand(productCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// productCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	productCmd.Flags().StringP("name", "n", "", "Name of the product")
	productCmd.Flags().StringP("version", "v", "", "Version of the product")
	productCmd.Flags().String("min", "", "Minimum version to query")
	productCmd.Flags().String("max", "", "Maximum version to query")
}
