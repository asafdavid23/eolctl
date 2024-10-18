/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"eolctl/internal/logging"
	"eolctl/internal/scanner"
	"fmt"
	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Identify and retrieve EOL information for a project based on its codebase.",
	Long: `The 'project' command analyzes the codebase in a specified project directory to identify the product and its version. 
It then retrieves End-of-Life (EOL) information for the identified product, providing you with up-to-date status and version details.`,
	Run: func(cmd *cobra.Command, args []string) {

		projectDir := args[0]
		logger := logging.NewLogger()
		// product, productFile := helpers.IdentifyProduct(projectDir)
		// version := helpers.IdentifyProductVersion(product, projectDir, productFile)

		// helpers.GetProduct(product, version)
		language, err := scanner.DetectLanguage(projectDir)

		if err != nil {
			logger.Fatal(err)
		}

		fmt.Println(language)
		packageFile, err := scanner.DetectPackgesFile(projectDir)

		if err != nil {
			logger.Fatal(err)
		}

		if language == "JavaScript" {
			version, err := scanner.DetectVersionFromPackageJSON(packageFile)

			if err != nil {
				logger.Fatal(err)
			}

			fmt.Print(version)
		} else if language == "Python" {
			version, err := scanner.DetectVersionFromRequirementsTxt(packageFile)

			if err != nil {
				logger.Fatal(err)
			}

			fmt.Print(version)
		} else if language == "Go" {
			version, err := scanner.DetectVersionFromGoMod(packageFile)

			if err != nil {
				logger.Fatal(err)
			}

			fmt.Print(version)
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
