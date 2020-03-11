package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tempor1s/gobackup/upload"
)

func init() {
	rootCmd.AddCommand(uploadCommand)
	downloadCommand.Flags().StringVarP(&Token, "token", "t", "", "Your personal access token. Needed to be able to upload to a repository...")
}

// uploadCommand is the command register for backing up repos
var uploadCommand = &cobra.Command{
	Use:   "upload [platform] [username]",
	Short: "Upload all the repos in the current folder to the platform of your choice.",
	Long:  "This command will upload all of the repositories that are in your current directory to an account on another platform. Make sure you know what you are doing before you use this!",
	// Run the backup command
	Run: uploadCmd,
}

func uploadCmd(cmd *cobra.Command, args []string) {
	upload.Start(Token, args)
}
