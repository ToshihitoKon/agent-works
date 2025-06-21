package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gdrive-props",
	Short: "Google Drive Custom Properties Manager",
	Long: `A CLI tool to manage Google Drive file custom properties (appProperties).
This tool allows you to add, update, delete, and search files based on custom properties.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add subcommands here
}