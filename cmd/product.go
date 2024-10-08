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

var versionsToCompare []string

// productCmd represents the product command
var productCmd = &cobra.Command{
	Use:   "product",
	Short: "Query for specific product EOL information.",
	Long:  `This command will specify the exact product to qury data from the API.`,
	Run: func(cmd *cobra.Command, args []string) {

		initConfig()

		name, _ := cmd.Flags().GetString("name")
		version, _ := cmd.Flags().GetString("version")
		outputFolder, _ := cmd.Flags().GetString("output")
		minVersion, _ := cmd.Flags().GetString("min")
		maxVersion, _ := cmd.Flags().GetString("max")
		version1, _ := cmd.Flags().GetString("existing-version")
		version2, _ := cmd.Flags().GetString("future-version")

		if name == "" {
			log.Fatal("Product name is required.")
		}

		cycle1 := helpers.GetProduct(name, version1, outputFolder, minVersion, maxVersion)
		cycle2 := helpers.GetProduct(name, version2, outputFolder, minVersion, maxVersion)

		if version1 != "" && version2 != "" {
			cycle1, cycle2 := helpers.CompareTwoVersions(cycle1, cycle2)
			fmt.Printf("%s\n", string(cycle1))
			fmt.Printf("%s\n", string(cycle2))
		}

		helpers.GetProduct(name, version, outputFolder, minVersion, maxVersion)
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
