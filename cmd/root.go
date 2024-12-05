/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var Version = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "eolctl",
	Short: "Access End-of-Life (EOL) dates and support lifecycles for various products.",
	Long: `The 'eolctl' command-line tool provides users with comprehensive information on End-of-Life (EOL) dates and support lifecycles for a wide range of products. 
It aggregates data from various reliable sources, presenting it in a clear and concise manner. 
Additionally, the tool offers an easily accessible API for data retrieval and supports iCalendar format for seamless integration into your scheduling applications.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
			fmt.Println("eolctl version:", Version)
			os.Exit(0)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.eolctl.yaml)")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Display the version of this CLI tool")
	rootCmd.PersistentFlags().StringP("output", "o", "table", "Output type table/json/yaml")
	rootCmd.PersistentFlags().String("log-level", "", "Set log level")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		// Use the specified config file
		viper.SetConfigFile(cfgFile)
	} else {
		// Use the default config file
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config/")
	}

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		// Ignore "file not found" errors; log other errors
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Error reading config file: %v", err)
		}
	}

	// Log the config file being used
	if viper.ConfigFileUsed() != "" {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}
}
