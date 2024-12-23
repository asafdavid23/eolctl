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
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		var result interface{}

		name, _ := cmd.Flags().GetString("name")
		version, _ := cmd.Flags().GetString("version")
		outputFolder := viper.GetString("output.path")
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

		if minVersion != "" && maxVersion != "" {
			enableCustomRange = true
		}

		if enableCustomRange && version != "" {
			logger.Fatal("Custom range can't be run alongside with specific version")
		}

		c, err := localCache.InitializeCacheFile()

		if err != nil {
			logger.Fatalf("Failed to initialize cache file: %v", err)
		}

		if name != "" && version == "" {
			productCacheKey := fmt.Sprintf("product-%s", name)
			if cacheData, found := c.Get(productCacheKey); found {
				logger.Info("Cache hit for product")
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

				if err := json.Unmarshal(cacheDataBytes, &outputData); err != nil {
					logger.Fatalf("Failed to parse JSON response: %v", err)
				}
			} else {
				logger.Info("Cache miss for product")
				logger.Info("Fetching product data from the API")

				if enableCustomRange {
					logger.Debug("Custom range mode is enabled, fetching data from the API for product version from min to max")
					outputData, _ = helpers.FilterVersions(outputData, minVersion, maxVersion)
				} else {
					outputData, err = helpers.GetProduct(name, version)
				}

				if err != nil {
					logger.Fatalf("Failed to fetch data for proudct %s\n\n", name)
				}

				if err := json.Unmarshal(outputData, &result); err != nil {
					logger.Fatalf("Failed to parse JSON response: %v", err)
				}

				logger.Debug("Caching product")
				c.Set(productCacheKey, outputData, cache.DefaultExpiration)
				localCache.SaveCacheFile()
			}
		} else if name != "" && version != "" {
			productCycleCacheKey := fmt.Sprintf("product-%s-%s", name, version)

			if cacheData, found := c.Get(productCycleCacheKey); found {
				logger.Info("Cache hit for product cycle")
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

				if err := json.Unmarshal(cacheDataBytes, &result); err != nil {
					logger.Fatalf("Failed to parse JSON response: %v", err)
				}
			} else {
				logger.Info("Cache miss for product cycle")
				logger.Info("Fetching product cycle data from the API")

				outputData, err = helpers.GetProduct(name, version)

				if err != nil {
					logger.Fatalf("Failed to fetch data for proudct %s\n\n", name)
				}

				if err := json.Unmarshal(outputData, &result); err != nil {
					logger.Fatalf("Failed to parse JSON response: %v", err)
				}

				logger.Debug("Caching product cycle")
				c.Set(productCycleCacheKey, outputData, cache.DefaultExpiration)
				localCache.SaveCacheFile()
			}
		}

		if output == "table" {
			var headers []string
			table := tablewriter.NewWriter(os.Stdout)

			switch v := result.(type) {
			case []interface{}:

				if len(v) > 0 {
					// Get headers from the first item
					if firstItem, ok := v[0].(map[string]interface{}); ok {

						for key := range firstItem {
							headers = append(headers, key)
						}

						table.SetHeader(headers)
					}
				}

				for _, item := range v {
					if record, ok := item.(map[string]interface{}); ok {
						row := []string{}

						for _, key := range headers {
							row = append(row, helpers.GetStringValue(record[key]))
						}
						table.Append(row)
					}
				}

				table.Render()
			case map[string]interface{}:
				// Handle single object
				headers := []string{}
				row := []string{}

				for key, value := range v {
					headers = append(headers, key)
					row = append(row, helpers.GetStringValue(value))
				}

				table.SetHeader(headers)
				table.Append(row)
				table.Render()
			}

		} else if output == "json" {
			productsJSON, err := json.Marshal(result)

			if err != nil {
				logger.Fatalf("Failed to marshal products to JSON: %v", err)
			}

			fmt.Print(string(productsJSON))
		} else {
			logger.Fatal("output type is not valid.")
		}

		if outputFolder != "" {
			productsJSON, err := json.Marshal(result)

			if err != nil {
				logger.Fatalf("Failed to marshal products to JSON: %v", err)
			}

			helpers.ExportToFile(productsJSON, outputFolder)
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
