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
	Run:   backup,
}

func backup(cmd *cobra.Command, args []string) {
	fmt.Println("Hello, from the backup command!")
}
