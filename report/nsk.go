package report

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// SearchKubeDir searches the $HOME/.kube directory and it's subdirectories looking for JSON files with nsk in the filename
// and returns a slice of strings containing the filenames without the file extention
func SearchKubeDir() ([]string, error) {
	var filePaths []string

	// We compile a regex to look for any filenames for either nsk-*-(prod|nonprd)
	// This way we only focus on clusters following our normal naming convention
	// avoiding clusters named things like nsk-main etc.
	re := regexp.MustCompile(`^nsk-.+(prod|nonprod).yaml$`)

	// The os.UserHomeDir is called to get the users home directory (works for MacOS and Windows)
	// This directory is used as the base for the search.
	homeDir, err := os.UserHomeDir()
	// If there is an error getting the home directory, it is printed and nil is returned
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// The filepath.Join() function is called to join the home directory and the .kube directory into a single file path
	// This is the directory that will be searched
	kubeDir := filepath.Join(homeDir, ".kube")

	// Walk the directory tree using filepath.Walk, passing each item found to a call back funtion defined within Walk()
	// This callback funtion takes three parameters: the path to the file in the form of a string -
	// the file's info which gets extracted using os.FileInfo and stored in a struct called info -
	// and an error called err - this function also returns an error value as indicated by the error keywork preceding the opening curly brace.
	err = filepath.Walk(kubeDir, func(path string, info os.FileInfo, err error) error {

		// If at any point there is an issue reading an entry within the directory we throw an error
		// This err refers to the 'err' argument in the Walk() function.
		if err != nil {
			return err
		}

		// As we search through the .kube directory we look for any file that has nsk-*-(prod|nonprodls) in the title AND DOES NOT have a yaml file extention
		if re.FindString(info.Name()) != "" {
			// Once it finds something that matches our criteria, it appends it to our String slice
			filePaths = append(filePaths, path)
		}
		// This funtion doesn't need to return any of the values collected so it just returns nil to meet the requirement set by the return value indicator (error)
		return nil
	})

	// Here we check the value of this err variable: 'err = filepath.Walk'
	if err != nil {
		return nil, err
	}
	return filePaths, nil
}

// CreateKubectlCommands assembles the kubectl commands that will be used to query the NSK clusters
func CreateKubectlCommands(filePaths []string) []string {

	fmt.Println("Creating kubectl commands...")

	var cmds []string

	for _, filePath := range filePaths {

		if filePath != "" {

			c := fmt.Sprintf("kubectl --kubeconfig %s get ns", filePath)

			cmds = append(cmds, c)
		}
	}

	return cmds
}

// GoGrabAppIds queries the clusters for a list of appids/namespaces
func GoGrabAppIds(cmds []string) [][]string {
	var c KubectlCommand

	var clusterNamespaces [][]string

	fmt.Println("Querying clusters for appid's...")

	for _, command := range cmds {

		s := strings.Split(command, " ")

		c.kubectl = s[0]          // kubectl
		c.kubeconfigFlag = s[1]   // --kubeconfig
		c.pathToKubeConfig = s[2] // path to kubeconfig e.g. /Users/gx8k/.kube/nsk-beet-prod
		c.get = s[3]              // get
		c.ns = s[4]               // ns

		clusterName := isItACluster.FindString(c.pathToKubeConfig)

		fmt.Println(strings.Trim(clusterName, ".yaml"))

		cmd := exec.Command(c.kubectl, c.kubeconfigFlag, c.pathToKubeConfig, c.get, c.ns)
		b, err := cmd.CombinedOutput()
		if err != nil {
			log.Println("Failed to assemble kubectl command")
			log.Println("Have you run aws-okta login?")
			log.Println(string(b))
			log.Fatal(err)
		}

		output := strings.Split(string(b), "\n")

		var sliceOfAppIds []string
		sliceOfAppIds = append(sliceOfAppIds, clusterName)

		for _, appid := range output {

			if a := isItANamespace.FindString(appid); a != "" {
				sliceOfAppIds = append(sliceOfAppIds, a)
			}
		}

		clusterNamespaces = append(clusterNamespaces, sliceOfAppIds)

	}

	return clusterNamespaces

}
