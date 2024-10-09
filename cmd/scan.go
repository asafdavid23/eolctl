/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan your code project / dockerfile for End-of-Life (EOL) information.",
	Long: `The 'scan' command analyzes your code project to identify its programming language and retrieves the relevant End-of-Life (EOL) information. 
	This tool helps ensure that you are aware of the support status and lifecycle for the programming languages and frameworks used in your project, enabling better maintenance and planning.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.AddCommand(projectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
