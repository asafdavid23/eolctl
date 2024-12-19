/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	helpers "eolctl/internal"
	"eolctl/internal/logging"
	"eolctl/internal/scanner"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type ProjectInfo struct {
	Product string `json:"langugage"`
	Version string `json:"version"`
	Eol     string `string:"eol"`
}

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Identify and retrieve EOL information for a project based on its codebase.",
	Long: `The 'project' command analyzes the codebase in a specified project directory to identify the product and its version. 
	It then retrieves End-of-Life (EOL) information for the identified product, providing you with up-to-date status and version details.`,
	Run: func(cmd *cobra.Command, args []string) {
		// var outputData []byte

		projectDir := args[0]
		logLevel, _ := cmd.Flags().GetString("log-level")

		logger := logging.NewLogger(logLevel)
		output, _ := cmd.Flags().GetString("output")

		if len(projectDir) < 1 {
			logger.Fatal("please specify project dir")
		}

		logger.Debug("Detecting project programming language")
		languages, projects, err := scanner.DetectLanguages(projectDir)

		if err != nil {
			logger.Fatal(err)
		}

		logger.Debugf("Project language is: %s and projects are: %s", languages, projects)

		var results []ProjectInfo

		for projIndex, project := range projects {
			language := languages[projIndex]
			version, err := scanner.DetectVersion(project)
			if err != nil {
				logger.Errorf("failed to detect version for project %s: %v", project, err)
				continue
			}

			var result map[string]interface{}
			productData, err := helpers.GetProduct(language, version)

			if err := json.Unmarshal(productData, &result); err != nil {
				logger.Fatalf("faild to parse JSON response: %v", err)
			}

			if err != nil {
				logger.Errorf("failed to get proudct info for language %s and version %s: %v", language, version, err)
				continue
			}

			product := string(productData)

			results = append(results, ProjectInfo{
				Product: language,
				Version: version,
				Eol:     helpers.GetStringValue(result["eol"]),
			})

			logger.Infof("Detected: Language=%s, Version=%s, Product=%s", language, version, product)
		}

		// Handle outputs
		if output == "json" {
			// Marshal the results to JSON
			jsonOutput, err := json.MarshalIndent(results, "", "  ")
			if err != nil {
				logger.Fatalf("Failed to marshal results to JSON: %v", err)
			}
			fmt.Println(string(jsonOutput))
		} else if output == "table" {
			// Print as a table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Product", "Version", "Eol"})

			for _, result := range results {
				table.Append([]string{result.Product, result.Version, result.Eol})
			}
			table.Render()
		} else {
			logger.Fatal("Invalid output type specified. Use 'table' or 'json'.")
		}
	},
}

func init() {
	// rootCmd.AddCommand(projectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	projectCmd.Flags().Bool("suggest-version", false, "Suggest a version upgrade based on the current project version")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
