package upload

import "fmt"

// Start will start the upload process of all repos in the current directory
func Start(token string, args []string) {
	fmt.Println("Placeholder")
	fmt.Println(token)
}

// gitHub will allow you to upload all github repos in the current directory into the github repo that is associated with your personal access token
func gitHub(token, username string) {
	// TODO
}

// gitLab will allow you to upload all gitlab repos in the current directory into the gitlab repo that is associated with your personal access token
func gitLab(token, username string) {
	// TODO
}
