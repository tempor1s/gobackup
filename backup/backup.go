package backup

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func GitHub(token string, args []string) {
	if len(args) == 0 {
		fmt.Println("Please enter a URL for a GitHub user. Example: `gobackup backup github.com/tempor1s`")
		return
	}

	if token == "" {
		// TODO: Support non-access token based request
		fmt.Println("WARNING: Personal token was not passed in. Please pass in a token using --token - Support for NON-token download coming soon.")
		return
	}

	// Get the URL to clone
	repoURL := args[0]

	// Get the users github name for the directory
	dirName := path.Base(repoURL)

	repos := getRepos(token)

	// Clone all repos
	cloneRepos(repos, dirName, token)
}

func getRepos(token string) []string {

	// Set up OAuth token stuff
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
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

func cloneRepos(repos []string, dirName, token string) {
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
				Password: token,
			},
			URL:      repo,
			Progress: os.Stdout,
		})

		if err != nil {
			log.Fatal(err)
		}
	}
}
