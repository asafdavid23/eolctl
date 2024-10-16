/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"eolctl/internal"
	"fmt"
	"github.com/spf13/cobra"
)

// availableProductsCmd represents the availableProducts command
var availableProductsCmd = &cobra.Command{
	Use:   "available-products",
	Short: "List all products supported by the API.",
	Long: `The 'available-products' command retrieves and displays a list of all products currently supported by the API. 
You can filter the list to find relevant products that meet your specific needs, allowing you to quickly identify which products are available for interaction with the API.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := helpers.NewLogger()
		outputData, err := helpers.GetAvailableProducts()

		if err != nil {
			logger.Fatalf("Failed to fetch available products from the API: %v", err)
		}

		fmt.Println(string(outputData))
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
