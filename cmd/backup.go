package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func init() {
	rootCmd.AddCommand(backupCommand)
	backupCommand.Flags().StringVarP(&Token, "token", "t", "", "Your personal access token. Needed to be able to clone private repos.")
}

var Token string

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

	repos := getRepos()

	// Clone all repos
	cloneRepos(repos, dirName)
}

func getRepos() []string {
	if Token == "" {
		// TODO: Support non-access token based request
		fmt.Println("WARNING: Personal Token was not passed in. Please pass in a token using --token - Support for NON-Token download coming soon.")
		return []string{}
	}

	// Set up OAuth Token stuff
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: Token},
	)

	tc := oauth2.NewClient(ctx, ts)

	// Create a new github client using the OAuth2 token
	client := github.NewClient(tc)

	// Get all repos that the user owns
	opt := &github.RepositoryListOptions{
		Affiliation: "owner",
		ListOptions: github.ListOptions{
			PerPage: 100000,
		},
	}

	repos, _, err := client.Repositories.List(ctx, "", opt)

	if err != nil {
		log.Fatal(err)
	}

	var ret []string
	for _, repo := range repos {
		ret = append(ret, *repo.HTMLURL)
	}

	return ret
}

func cloneRepos(repos []string, dirName string) {

	// Create username dir to put all cloned repos in.
	err := os.MkdirAll(dirName, os.ModePerm)

	if err != nil {
		log.Fatal(err)
	}

	for _, repo := range repos {
		fmt.Printf("[gobackup] cloning %s\n", repo)
		repoName := path.Base(repo)
		_, err := git.PlainClone(dirName+"/"+repoName, false, &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: "gobackup",
				Password: Token,
			},
			URL:      repo,
			Progress: os.Stdout,
		})

		if err != nil {
			log.Fatal(err)
		}
	}
}
