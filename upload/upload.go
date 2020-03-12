package upload

import (
	"fmt"
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

	fmt.Printf("Backing up your repos... Please wait - Don't worry if the bar freezes, this could take a few minutes :)\n\n")

	// Do different things based off of the platform they want to use
	switch args[0] {
	case "github":
		gitHub(token, args[1])
	case "gitlab":
		gitLab(token, args[1])
	}
}
