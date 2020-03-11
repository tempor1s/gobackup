package backup

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/go-github/github"
	"github.com/xanzy/go-gitlab"

	"github.com/schollz/progressbar/v2"
	"golang.org/x/oauth2"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// getGithubRepos will get all a users github repos
func getGithubRepos(token, username string, repoChan chan string, totalRepos chan int, wg *sync.WaitGroup) {
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
		repos, _, err = client.Repositories.List(ctx, username, opt)
	}
	checkIfError(err)

	// Send repo count to chanel to use for ProgressBar
	totalRepos <- len(repos)

	// Add all repos to the channel and increment WaitGroup
	for _, repo := range repos {
		repoChan <- *repo.HTMLURL
		wg.Add(1)
	}

	// Close repo list channel, our other function will still be able to receive what is already inside of it, and this prevents a deadlock or infinite loop
	close(repoChan)
}

func getGitlabRepos(token, username string, repoChan chan string, totalRepos chan int, wg *sync.WaitGroup) {
	// Create a new Gitlab client with our token to make requests
	client := gitlab.NewClient(nil, token)

	// Get all the repos for a user
	repos, _, err := client.Projects.ListUserProjects(username, nil)
	checkIfError(err)

	// Send repo count to chanel to use for ProgressBar
	totalRepos <- len(repos)

	// Add all repo urls to the channel and increment WaitGroup
	for _, repo := range repos {
		repoChan <- repo.HTTPURLToRepo
		wg.Add(1)
	}

	// Close repo list channel, our other functions will still be able to receive what is already inside of it, and this prevents a deadlock or infite loop
	close(repoChan)
}

// cloneRepos will clone all the repos in a given string slice of repo URLS
func cloneRepos(repos chan string, bar *progressbar.ProgressBar, dirName, token string, wg *sync.WaitGroup) {
	// Create username dir to put all cloned repos in.
	err := os.MkdirAll(dirName, os.ModePerm)
	checkIfError(err)

	// Clone each repo in the channel
	for repo := range repos {
		go cloneWorker(repo, dirName, token, wg, bar)
	}
}

// cloneWorker will clone the given repository
func cloneWorker(repo, dirName, token string, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
	// fmt.Printf("[gobackup] cloning %s\n", repo)

	// Decrement the waitgroup count when we are finished cloning the repository and increment our progress bar
	defer bar.Add(1)
	defer wg.Done()
	// Get the name of the repo we are cloning
	repoName := path.Base(repo)

	repoName = strings.TrimSuffix(repoName, filepath.Ext(repoName))
	// Dirname which will be <github_username>/<repo_name>
	dirName = dirName + "/" + repoName

	// Setup auth if we have a token
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

	checkIfError(err)
}
