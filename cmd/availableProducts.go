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
	"os"
)

// availableProductsCmd represents the availableProducts command
var availableProductsCmd = &cobra.Command{
	Use:   "available-products",
	Short: "List all products supported by the API.",
	Long: `The 'available-products' command retrieves and displays a list of all products currently supported by the API. 
You can filter the list to find relevant products that meet your specific needs, allowing you to quickly identify which products are available for interaction with the API.`,
	Run: func(cmd *cobra.Command, args []string) {
		logLevel, _ := cmd.Flags().GetString("log-level")
		output, _ := cmd.Flags().GetString("output")

		logger := logging.NewLogger(logLevel)
		outputData, err := helpers.GetAvailableProducts(output)

		if err != nil {
			logger.Fatalf("Failed to fetch available products from the API: %v", err)
		}

		var products []interface{}
		if err := json.Unmarshal(outputData, &products); err != nil {
			logger.Fatalf("faild to parse JSON response: %d", err)
		}

		if output == "table" {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Product"})

			for _, product := range products {
				if str, ok := product.(string); ok {
					table.Append([]string{str})
				}
			}

			table.Render()
		} else if output == "json" {
			fmt.Print(string(outputData))
		} else {
			logger.Fatal("Output type is not valid.")
		}
	},
}

func init() {
	// rootCmd.AddCommand(availableProductsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// availableProductsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// availableProductsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
