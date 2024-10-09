/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// compareCmd represents the compare command
var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare two cycles retrieved from the API.",
	Long: `The 'compare' command allows you to compare two cycles retrieved from the API. 
You can use this command to identify differences, similarities, or perform detailed comparisons between the cycles, 
providing insights into their attributes or configurations. 
Examples of usage might include comparing versions, configurations, or other attributes of resources.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(compareCmd)
	compareCmd.AddCommand(productCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// compareCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// compareCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
