/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"eolctl/internal"
	"eolctl/internal/logging"
	"fmt"
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
		var customRangeOutput []byte
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
			availbleProducts, err := helpers.GetAvailableProducts()

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
		outputData, err := helpers.GetProduct(name, version)

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
			customRangeOutput, _ = helpers.FilterVersions(outputData, minVersion, maxVersion)

			if customRangeOutput != nil && outputFolder != "" {
				helpers.ExportToFile(customRangeOutput, outputFolder)
				os.Exit(0)
			} else if output != "" {
				helpers.ConvertOutput(customRangeOutput, output)
				os.Exit(0)
			} else {
				fmt.Print(string(customRangeOutput))
				os.Exit(0)
			}
		}

		if outputFolder != "" && output == "" {
			helpers.ExportToFile(outputData, outputFolder)
			os.Exit(0)
		}

		if output != "" && outputFolder == "" {
			helpers.ConvertOutput(outputData, output)
			os.Exit(0)
		} else {
			fmt.Print(string(outputData))
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
	productCmd.Flags().String("existing-version", "", "Existing version to compare")
	productCmd.Flags().String("future-version", "", "Future version to compare")
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
