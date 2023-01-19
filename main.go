package main

import (
	"flag"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/JeremyPDonahue/custom_reporter/apiaudit"
	"github.com/JeremyPDonahue/custom_reporter/report"
)

func main() {

	checkAWSLogin()

	// This code creates two variables flagNSR and flagAPI, both of type bool, which will be used to store the value of the command line flags "nsr" and "api" respectively.
	// The flag package is used to register the flags "nsr" and "api" with their default values set to "false"
	// The flag package is then used to parse the command line arguments and assign values to the variables flagNSR and flagAPI.
	// The code checks the value of the flags by using the Lookup function, if the value of "nsr" flag is "true" it calls the nsReport() function, if the value of "api" flag is "true" it calls the apiAudit() function.
	// If no flag is passed at the command line, it prints "no flag passed (-nsr || -api)"

	var flagNSR, flagAPI bool
	// Register the flags
	flag.BoolVar(&flagNSR, "nsr", false, "run nsreport")
	flag.BoolVar(&flagAPI, "api", false, "run api audit")
	// Parse the flags
	flag.Parse()
	// Check the value of the flags
	if flag.Lookup("nsr").Value.String() == "true" {
		nsReport()
	} else if flag.Lookup("api").Value.String() == "true" {
		resources := []string{"ingress", "role", "rolebinding"}
		apiAudit(resources)
	} else {
		fmt.Println("no flag passed (-nsr || -api)")
	}

}

// The nsReport function is used to generate the Namespace Report when the flag -nsr is passed.
// It performs the following operations:
// 1. Calls the SearchKubeDir function to get a list of paths to kube config files located in the "./kube" directory
// 2. If there is an error getting the paths, it prints the error
// 3. Calls the CreateKubectlCommands function passing the list of paths, to create a list of kubectl commands
// 4. Calls the GoGrabAppIds function passing the list of kubectl commands, to get a list of clusters and their associated app-ids
// 5. Calls the GoGrabServiceNowData function passing the list of clusters and app-ids to get additional data about the clusters from ServiceNow
// 6. Calls the CreateCsvFile function passing the data returned from the previous step and the list of paths to kube config files, to generate an excel file containing the report.
func nsReport() {

	// Grabbing the list of clusters from your ./kube directory
	pathsToKubeConfigs, err := report.SearchKubeDir()
	if err != nil {
		fmt.Println(err)
	}

	// Using that list to populate the kubectl commands
	commandList := report.CreateKubectlCommands(pathsToKubeConfigs)

	// Using those kubectl commands to grab the appids
	clusterNamespaces := report.GoGrabAppIds(commandList)

	// Using those appids to query servicenow
	customerData := report.GoGrabServiceNowData(clusterNamespaces)

	clusterList := extractLastPath(pathsToKubeConfigs)

	// Generate excel files
	report.CreateCsvFile(customerData, clusterList)
}

// This function performs an API audit on Kubernetes resources.
// It uses the report.SearchKubeDir() function to source your kubeconfig, and if there is an error it will print it out.
// The function then iterates over the kubeconfigs using the apiaudit.IterateOverPaths() function, passing in the resource type.
// Finally, it creates a CSV file with the results of the audit using the apiaudit.CreateAPIAuditCSV() function.
func apiAudit(resources []string) {
	var filePathsToKube, err = report.SearchKubeDir()
	if err != nil {
		fmt.Println(err)
	}

	apiaudit.RefreshFiles("jsonFiles")

	for _, resource := range resources {

		apiaudit.IterateOverPaths(filePathsToKube, resource)
	}

	apiaudit.CreateAPIAuditCSV()
}

// This function takes in a slice of strings representing file paths and returns a new slice of strings containing the base of each file path.
// The base of a file path is the last element of the path, which is usually the file name.
// For example, if the input is a slice of paths ["/Users/yx4h/.kube/nsk-beet-prod", "/Users/yx4h/.kube/nsk-curry-nonprod"],
// the output will be a slice containing ["nsk-beet-prod", "nsk-beet-nonprod"]
func extractLastPath(paths []string) []string {
	var lastPaths []string
	for _, path := range paths {
		lastPaths = append(lastPaths, filepath.Base(path))
	}
	return lastPaths
}

func checkAWSLogin() {
	// check if the user is logged in
	awsOkta := exec.Command("aws-okta", "info")

	awk := exec.Command("awk", "NR==3 {print $NF}")

	awk.Stdin, _ = awsOkta.StdoutPipe()

	awk.Start()

	awsOkta.Run()

	awk.Wait()

	output, _ := awk.Output()

	// check the output for "logged in"
	outputString := string(output)
	if strings.Contains(outputString, "-") {
		fmt.Println("Not logged in, running 'aws-okta login'...")
		// if not logged in, run the login command
		loginCmd := exec.Command("aws-okta", "login")
		loginOutput, loginErr := loginCmd.Output()
		if loginErr != nil {
			fmt.Println("Error running 'aws-okta login':", loginErr)
			return
		}
		fmt.Println(string(loginOutput))
	} else {
		fmt.Println("Already logged in.")
	}
}
