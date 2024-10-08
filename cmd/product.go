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
)

// productCmd represents the product command
var productCmd = &cobra.Command{
	Use:   "product",
	Short: "Query for specific product EOL information.",
	Long:  `This command will specify the exact product to qury data from the API.`,
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()

		var enableCustomRange bool
		var customRangeOutput []byte

		name, _ := cmd.Flags().GetString("name")
		version, _ := cmd.Flags().GetString("version")
		outputFolder, _ := cmd.Flags().GetString("output")
		minVersion, _ := cmd.Flags().GetString("min")
		maxVersion, _ := cmd.Flags().GetString("max")

		if name == "" {
			log.Fatal("Product name is required.")
		}

		outputData := helpers.GetProduct(name, version)

		if minVersion != "" && maxVersion != "" {
			enableCustomRange = true
		}

		if enableCustomRange && version != "" {
			log.Fatal("Custom range can't be run alongside with specific version")
		} else if enableCustomRange {
			log.Print("Executing custom range")
			customRangeOutput, _ = helpers.FilterVersions(outputData, minVersion, maxVersion)

			if outputFolder != "" {
				helpers.ExportToFile(customRangeOutput, outputFolder)
			} else {
				fmt.Println(string(customRangeOutput))
			}
		}

		if outputFolder != "" {
			helpers.ExportToFile(outputData, outputFolder)
		} else {
			fmt.Println(string(outputData))
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
	productCmd.Flags().StringP("output", "o", "", "Export to file")
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
