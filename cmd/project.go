/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	ai "github.com/asafdavid23/eolctl/pkg/ai"
	helpers "github.com/asafdavid23/eolctl/pkg/helpers"

	"github.com/asafdavid23/eolctl/internal/logging"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type ProjectInfo struct {
	Product      string `json:"language"`
	Version      string `json:"version"`
	Eol          string `json:"eol"`
	Risk         string `json:"risk"`
	DaysUntilEOL int    `json:"days_until_eol,omitempty"`
}

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Identify and retrieve EOL information for a project based on its codebase.",
	Long: `The 'project' command analyzes the codebase in a specified project directory to identify the product and its version.
	It then retrieves End-of-Life (EOL) information for the identified product, providing you with up-to-date status and version details.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectDir := args[0]
		logLevel, _ := cmd.Flags().GetString("log-level")
		logger := logging.NewLogger(logLevel)
		output, _ := cmd.Flags().GetString("output")

		logger.Debug("Detecting project programming language")
		stacks, err := ai.DetectStack(projectDir)
		if err != nil {
			logger.Fatalf("failed to detect project stack: %v", err)
		}
		logger.Debugf("Detected %d stack(s)", len(stacks))

		var results []ProjectInfo

		for _, stack := range stacks {
			var result map[string]interface{}
			productData, err := helpers.GetProduct(stack.Language, stack.Version)
			if err != nil {
				logger.Errorf("failed to get product info for language %s and version %s: %v", stack.Language, stack.Version, err)
				continue
			}

			if err := json.Unmarshal(productData, &result); err != nil {
				logger.Errorf("failed to parse JSON response for %s: %v", stack.Language, err)
				continue
			}

			riskInfo := helpers.CalculateRisk(result["eol"])

			results = append(results, ProjectInfo{
				Product:      stack.Language,
				Version:      stack.Version,
				Eol:          helpers.GetStringValue(result["eol"]),
				Risk:         string(riskInfo.Level),
				DaysUntilEOL: riskInfo.DaysUntilEOL,
			})

			logger.Infof("Detected: Language=%s, Version=%s", stack.Language, stack.Version)
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
			table.SetHeader([]string{"Product", "Version", "Eol", "Risk"})

			for _, result := range results {
				table.Append([]string{result.Product, result.Version, result.Eol, result.Risk})
			}
			table.Render()
		} else {
			logger.Fatal("Invalid output type specified. Use 'table' or 'json'.")
		}

		riskReport, _ := cmd.Flags().GetBool("risk-report")

		var items []ai.RiskItem

		if riskReport {
			for _, item := range results {
				items = append(items, ai.RiskItem{
					Product:      item.Product,
					Version:      item.Version,
					EOL:          item.Eol,
					RiskLevel:    item.Risk,
					DaysUntilEOL: item.DaysUntilEOL,
				})
			}

			if len(items) == 0 {
				logger.Warn("no risk data to summarize — no components were successfully scanned")
			} else {
				narrative, err := ai.GenerateRiskNarrative(items)
				if err != nil {
					logger.Errorf("failed to generate risk narrative: %v", err)
				} else {
					fmt.Println("\n--- AI Risk Summary ---")
					fmt.Println(narrative)
				}
			}
		}
	},
}

func init() {
	// rootCmd.AddCommand(projectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	projectCmd.Flags().Bool("suggest-version", false, "Suggest a version upgrade based on the current project version")
	projectCmd.Flags().Bool("risk-report", false, "Generate an AI-powered risk narrative using Claude")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
