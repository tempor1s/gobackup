package upload

import (
	"io/ioutil"
	"log"
	"sync"
)

// getDirNames will get all directories in a given repo and send them to a given channel
func getDirNames(dir string, repoNames chan string, repoCountChan chan int, wg *sync.WaitGroup) {
	// Get all the directories within the directory directory
	directories, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	//TODO: Handle this length also counting files
	repoCountChan <- len(directories)

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
