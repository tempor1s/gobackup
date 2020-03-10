package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tempor1s/gobackup/backup"
)

func init() {
	rootCmd.AddCommand(backupCommand)
	backupCommand.Flags().StringVarP(&Token, "token", "t", "", "Your personal access token. Needed to be able to clone private repos.")
}

// Token is a github personal token passed in with --token
var Token string

// backupCommand is the command register for backing up repos
var backupCommand = &cobra.Command{
	Use:   "backup [github.com/username]",
	Short: "Backup all of the repos at a given github URL",
	Long:  "This command will backup all of the repos from a user into your CURRENT directory.",
	// Run the backup command
	Run: backupCmd,
}

// backupCmd just allows us to hook into our backup module
func backupCmd(cmd *cobra.Command, args []string) {
	backup.GitHub(Token, args)
}
