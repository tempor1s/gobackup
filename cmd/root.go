package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// If you want to have explicit logging when running commands.
var Verbose bool

var rootCmd = &cobra.Command{
	Use:   "gobackup",
	Short: "A simple CLI tool that backs up all your github repos to your local machine, and uploads them to another repo host like GitLab!",
	Long:  "GoBackup will allow you to easily backup, or download, all of your repos onto your local machine very quickly. It can also be used to transfer these repos into another cloud hosting platform like GitLab or BitBucket if you want to keep a backup there, or switch services!",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Usage(); err != nil {
			log.Fatal(err)
		}
	},
}

// Execute a command
func Execute() {
	// Global Flags
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose output mode")

	// Execute Command
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
