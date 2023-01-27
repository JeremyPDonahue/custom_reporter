package apiaudit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type ResourceJSON struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Annotations struct {
			KonghqComStripPath                          string `json:"konghq.com/strip-path"`
			KubectlKubernetesIoLastAppliedConfiguration string `json:"kubectl.kubernetes.io/last-applied-configuration"`
			KubernetesIoIngressClass                    string `json:"kubernetes.io/ingress.class"`
		} `json:"annotations"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Generation        int       `json:"generation"`
		Labels            struct {
			App        string `json:"app"`
			AppID      string `json:"app-id"`
			AppVersion string `json:"app-version"`
			ProjectID  string `json:"project-id"`
			SpRelease  string `json:"sp-release"`
		} `json:"labels"`
		Name            string `json:"name"`
		Namespace       string `json:"namespace"`
		ResourceVersion string `json:"resourceVersion"`
		UID             string `json:"uid"`
	} `json:"metadata"`
	Spec struct {
		Rules []struct {
			Host string `json:"host"`
			HTTP struct {
				Paths []struct {
					Backend struct {
						Service struct {
							Name string `json:"name"`
							Port struct {
								Number int `json:"number"`
							} `json:"port"`
						} `json:"service"`
					} `json:"backend"`
					Path     string `json:"path"`
					PathType string `json:"pathType"`
				} `json:"paths"`
			} `json:"http"`
		} `json:"rules"`
		TLS []struct {
			Hosts      []string `json:"hosts"`
			SecretName string   `json:"secretName"`
		} `json:"tls"`
	} `json:"spec"`
	Status struct {
		LoadBalancer struct {
			Ingress []struct {
				Hostname string `json:"hostname"`
			} `json:"ingress"`
		} `json:"loadBalancer"`
	} `json:"status"`
}

func GetNamespaces(paths []string) [][]string {
	var namespaces [][]string
	for _, path := range paths {

		output, err := exec.Command("sh", "-c", fmt.Sprintf("kubectl --kubeconfig %s get namespaces --output=name", path)).Output()
		if err != nil {
			fmt.Println("Could not get namespaces")
			os.Exit(1)
		}
		namespaces = append(namespaces, strings.Split(string(output), "\n"))
	}
	return namespaces
}

//  collectReourceJSON function is used to gather information from our kubernetes clusters in json format.
//	It takes the path of the kubeconfig file and the resource to be collected as input.
//	The function first creates a directory named "jsonFiles" if it does not already exist.
//	Then it uses the kubectl command to collect information about the specified resource in json format.
//	The output of the command is saved to a file named "<cluster>-<resource>.json" in the jsonFiles directory.
//	Finally, the function prints a message indicating that the data has been received for the specified cluster and resource.
func CollectReourceJSON(paths []string, resource string) []ResourceJSON {
	var resourceStructs []ResourceJSON
	for _, path := range paths {

		cmd := exec.Command("sh", "-c", fmt.Sprintf("kubectl --kubeconfig %s get %s -A -o json | jq -c '(.items[])'", path, resource))

		var jr ResourceJSON

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

type APIAudit struct {
	Cluster      string
	ResourceName string
	Namespace    string
	Type         string
	AppID        string
	ProjectID    string
	ApiVersion   string
}

func FilterForAPI(resourceJSON []ResourceJSON, apiVersion string) []APIAudit {

	c := regexp.MustCompile("nsk-(.*)-(prod|nonprod)")
	a := regexp.MustCompile(apiVersion)

	var filteredJSONResponses []APIAudit

	for _, resource := range resourceJSON {

		host := resource.Spec.Rules[0].Host
		cluster := c.FindString(host)
		lastAppliedConfiguration := resource.Metadata.Annotations.KubectlKubernetesIoLastAppliedConfiguration

		if api := a.FindString(lastAppliedConfiguration); api != "" {
			cmd := exec.Command("bash", "-c", fmt.Sprintf("echo '%s' | jq '([\"%s\", .metadata.name, .metadata.namespace, .kind, .metadata.labels.\"app-id\", .metadata.labels.\"project-id\", .apiVersion]) | @csv'", lastAppliedConfiguration, cluster))

			output, err := cmd.Output()
			if err != nil {
				fmt.Println(err)
			}
			data := strings.Split(string(output), ",")
			var api APIAudit

			api.Cluster = data[0]
			api.ResourceName = data[1]
			api.Namespace = data[2]
			api.Type = data[3]
			api.AppID = data[4]
			api.ProjectID = data[5]
			api.ApiVersion = data[6]

			filteredJSONResponses = append(filteredJSONResponses, api)

		}
	}

	return filteredJSONResponses
}
