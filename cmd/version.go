package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCommand)
}

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print the version of GoBackup",
	Long:  "All software has a version, and this is GoBackup's.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GoBackup's Version: Unreleased")
	},
}
