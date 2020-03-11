package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tempor1s/gobackup/download"
)

func init() {
	rootCmd.AddCommand(downloadCommand)
	downloadCommand.Flags().StringVarP(&Token, "token", "t", "", "Your personal access token. Needed to be able to clone private repos.")
}

// downloadCommand is the command register for backing up repos
var downloadCommand = &cobra.Command{
	Use:   "download [platform] [username]",
	Short: "Download all the repos for a given user.",
	Long:  "This command will download all of the repos from a user into a <username>/ repository. Works with Github/Gitlab currently.",
	// Run the backup command
	Run: downloadCmd,
}

// downloadCmd just allows us to hook into our backup module
func downloadCmd(cmd *cobra.Command, args []string) {
	download.Start(Token, args)
}
