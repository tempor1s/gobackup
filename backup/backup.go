package backup

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	"github.com/google/go-github/github"
	"github.com/schollz/progressbar/v2"
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
	}

	fmt.Printf("Backing up your repos... Please wait - Don't worry if the bar freezes, this could take a few minutes :)\n\n")

	// Get the URL to clone
	repoURL := args[0]

	// Get the users github name for the directory
	dirName := path.Base(repoURL)

	// Start timer, create wait group so we dont exit early, and create channel for URL's
	repos := make(chan string)
	repoCount := make(chan int)
	var wg sync.WaitGroup

	// Get all repos for the user
	go getRepos(token, dirName, repos, repoCount, &wg)

	// Get length of repos for the max of our progress bar
	bar := progressbar.NewOptions(<-repoCount, progressbar.OptionSetRenderBlankState(true))

	// Clone all repos
	cloneRepos(repos, bar, dirName, token, &wg)
	// Wait until all repos have been cloned before printing time and exiting
	wg.Wait()

	fmt.Printf("\n\nCloning repos complete. Thanks for using GoClones!\n")
}

// getRepos will get all the repos for a user
func getRepos(token, userName string, c chan string, count chan int, wg *sync.WaitGroup) {
	// Set up OAuth token stuff
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	// Create a new github client using the OAuth2 token or no token
	var client *github.Client
	if token != "" {
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	// Options for our request to GitHub. This will have the API only return repos that we own, and with no pagination
	opt := &github.RepositoryListOptions{
		Affiliation: "owner",
		ListOptions: github.ListOptions{
			PerPage: 100000,
		},
	}

	// Get all repos that the user owns
	var repos []*github.Repository
	var err error

	if token != "" {
		// Get private repos if we have a token
		repos, _, err = client.Repositories.List(ctx, "", opt)
	} else {
		// Get only public repos if we have no token, allowing us to clone other peoples repos as well
		repos, _, err = client.Repositories.List(ctx, userName, opt)
	}

	if err != nil {
		log.Fatal(err)
	}

	// Send length of bar to channel to use for ProgressBar
	count <- len(repos)

	// Add all repos to the channel
	for _, repo := range repos {
		c <- *repo.HTMLURL
		wg.Add(1)
	}

	// Close the channel, our other function will still be able to receive what is already inside of it, and this prevents a deadlock
	close(c)
}

// cloneRepos will clone all the repos in a given string slice of repo URLS
func cloneRepos(repos chan string, bar *progressbar.ProgressBar, dirName, token string, wg *sync.WaitGroup) {
	// Create username dir to put all cloned repos in.
	err := os.MkdirAll(dirName, os.ModePerm)

	if err != nil {
		log.Fatal(err)
	}

	// Clone each repo in the channel
	for repo := range repos {
		go cloneWorker(repo, dirName, token, wg, bar)
	}
}

// cloneWorker will clone the given repository
func cloneWorker(repo, dirName, token string, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
	// fmt.Printf("[gobackup] cloning %s\n", repo)

	// Decrement the waitgroup count when we are finished cloning the repository
	defer bar.Add(1)
	defer wg.Done()
	// Get the name of the repo we are cloning
	repoName := path.Base(repo)
	// Dirname which will be <github_username>/<repo_name>
	dirName = dirName + "/" + repoName

	// Setup auth
	var auth *http.BasicAuth
	if token != "" {
		// If we have a token
		auth = &http.BasicAuth{
			Username: "gobackup",
			Password: token,
		}
	} else {
		// If we have no token, we dont want to use any auth
		auth = nil
	}
	// Clone the repository
	_, err := git.PlainClone(dirName, false, &git.CloneOptions{
		Auth: auth,
		URL:  repo,
	})

	if err != nil {
		log.Fatal(err)
	}
}
