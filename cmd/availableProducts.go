/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	helpers "eolctl/internal"
	localCache "eolctl/internal/cache"
	"eolctl/internal/logging"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/cobra"
)

// availableProductsCmd represents the availableProducts command
var availableProductsCmd = &cobra.Command{
	Use:   "available-products",
	Short: "List all products supported by the API.",
	Long: `The 'available-products' command retrieves and displays a list of all products currently supported by the API. 
You can filter the list to find relevant products that meet your specific needs, allowing you to quickly identify which products are available for interaction with the API.`,
	Run: func(cmd *cobra.Command, args []string) {
		var outputData []byte
		var err error
		var products []json.RawMessage

		logLevel, _ := cmd.Flags().GetString("log-level")
		output, _ := cmd.Flags().GetString("output")

		logger := logging.NewLogger(logLevel)

		c := localCache.InitializeCacheFile()

		cacheKey := "available-products"

		if cacheData, found := c.Get(cacheKey); found {
			logger.Info("Cache hit for available products")
			logger.Debugf("Type of cacheData: %T", cacheData)

			// Assert cacheData to cache.Item
			cacheItem, ok := cacheData.(cache.Item)
			if !ok {
				logger.Fatalf("Failed to assert cache data to cache.Item")
			}

			// Access the data within cache.Item
			cacheDataBytes, ok := cacheItem.Object.([]byte)
			if !ok {
				logger.Fatalf("Failed to assert cache item object to []byte")
			}

			// Parse API response into products slice
			if err := json.Unmarshal(cacheDataBytes, &products); err != nil {
				logger.Fatalf("Failed to parse JSON response: %v", err)
			}

		} else {
			logger.Info("Cache miss for avaiable products")
			logger.Info("Fetching available products from the API")
			outputData, err = helpers.GetAvailableProducts(output)

			if err != nil {
				logger.Fatalf("Failed to fetch available products from the API: %v", err)
			}

			// Parse API response into products slice
			if err := json.Unmarshal(outputData, &products); err != nil {
				logger.Fatalf("Failed to parse JSON response: %v", err)
			}

			logger.Debug("Caching available products")
			c.Set(cacheKey, outputData, cache.DefaultExpiration)
			localCache.SaveCacheFile()
		}

		if output == "table" {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Product"})

			for _, product := range products {
				productStr := string(product)
				table.Append([]string{productStr})
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
