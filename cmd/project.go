/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"endoflifectl/internal"
	"fmt"
	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		projectDir := args[0]
		product, productFile := helpers.IdentifyProduct(projectDir)
		version := helpers.IdentifyProductVersion(product, projectDir, productFile)

		data := helpers.GetProduct(product, version)
		fmt.Println(string(data))
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
