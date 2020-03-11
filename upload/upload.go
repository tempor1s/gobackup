package upload

import (
	"fmt"
)

// Start will start the upload process of all repos in the current directory
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
}

// gitHub will allow you to upload all github repos in the current directory into the github repo that is associated with your personal access token
func gitHub(token, username string) {
	// TODO
}

// gitLab will allow you to upload all gitlab repos in the current directory into the gitlab repo that is associated with your personal access token
func gitLab(token, username string) {
	// TODO
}
