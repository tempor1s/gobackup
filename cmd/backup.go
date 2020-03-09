package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(backupCommand)
}

var backupCommand = &cobra.Command{
	Use:   "backup [github.com/username]",
	Short: "Backup all of the repos at a given github URL",
	Long:  "This command will backup all of the repos from a user into your CURRENT directory.",
	// Run the backup command
	Run: backup,
}

// TODO: Put this into another file as to not muddy up the cmd package that should only be used for managing commands
func backup(cmd *cobra.Command, args []string) {
	fmt.Println("Hello, from the backup command!")
}
