package upload

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/xanzy/go-gitlab"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// Start will start the upload process of all repos in the "directory" directory
func Start(token string, args []string) {
	// Need service provider & directory
	if len(args) == 0 {
		fmt.Println("Please pass the provider that you want to upload to. (github/gitlab) Example: `upload github`")
		return
	} else if len(args) == 1 {
		fmt.Println("Please pass in the directory that you would like to upload to. Example: `upload github directory`")
		return
	}

	// We need a token for uploading to a service
	if token == "" {
		fmt.Println("For uploading to Github you need to provide a personal access token using --token. Example: `upload github directory --token=123asd`")
		return
	}

	// Do different things based off of the platform they want to use
	switch args[0] {
	case "github":
		gitHub(token, args[1])
	case "gitlab":
		gitLab(token, args[1])
	}
}

// gitHub will allow you to upload all repos in the given directory into the github repo that is associated with your personal access token
func gitHub(token, directory string) {
	// TODO
}

// gitLab will allow you to upload all repos in the given directory into the gitlab repo that is associated with your personal access token
func gitLab(token, directory string) {
	// Create a new gitlab client that will be our hook into the GitLab api
	client := gitlab.NewClient(nil, token)
	var wg sync.WaitGroup

	// Create a channel to keep all our repos
	repoChan := make(chan string)

	go getDirNames(directory, repoChan, &wg)

	// Loop through all repos in the directory directory and upload them all to GitLab as new projects
	for repoName := range repoChan {
		go uploadRepos(directory, repoName, token, client, &wg)
	}

	wg.Wait()

	fmt.Println("Upload Complete")
}

// uploadRepos is designed to be a concurrent worker that will upload the current repo
func uploadRepos(directory, repoName, token string, client *gitlab.Client, wg *sync.WaitGroup) {
	// Decrease the waitgroup after we are done uploading the current repo because we are done with all work
	defer wg.Done()
	// The path to the repo, exa: tempor1s/gobackup
	path := directory + "/" + repoName

	// Create a new project with the name of the current directory (the repo)
	project := createProject(client, repoName)

	// If the project already exists, we dont wanna do anything to it. Otherwise, create remote and push to the new project
	if project != nil {
		createRemoteAndPush(path, token, project)
	}
}

// getDirNames will get all directories in a given repo and send them to a given channel
func getDirNames(dir string, repoNames chan string, wg *sync.WaitGroup) {
	// Get all the directories within the directory directory
	directories, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	// Loop through all the directories thqat we read
	for _, directory := range directories {
		// Add the directory to our channel and increase the WaitGroup
		if directory.IsDir() {
			repoNames <- directory.Name()
			wg.Add(1)
		}
	}

	// Close the channel after we add all the names, other functions will still be able to access it :)
	close(repoNames)
}

// createProject will create a new gitlab project with a given name, and then return it
func createProject(client *gitlab.Client, name string) *gitlab.Project {
	// TODO: Do something different with description
	// TODO: Respect visability of the repo we cloned
	// Options for our new project (repo)
	opt := &gitlab.CreateProjectOptions{
		Name:                 gitlab.String(name),
		Description:          gitlab.String("Placeholder"),
		MergeRequestsEnabled: gitlab.Bool(true),
		SnippetsEnabled:      gitlab.Bool(true),
		Visibility:           gitlab.Visibility(gitlab.PublicVisibility),
	}

	// Create the new project
	project, _, err := client.Projects.CreateProject(opt)

	// If the repo already exists, just tell the user and move on
	if err != nil {
		log.Println(err)
		return nil
	}

	return project
}

// createRemoteAndPush will create a new remote to the backup repository and then push the code to that remote (gitlab repo we create above)
func createRemoteAndPush(path, token string, project *gitlab.Project) {
	// Open the github repo at our current path
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal(err)
	}

	// Check to see if our backup remote exists
	remoteExists, err := r.Remote("backup")

	if err != nil {
		fmt.Println("Remote does not exist... Creating..")
	}

	// Only create a new remote if one does not already exist - this will future proof for doing backups to existing repos
	if remoteExists == nil {
		// Create a new remote to push to so that we maintain the old URL
		r.CreateRemote(&config.RemoteConfig{
			Name: "backup",
			URLs: []string{project.HTTPURLToRepo},
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

	// TODO: We can remove this when we create a progress bar
	fmt.Printf("New Project was created and pushed with the name %s\n", project.Name)
}
