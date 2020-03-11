package upload

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/xanzy/go-gitlab"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// TODO: Change username to just be directory - internally we can still just use the username for backup command and such
// Start will start the upload process of all repos in the "username" directory
func Start(token string, args []string) {
	if len(args) == 0 {
		fmt.Println("Please pass the provider that you want to upload to. (github/gitlab) Example: `upload github`")
		return
	} else if len(args) == 1 {
		fmt.Println("Please pass in the username that you would like to upload to. Example: `upload github tempor1s`")
		return
	}

	if token == "" {
		fmt.Println("For uploading to Github you need to provide a personal access token using --token. Example: `upload github tempor1s --token=123asd`")
		return
	}

	switch args[0] {
	case "github":
		gitHub(token, args[1])
	case "gitlab":
		gitLab(token, args[1])
	}
}

// gitHub will allow you to upload all github repos in the "username" directory into the github repo that is associated with your personal access token
func gitHub(token, username string) {
	// TODO
}

// gitLab will allow you to upload all gitlab repos in the "username" directory into the gitlab repo that is associated with your personal access token
func gitLab(token, username string) {
	// Create a new gitlab client that will be our hook into the GitLab api
	client := gitlab.NewClient(nil, token)
	// Get all the directories within the username directory
	directories, err := ioutil.ReadDir(username)

	repoNames := []string{}

	if err != nil {
		log.Fatal(err)
	}

	// Double for loop here right now is gross, but will give us less work in the future when we make this concurrent
	for _, directory := range directories {
		repoNames = append(repoNames, directory.Name())
	}

	// Loop through all repos in the username directory and upload them all to GitLab as new projects
	for _, repoName := range repoNames {
		path := username + "/" + repoName

		project := createProject(client, repoName)

		if project != nil {
			createRemoteAndPush(path, token, project)
		}
	}

}

func createProject(client *gitlab.Client, name string) *gitlab.Project {
	opt := &gitlab.CreateProjectOptions{
		Name:                 gitlab.String(name),
		Description:          gitlab.String("Placeholder"),
		MergeRequestsEnabled: gitlab.Bool(true),
		SnippetsEnabled:      gitlab.Bool(true),
		Visibility:           gitlab.Visibility(gitlab.PublicVisibility),
	}

	project, _, err := client.Projects.CreateProject(opt)

	if err != nil {
		log.Println(err)
		return nil
	}

	return project
}

func createRemoteAndPush(path, token string, project *gitlab.Project) {
	r, err := git.PlainOpen(path)

	if err != nil {
		log.Fatal(err)
	}

	// Create a new remote to push to so that we maintain the old URL
	r.CreateRemote(&config.RemoteConfig{
		Name: "backup",
		URLs: []string{project.HTTPURLToRepo},
	})

	// Create auth with other token
	auth := &http.BasicAuth{
		Username: "gobackup",
		Password: token,
	}

	p := &git.PushOptions{
		RemoteName: "backup",
		Auth:       auth,
	}

	err = r.Push(p)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("New Project was created and pushed with the name %s\n", project.Name)
}
