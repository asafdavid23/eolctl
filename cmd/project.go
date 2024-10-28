/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"eolctl/internal"
	"eolctl/internal/logging"
	"eolctl/internal/scanner"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Identify and retrieve EOL information for a project based on its codebase.",
	Long: `The 'project' command analyzes the codebase in a specified project directory to identify the product and its version. 
It then retrieves End-of-Life (EOL) information for the identified product, providing you with up-to-date status and version details.`,
	Run: func(cmd *cobra.Command, args []string) {
		var outputData []byte

		projectDir := args[0]
		logLevel, _ := cmd.Flags().GetString("log-level")

		logger := logging.NewLogger(logLevel)
		output, _ := cmd.Flags().GetString("output")

		logger.Debug("Detecting project programming language")
		language, err := scanner.DetectLanguage(projectDir)
		logger.Debugf("Project language is %s: ", language)

		if err != nil {
			logger.Fatal(err)
		}

		logger.Debug("Detecting package file")
		packageFile, err := scanner.DetectPackgesFile(projectDir)

		if err != nil {
			logger.Fatal(err)
		}

		if language == "JavaScript" {
			logger.Debug("Detect version from pacakge.json")
			version, err := scanner.DetectVersionFromPackageJSON(packageFile)

			if err != nil {
				logger.Fatal(err)
			}

			language = "nodejs"
			parts := strings.Split(version, ".")
			shortVersion := parts[0]

			logger.Debug("Fetching project product version from the API")
			outputData, err = helpers.GetProduct(language, shortVersion, output)

			if err != nil {
				logger.Fatal(err)
			}

		} else if language == "Python" {
			version, err := scanner.DetectPythonVersion(projectDir)

			if err != nil {
				logger.Fatal(err)
			}

			outputData, err = helpers.GetProduct(language, version, output)

			if err != nil {
				logger.Fatal(err)
			}

		} else if language == "Go" {
			logger.Debug("Detect version from go.mod")
			version, err := scanner.DetectVersionFromGoMod(packageFile)

			if err != nil {
				logger.Fatal(err)
			}

			parts := strings.Split(version, ".")
			shortVersion := parts[0] + "." + parts[1]

			logger.Debug("Fetching project product version from the API")
			outputData, err = helpers.GetProduct(strings.ToLower(language), shortVersion, output)

			if err != nil {
				logger.Fatal(err)
			}
		}

		var result map[string]interface{}
		if err := json.Unmarshal(outputData, &result); err != nil {
			logger.Fatalf("faild to parse JSON response: %v", err)
		}

		if output == "table" {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Latest", "LatestReleaseDate", "ReleaseDate", "LTS", "EOL", "SUPPORT"})

			row := []string{
				helpers.GetStringValue(result["latest"]),
				helpers.GetStringValue(result["latestReleaseDate"]),
				helpers.GetStringValue(result["releaseDate"]),
				helpers.GetStringValue(result["lts"]),
				helpers.GetStringValue(result["eol"]),
				helpers.GetStringValue(result["support"]),
			}
			table.Append(row)
			table.Render()
		} else if output == "json" {
			fmt.Print(string(outputData))
		} else {
			logger.Fatal("output type is not valid.")
		}
	},
}

func init() {
	// rootCmd.AddCommand(projectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
