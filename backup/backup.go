package backup

import (
	"fmt"
	"sync"

	"github.com/schollz/progressbar/v2"
)

// Start will start the command, handle arguments, and dispatch the correct handler
func Start(token string, args []string) {
	if len(args) == 0 {
		fmt.Println("Please pass the provider that you want to clone from. (github/gitlab/bitbucket) Example: `backup github`")
		return
	} else if len(args) == 1 {
		fmt.Println("Please pass in the username that you would like to clone from. Example: `backup github tempor1s`")
		return
	}

	if token == "" {
		fmt.Println("WARNING: Personal token was not passed in, only cloning public repos. If you want to clone private repos, please supply a token using --token")
	}

	fmt.Printf("Backing up your repos... Please wait - Don't worry if the bar freezes, this could take a few minutes :)\n\n")

	// Check the service to backup and then dispatch to the correct handler with token and username
	switch args[0] {
	case "github":
		gitHub(token, args[1])
	case "gitlab":
		gitLab(token, args[1])
	case "bitbucket":
		bitBucket(token, args[1])
	}
}

// gitHub will clone all of a users github repos to your local machine. Puts them in a <github_username>/ folder.
func gitHub(token string, username string) {
	// Start timer, create wait group so we dont exit early, and create channel for URL's
	repos := make(chan string)
	repoChan := make(chan int)
	var wg sync.WaitGroup

	// Get all repos for the user
	go getGithubRepos(token, username, repos, repoChan, &wg)

	// Get length of repos for the max of our progress bar
	repoCount := <-repoChan
	bar := progressbar.NewOptions(repoCount, progressbar.OptionSetRenderBlankState(true))

	// Clone all repos
	cloneRepos(repos, bar, username, token, &wg)
	// Wait until all repos have been cloned before printing time and exiting
	wg.Wait()

	// Get the total size of all the cloned directories and print information
	size := getDirSizeStr(username)

	fmt.Printf("\n\nCloning repos complete. Cloned %d repos with a total size of %s\n", repoCount+1, size)
}

// gitLab will clone all of a users gitlab repos to your local machine. Puts them in a <gitlab_username>/ folder.
func gitLab(token string, username string) {
	// TODO
	fmt.Println("Coming soon!")
}

// bitBucket will clone all of a users bitbucket repos to your local machine. Puts them in a <bitbucket_username>/ folder.
func bitBucket(token string, username string) {
	// TODO
	fmt.Println("Coming soon!")
}
