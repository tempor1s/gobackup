package backup

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// GitHub will clone all users github repos to your local machine. Puts them in a <github_username>/ folder.
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

	// Start timer, create wait group so we dont exit early, and create channel for URL's
	start := time.Now()
	repos := make(chan string)
	var wg sync.WaitGroup
	// Get all repos for the user
	go getRepos(token, repos, &wg)

	// Clone all repos
	cloneRepos(repos, dirName, token, &wg)
	// Wait until all repos have been cloned before printing time and exiting
	wg.Wait()
	fmt.Println(time.Since(start))
}

// getRepos will get all the repos for a user
func getRepos(token string, c chan string, wg *sync.WaitGroup) {
	// Set up OAuth token stuff
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	// Create a new github client using the OAuth2 token
	client := github.NewClient(tc)

	// Options for our request to GitHub. This will have the API only return repos that we own, and with no pagination
	opt := &github.RepositoryListOptions{
		Affiliation: "owner",
		ListOptions: github.ListOptions{
			PerPage: 100000,
		},
	}

	// Get all repos that the user owns
	repos, _, err := client.Repositories.List(ctx, "", opt)

	if err != nil {
		log.Fatal(err)
	}

	// Add all repos to the channel
	for _, repo := range repos {
		c <- *repo.HTMLURL
		wg.Add(1)
	}

	// Close the channel, our other function will still be able to receive what is already inside of it, and this prevents a deadlock
	close(c)
}

// cloneRepos will clone all the repos in a given string slice of repo URLS
func cloneRepos(repos chan string, dirName, token string, wg *sync.WaitGroup) {
	// Create username dir to put all cloned repos in.
	err := os.MkdirAll(dirName, os.ModePerm)

	if err != nil {
		log.Fatal(err)
	}

	// Clone each repo in the channel
	for repo := range repos {
		go cloneWorker(repo, dirName, token, wg)
	}
}

// cloneWorker will clone the given repository
func cloneWorker(repo, dirName, token string, wg *sync.WaitGroup) {
	fmt.Printf("[gobackup] cloning %s\n", repo)

	// Decrement the waitgroup count when we are finished cloning the repository
	defer wg.Done()
	// Get the name of the repo we are cloning
	repoName := path.Base(repo)
	// Dirname which will be <github_username>/<repo_name>
	dirName = dirName + "/" + repoName
	// Clone the repository
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
