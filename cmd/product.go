/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"eolctl/internal"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
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
		version1, _ := cmd.Flags().GetString("existing-version")
		version2, _ := cmd.Flags().GetString("future-version")
		output, _ := cmd.Flags().GetString("output")

		if name == "" {
			log.Fatal("Product name is required.")
		}

		outputData = helpers.GetProduct(name, version)

		if minVersion != "" && maxVersion != "" {
			enableCustomRange = true
		}

		if enableCustomRange && version != "" {
			log.Fatal("Custom range can't be run alongside with specific version")
		} else if enableCustomRange {
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

		if version1 != "" && version2 != "" {
			cycle1 := helpers.GetProduct(name, version1)
			cycle2 := helpers.GetProduct(name, version2)

			cycle1, cycle2 = helpers.CompareTwoVersions(cycle1, cycle2)
			fmt.Printf("%s\n", string(cycle1))
			fmt.Printf("%s\n", string(cycle2))
			os.Exit(0)
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
	productCmd.Flags().StringP("output", "o", "", "Output type table/json/yaml")
	productCmd.Flags().String("output-path", "", "Export to file")
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
