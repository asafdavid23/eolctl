/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"eolctl/internal"
	"github.com/spf13/cobra"
)

// availableProductsCmd represents the availableProducts command
var availableProductsCmd = &cobra.Command{
	Use:   "available-products",
	Short: "Get all available products supported by the API",
	Long:  `Using this command you can look and filter for your relevant product API support.`,
	Run: func(cmd *cobra.Command, args []string) {
		helpers.GetAvailableProducts()
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
