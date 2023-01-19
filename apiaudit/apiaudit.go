package apiaudit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func RefreshFiles(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	} else {
		d, err := os.Open(dir)
		if err != nil {
			fmt.Println("Could not open jsonFiles directory")
			fmt.Println(err)
		}
		defer d.Close()

		names, err := d.Readdirnames(-1)
		if err != nil {
			fmt.Println("Could not enumerate the contents of jsonFiles")
			fmt.Println(err)
		}

		for _, name := range names {
			err = os.Remove(filepath.Join(dir, name))
			if err != nil {
				fmt.Printf("Could not remove %s", name)
				fmt.Println(err)
			}
		}
	}
}

//  collectReourceJSON function is used to gather information from our kubernetes clusters in json format.
//	It takes the path of the kubeconfig file and the resource to be collected as input.
//	The function first creates a directory named "jsonFiles" if it does not already exist.
//	Then it uses the kubectl command to collect information about the specified resource in json format.
//	The output of the command is saved to a file named "<cluster>-<resource>.json" in the jsonFiles directory.
//	Finally, the function prints a message indicating that the data has been received for the specified cluster and resource.
func collectReourceJSON(path, resource string) {
	jsonFiles := "jsonFiles"

	cmd := exec.Command("sh", "-c", fmt.Sprintf("kubectl --kubeconfig %s get %s -A -o json | jq -c", path, resource))
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	directories := strings.Split(path, "/")
	cluster := directories[len(directories)-1]

	file, err := os.Create(jsonFiles + "/" + cluster + "-" + resource + ".json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	_, err = file.Write(stdout)
	if err != nil {
		fmt.Println(err)
	}

	current := fmt.Sprintf("data received for: %s for %s", cluster, resource)
	fmt.Println(current)
}

// IterateOverPaths function takes in a list of paths and a resource as input and performs the following operations:
// 1. Iterates over each path in the 'paths' slice
// 2. For each path, calls the 'collectReourceJSON' function passing the current path and the resource
// 3. For each path in the 'paths' slice, it collects the data for the specified resource
func IterateOverPaths(paths []string, resource string) {

	for _, path := range paths {
		collectReourceJSON(path, resource)
	}

}

// GetFileList function is used to get a list of files in a specified directory.
// It takes in a 'directory' string as input and performs the following operations:
// 1. Runs a command "ls directory/ | grep nsk" using the exec package
// 2. Captures the output of the command execution
// 3. Using the strings package, splits the output into a list of individual file names
// 4. Returns the list of file names as a slice of strings
func GetFileList(directory string, toGrepFor string) []string {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("ls %s/ | grep %s", directory, toGrepFor))
	output, _ := cmd.Output()
	files := strings.Fields(string(output))
	return files
}

// WriteToFile function is used to write to a file with a specific name.
// It takes in a 'fileName' and a 'header' string as input and performs the following operations:
// 1. Get the current date
// 2. Append the date to the fileName
// 3. Open or create a new file with the fileName
// 4. If there is an error opening the file, it will return an error message
// 5. Append the header + new line to the file
// 6. If there is an error writing to the file, it will return an error message
// 7. Close the file
// 8. Returns the fileName
// This function is useful for creating a file with a specific name and adding a header to it.
func WriteToFile(fileName string, header string) string {
	now := time.Now()
	year, month, day := now.Date()
	fileName = fmt.Sprintf("%s_%d_%d_%d.csv", fileName, year, month, day)

	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		os.Remove(fileName)
	}

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Sprintf("Unable to open %s", fileName)
		fmt.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(header + "\n"); err != nil {
		fmt.Sprintf("Unable to write to %s", fileName)
		fmt.Println(err)
	}

	return fileName
}

// CreateAPIAuditCSV function is used to create a csv file containing information regarding those with v1beta1.
// TODO: Abstract this away from being specific to v1beta1 so we can look for anything.
// It performs the following operations:
// 1. Calls the GetFileList function to get a list of files in the "jsonFiles" directory with the substring "nsk"
// 2. Creates a header for the csv file with specific columns
// 3. Calls the WriteToFile function to create and/or open a file named "apiAudit" and append the header to it.
// 4. Iterates over the list of files returned from step 1
// 5. For each file, it runs a jq command to extract specific information from the file
// 6. Appends the extracted information to the "apiAudit" file in csv format
func CreateAPIAuditCSV() {
	filesToParse := GetFileList("jsonFiles/", "nsk")
	header := "\"Cluster\",\"Name\",\"Namespace\",\"Kind\",\"App-ID\",\"Project-ID\",\"API Version\""
	fileName := WriteToFile("apiAudit", header)

	for _, file := range filesToParse {

		withoutSuffix := strings.TrimSuffix(file, ".json")
		lastIndex := strings.LastIndex(withoutSuffix, "-")
		cluster := file[:lastIndex]

		jqCommand := fmt.Sprintf("jq -r '(.items[] | [.metadata.annotations.\"kubectl.kubernetes.io/last-applied-configuration\" ]) | select(try(.[] | test(\"v1beta1\"))) | .[] | fromjson | [\"%s\", .metadata.name, .metadata.namespace, .kind, .metadata.labels.\"app-id\", .metadata.labels.\"project-id\", .apiVersion] | @csv' jsonFiles/%s >> %s", cluster, file, fileName)
		cmd := exec.Command("bash", "-c", jqCommand)
		cmd.Run()

	}
}
