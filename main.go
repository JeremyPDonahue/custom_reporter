package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/JeremyPDonahue/custom_reporter/apiaudit"
	"github.com/JeremyPDonahue/custom_reporter/files"
	"github.com/JeremyPDonahue/custom_reporter/report"
)

// var f []string = files.GetFileList("jsonFiles", "nsk")

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

	clusterList := files.ExtractLastPaths(pathsToKubeConfigs)

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

	var i int
	for _, resource := range resources {

		filtered := apiaudit.CollectAndFilterResourceJSON(filePathsToKube, resource, "v1beta1")
		filterInterface := make([]interface{}, len(filtered))
		for i, v := range filtered {
			filterInterface[i] = v
		}
		if i == 0 {
			files.WriteToFile("apiAudit", files.HeadersFromStructFields(filterInterface[0]))
			files.StructsToAPICSV(filterInterface, "apiAudit")
			i++
		} else {
			files.StructsToAPICSV(filterInterface, "apiAudit")
		}

	}

}

func checkAWSLogin() {
	cmd := "aws-okta info | awk 'NR==3 {print $NF}'"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// check the output for "logged in"
	outputString := string(out)

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
