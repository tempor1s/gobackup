package upload

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/google/go-github/github"
	"github.com/schollz/progressbar/v2"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// gitHub will allow you to upload all repos in the given directory into the github repo that is associated with your personal access token
func gitHub(token, directory string) {
	// Set up OAuth token stuff
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// Create a channel to keep all our repos
	repoChan := make(chan string)
	repoCountChan := make(chan int)
	var wg sync.WaitGroup

	// Get all the directory names in the given directory
	go getDirNames(directory, repoChan, repoCountChan, &wg)

	repoCount := <-repoCountChan

	if repoCount == 0 {
		log.Fatal("Error; no repos found to upload")
	}

	// Build basic progress bar with the amount of repos that we have
	bar := progressbar.NewOptions(repoCount, progressbar.OptionSetRenderBlankState(true))

	// Loop through all repos in the directory directory and upload them all to GitLab as new projects
	for repoName := range repoChan {
		go uploadGithubRepos(directory, repoName, token, client, &wg, bar, ctx)
	}

	wg.Wait()

	fmt.Printf("\nUpload Complete - Uploaded %d repos to GitHub.\n", repoCount)
}

func uploadGithubRepos(directory, repoName, token string, client *github.Client, wg *sync.WaitGroup, bar *progressbar.ProgressBar, ctx context.Context) {
	defer wg.Done()
	defer bar.Add(1)

	// The path to the repo, exa: tempor1s/gobackup
	path := directory + "/" + repoName

	repo := createRepo(ctx, client, repoName)

	if repo != "" {
		createGithubRemotePush(path, token, repo)
	}
}

func createRepo(ctx context.Context, client *github.Client, repoName string) string {
	isPrivate := false
	repoDescription := ""

	user, _, err := client.Users.Get(ctx, "")

	r := &github.Repository{Name: &repoName, Private: &isPrivate, Description: &repoDescription}

	// Delete old repo and then create a new one
	client.Repositories.Delete(ctx, user.GetLogin(), repoName)
	repo, _, _ := client.Repositories.Create(ctx, "", r)

	if repo == nil {
		return repo.GetURL()
	}

	if err != nil {
		log.Fatal(err)
	}

	return user.GetHTMLURL() + "/" + repoName
}

func createGithubRemotePush(path, token, repoURL string) {
	// Open the github repo at our current path
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal(err)
	}

	// Check to see if our backup remote exists
	remoteExists, _ := r.Remote("backup")

	// Create a new remote to push to so that we maintain the old URL
	if remoteExists == nil {
		r.CreateRemote(&config.RemoteConfig{
			Name: "backup",
			URLs: []string{repoURL + ".git"},
		})
	} else {
		r.DeleteRemote("backup")

		r.CreateRemote(&config.RemoteConfig{
			Name: "backup",
			URLs: []string{repoURL + ".git"},
		})
	}

	// Create auth with other token
	auth := &http.BasicAuth{
		Username: "gobackup",
		Password: token,
	}

	// Push to the remote we just created and use the auth we created above
	p := &git.PushOptions{
		RemoteName: "backup",
		Auth:       auth,
	}

	// Push all the code in our repo to the remote that we just created
	err = r.Push(p)

	if err != nil {
		log.Fatal(err)
	}
}
