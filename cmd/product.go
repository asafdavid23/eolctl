/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"endoflifectl/internal"
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

		name, _ := cmd.Flags().GetString("name")
		version, _ := cmd.Flags().GetString("version")
		export, _ := cmd.Flags().GetBool("export")
		outputFolder := viper.GetString("app.outputFolder")

		if name == "" {
			log.Fatal("Product name is required.")

		}

		outputData := helpers.GetProduct(name, version)

		if export {
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
	productCmd.Flags().BoolP("export", "e", false, "Export to file")
	// productCmd.Flags().Bool("top", false, "Get top 3 only")
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
		log.Fatalf("fatal error config file: default %v", err)
	}
}
