package apiaudit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func CollectAndFilterResourceJSON(paths []string, resource, apiVersion string) []APIAudit {
	var responses []APIAudit
	switch resource {
	case "ingress":
		responses = FilterIngressForAPI(CollectIngressJSON(paths, resource), apiVersion)

	case "role":
		responses = FilterRoleForAPI(CollectRoleJSON(paths, resource), apiVersion)

	case "rolebinding":
		responses = FilterRoleBindingForAPI(CollectRoleBindingJSON(paths, resource), apiVersion)

	}

	return responses
}

//	 collectReourceJSON function is used to gather information from our kubernetes clusters in json format.
//		It takes the path of the kubeconfig file and the resource to be collected as input.
//		The function first creates a directory named "jsonFiles" if it does not already exist.
//		Then it uses the kubectl command to collect information about the specified resource in json format.
//		The output of the command is saved to a file named "<cluster>-<resource>.json" in the jsonFiles directory.
//		Finally, the function prints a message indicating that the data has been received for the specified cluster and resource.
func CollectIngressJSON(paths []string, resource string) []IngressJSON {
	var resourceStructs []IngressJSON
	fmt.Print(".")

	for _, path := range paths {
		fmt.Print(".")

		cmd := exec.Command("sh", "-c", fmt.Sprintf("kubectl --kubeconfig %s get %s -A -o json | jq -c '(.items[])'", path, resource))

		var jr IngressJSON

		// fmt.Printf("Collecting resource JSON for %s \n", path)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}

		decoder := json.NewDecoder(bytes.NewReader(output))
		for {

			if err := decoder.Decode(&jr); err != nil {
				break
			}
			resourceStructs = append(resourceStructs, jr)
		}
	}

	return resourceStructs
}

func CollectRoleJSON(paths []string, resource string) []RoleJSON {
	var resourceStructs []RoleJSON
	fmt.Print(".")

	for _, path := range paths {
		fmt.Print(".")

		cmd := exec.Command("sh", "-c", fmt.Sprintf("kubectl --kubeconfig %s get %s -A -o json | jq -c '(.items[])'", path, resource))

		var jr RoleJSON

		// fmt.Printf("Collecting resource JSON for %s \n", path)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}

		decoder := json.NewDecoder(bytes.NewReader(output))
		for {

			if err := decoder.Decode(&jr); err != nil {
				break
			}
			jr.Cluster = filepath.Base(path)
			resourceStructs = append(resourceStructs, jr)
		}
	}

	return resourceStructs
}

func CollectRoleBindingJSON(paths []string, resource string) []RoleBindingJSON {
	var resourceStructs []RoleBindingJSON
	fmt.Print(".")

	for _, path := range paths {
		fmt.Print(".")

		cmd := exec.Command("sh", "-c", fmt.Sprintf("kubectl --kubeconfig %s get %s -A -o json | jq -c '(.items[])'", path, resource))

		var jr RoleBindingJSON

		// fmt.Printf("Collecting RoleBinding JSON for %s \n", path)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}

		decoder := json.NewDecoder(bytes.NewReader(output))
		for {

			if err := decoder.Decode(&jr); err != nil {
				break
			}
			jr.Cluster = filepath.Base(path)
			resourceStructs = append(resourceStructs, jr)
		}
	}

	return resourceStructs
}

func FilterIngressForAPI(resourceJSON []IngressJSON, apiVersion string) []APIAudit {

	c := regexp.MustCompile("nsk-(.*)-(prod|nonprod)")
	a := regexp.MustCompile(apiVersion)

	var filteredJSONResponses []APIAudit
	var cluster string
	fmt.Println(".")

	for _, resource := range resourceJSON {

		host := resource.Spec.Rules[0].Host
		sp := resource.Metadata.Labels.SpRelease
		cluster = c.FindString(host)
		lastAppliedConfiguration := resource.Metadata.Annotations.KubectlKubernetesIoLastAppliedConfiguration

		if api := a.FindString(lastAppliedConfiguration); api != "" {
			cmd := exec.Command("bash", "-c", fmt.Sprintf("echo '%s' | jq '([.metadata.name, .metadata.namespace, .kind, .metadata.labels.\"app-id\", .metadata.labels.\"project-id\", .apiVersion]) | @csv'", lastAppliedConfiguration))

			output, err := cmd.Output()
			if err != nil {
				fmt.Println(err)
			}
			data := strings.Split(string(output), ",")
			var api APIAudit

			api.Cluster = cluster
			api.ResourceName = strings.Trim(data[0], "\"\\")
			api.Namespace = strings.Trim(data[1], "\"\\")
			api.Type = strings.Trim(data[2], "\"\\")
			api.AppID = strings.Trim(data[3], "\"\\")
			api.ProjectID = strings.Trim(data[4], "\"\\")
			api.ApiVersion = strings.Trim(data[5], "\"\\")
			api.SpRelease = sp
			str, proj := GetGitLabJSON(api.ProjectID)
			b, err := GetArchivedStatus(str, proj)
			if err != nil {
				fmt.Println("Trouble getting archive status")
				os.Exit(1)
			}
			api.ArchiveStatus = fmt.Sprintf("%v", b)
			if api.ProjectID == "" {
				api.ArchiveStatus = "Unknown"
			}

			filteredJSONResponses = append(filteredJSONResponses, api)

		}

	}

	return filteredJSONResponses
}

func FilterRoleForAPI(resourceJSON []RoleJSON, apiVersion string) []APIAudit {

	a := regexp.MustCompile(apiVersion)

	var filteredJSONResponses []APIAudit
	fmt.Print(".")

	for _, resource := range resourceJSON {

		sp := resource.Metadata.Labels.SpRelease
		cluster := strings.Trim(resource.Cluster, ".yaml")
		lastAppliedConfiguration := resource.Metadata.Annotations.KubectlKubernetesIoLastAppliedConfiguration

		// fmt.Printf("Filtering Role JSON for %s \n", cluster)

		if api := a.FindString(lastAppliedConfiguration); api != "" {
			cmd := exec.Command("bash", "-c", fmt.Sprintf("echo '%s' | jq '([.metadata.name, .metadata.namespace, .kind, .metadata.labels.\"app-id\", .metadata.labels.\"project-id\", .apiVersion]) | @csv'", lastAppliedConfiguration))

			output, err := cmd.Output()
			if err != nil {
				fmt.Println(err)
			}
			data := strings.Split(string(output), ",")

			var api APIAudit

			api.Cluster = cluster
			api.ResourceName = strings.Trim(data[0], "\"\\")
			api.Namespace = strings.Trim(data[1], "\"\\")
			api.Type = strings.Trim(data[2], "\"\\")
			api.AppID = strings.Trim(data[3], "\"\\")
			api.ProjectID = strings.Trim(data[4], "\"\\")
			api.ApiVersion = strings.Trim(data[5], "\"\\")
			api.SpRelease = sp
			str, proj := GetGitLabJSON(api.ProjectID)
			b, err := GetArchivedStatus(str, proj)
			if err != nil {
				fmt.Println("Trouble getting archive status")
				os.Exit(1)
			}
			api.ArchiveStatus = fmt.Sprintf("%v", b)
			if api.ProjectID == "" {
				api.ArchiveStatus = "Unknown"
			}

			filteredJSONResponses = append(filteredJSONResponses, api)

		}

	}
	return filteredJSONResponses
}

func FilterRoleBindingForAPI(resourceJSON []RoleBindingJSON, apiVersion string) []APIAudit {

	a := regexp.MustCompile(apiVersion)

	var filteredJSONResponses []APIAudit
	fmt.Print(".")

	for _, resource := range resourceJSON {

		sp := resource.Metadata.Labels.SpRelease
		cluster := strings.Trim(resource.Cluster, ".yaml")
		lastAppliedConfiguration := resource.Metadata.Annotations.KubectlKubernetesIoLastAppliedConfiguration

		// fmt.Printf("Filtering RoleBinding JSON for %s \n", cluster)

		if api := a.FindString(lastAppliedConfiguration); api != "" {
			cmd := exec.Command("bash", "-c", fmt.Sprintf("echo '%s' | jq '([.metadata.name, .metadata.namespace, .kind, .metadata.labels.\"app-id\", .metadata.labels.\"project-id\", .apiVersion]) | @csv'", lastAppliedConfiguration))

			output, err := cmd.Output()
			if err != nil {
				fmt.Println(err)
			}
			data := strings.Split(string(output), ",")
			var api APIAudit

			api.Cluster = cluster
			api.ResourceName = strings.Trim(data[0], "\"\\")
			api.Namespace = strings.Trim(data[1], "\"\\")
			api.Type = strings.Trim(data[2], "\"\\")
			api.AppID = strings.Trim(data[3], "\"\\")
			api.ProjectID = strings.Trim(data[4], "\"\\")
			api.ApiVersion = strings.Trim(data[5], "\"\\")
			api.SpRelease = sp
			str, proj := GetGitLabJSON(api.ProjectID)
			b, err := GetArchivedStatus(str, proj)
			if err != nil {
				fmt.Println("Trouble getting archive status")
				os.Exit(1)
			}
			api.ArchiveStatus = fmt.Sprintf("%v", b)
			if api.ProjectID == "" {
				api.ArchiveStatus = "Unknown"
			}

			filteredJSONResponses = append(filteredJSONResponses, api)

		}

	}
	return filteredJSONResponses
}

func GetGitLabJSON(projectID string) (string, string) {
	fmt.Print(".")
	token, _ := os.LookupEnv("GITLAB_API_TOKEN")
	api_endpoint, _ := os.LookupEnv("GITLAB_API")
	cmd := exec.Command("sh", "-c", fmt.Sprintf("curl --header \"Authorization: Bearer %s\" %s%s", token, api_endpoint, projectID))
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Manually check projectID: %s (Will be set to FLASE in report)\n", projectID)
	}
	return string(output), projectID
}

func GetArchivedStatus(jsonStr, projectID string) (bool, error) {

	if jsonStr == "" || projectID == "" {
		return false, nil
	}
	fmt.Print(".")
	if jsonStr[0] == '[' && jsonStr[len(jsonStr)-1] == ']' {
		jsonStr = jsonStr[1 : len(jsonStr)-1]
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		fmt.Printf("\n Manually check projectID: %s (Will be set to FLASE in report)\n", projectID)
		if strings.Contains(err.Error(), "invalid character ',' after top-level value") {
			return false, nil
		}
	}

	if value, ok := data["archived"]; ok {
		if b, ok := value.(bool); ok {
			return b, nil
		}
		return false, fmt.Errorf("value of archived is not a boolean.")
	}

	if message, ok := data["message"]; ok {
		if message == "404 Project Not Found" {
			return true, nil
		}
		return false, fmt.Errorf("unexpected message value: %v", message)
	}
	return false, fmt.Errorf("archived key not found")
}
