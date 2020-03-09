package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
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
	if len(args) == 0 {
		fmt.Println("Please enter a URL for a GitHub user. Example: `gobackup backup github.com/tempor1s`")
		return
	}

	// Get the URL to clone
	repoURL := args[0]

	// Get the users github name for the directory
	dirName := path.Base(repoURL)

	// Make sure it has the https prefix, and add it if it does not
	if !strings.HasPrefix(repoURL, "https://") {
		repoURL = "https://" + repoURL
	}

	// Create username dir to put all cloned repos in.
	err := os.MkdirAll(dirName, os.ModePerm)

	if err != nil {
		log.Fatal(err)
	}

	// Clone the repo into username folder that we just created - will keep git information because not bare
	_, err = git.PlainClone(dirName, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})

	if err != nil {
		log.Fatal(err)
	}
}
