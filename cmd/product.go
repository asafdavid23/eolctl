/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"eolctl/internal"
	"eolctl/internal/logging"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
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
		outputFolder, _ := cmd.Flags().GetString("output-path")
		minVersion, _ := cmd.Flags().GetString("min")
		maxVersion, _ := cmd.Flags().GetString("max")
		output, _ := cmd.Flags().GetString("output")
		logLevel, _ := cmd.Flags().GetString("log-level")

		logger := logging.NewLogger(logLevel)

		if name != "" {
			logger.Debug("Fetching available products list from the API")
			availbleProducts, err := helpers.GetAvailableProducts(output)

			if err != nil {
				logger.Fatalf("Failed to fetch available products from the API: %v", err)
			}

			logger.Debug("Verifying product does exist on the API")
			if !strings.Contains(string(availbleProducts), name) {
				logger.Fatalf("%s doesn't exists on the API", name)
			}
		} else {
			logger.Fatal("Product name is required.")
		}

		logger.Debug("Fetching product data from the API")
		outputData, err := helpers.GetProduct(name, version, output)

		if err != nil {
			logger.Fatalf("Failed to fetch data for proudct %s\n\n", name)
		}

		if minVersion != "" && maxVersion != "" {
			enableCustomRange = true
		}

		if enableCustomRange && version != "" {
			logger.Fatal("Custom range can't be run alongside with specific version")
		} else if enableCustomRange {
			logger.Debug("Custom range mode is enabled, fetching data from the API for product version from min to max")
			outputData, _ = helpers.FilterVersions(outputData, minVersion, maxVersion)
		}

		if outputFolder != "" && output == "" {
			helpers.ExportToFile(outputData, outputFolder)
			os.Exit(0)
		}

		var result interface{}
		if err := json.Unmarshal(outputData, &result); err != nil {
			logger.Fatalf("faild to parse JSON response: %v", err)
		}

		if output == "table" {
			table := tablewriter.NewWriter(os.Stdout)
			switch v := result.(type) {
			case []interface{}:
				table.SetHeader([]string{"Cycle", "Latest", "LatestReleaseDate", "ReleaseDate", "LTS", "EOL", "SUPPORT"})
				for _, item := range v {
					if release, ok := item.(map[string]interface{}); ok {
						row := []string{
							helpers.GetStringValue(release["cycle"]),
							helpers.GetStringValue(release["latest"]),
							helpers.GetStringValue(release["latestReleaseDate"]),
							helpers.GetStringValue(release["releaseDate"]),
							helpers.GetStringValue(release["lts"]),
							helpers.GetStringValue(release["eol"]),
							helpers.GetStringValue(release["support"]),
						}
						table.Append(row)
					}
				}
			case map[string]interface{}:
				table.SetHeader([]string{"Latest", "LatestReleaseDate", "ReleaseDate", "LTS", "EOL", "SUPPORT"})

				row := []string{
					helpers.GetStringValue(v["latest"]),
					helpers.GetStringValue(v["latestReleaseDate"]),
					helpers.GetStringValue(v["releaseDate"]),
					helpers.GetStringValue(v["lts"]),
					helpers.GetStringValue(v["eol"]),
					helpers.GetStringValue(v["support"]),
				}
				table.Append(row)
			}
			table.Render()
		} else if output == "json" {
			fmt.Print(string(outputData))
		} else {
			logger.Fatal("output type is not valid.")
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

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config/")
	}

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
	}
}
