package report

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// GatherLegacyNamespaceInfo looks in the ~/.kube/config file for the current context and queries the cluster for namespaces
func GatherLegacyNamespaceInfo() [][]string {

	var legNS [][]string

	kubeconfig := os.Getenv("HOME") + "/.kube/config"

	// use the current context in kubeconfig
	config, buildErr := clientcmd.BuildConfigFromFlags("", kubeconfig)

	var (
		clientset *kubernetes.Clientset
		err       error
	)

	if buildErr == nil {
		// create the clientset to be returned
		clientset, err = kubernetes.NewForConfig(config)
	} else {
		fmt.Println("error getting config:", buildErr)
		os.Exit(1)
	}

	var (
		nsList *v1.NamespaceList
		opts   metav1.ListOptions
	)

	if err == nil {

		nsList, err = clientset.CoreV1().Namespaces().List(context.TODO(), opts)
		if err != nil {
			fmt.Println("Problem getting namespaces")
			log.Fatal(err)
		}

		exclude := []string{"default", "kube-node-lease", "kube-public", "kube-system", "shared-ingress", "k8s-sandbox"}

		for i := 0; i < len(nsList.Items); i++ {
			if Contains(exclude, nsList.Items[i].ObjectMeta.Name) {
				continue
			} else {
				var pair []string
				pair = append(pair, nsList.Items[i].ObjectMeta.Name)
				pair = append(pair, nsList.Items[i].ObjectMeta.Labels["nordappid"])
				legNS = append(legNS, pair)
				i++
			}
		}

	} else {
		fmt.Println("error getting ns:", err)
		os.Exit(1)
	}

	return legNS

}

// ChangeContext sets the current context to the value passed to it. It is a necessary helper funtion for GatherLegacyNSInfoForAllClusters
func ChangeContext(cluster string) {

	cmd := exec.Command("kubectl", "config", "use-context", cluster)

	_, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Could not change context on legacy cluster %s", cluster)
		fmt.Println(err)
		os.Exit(1)
	}
}

// GatherLegacyNSInfoForAllClusters uses GatherLegacyNamespaceInfo and ChangeContext to gather namespace info for all three legacy clusters
func GatherLegacyNSInfoForAllClusters() [][]string {

	var legNS [][]string
	clusters := []string{"steel", "hydrogen", "barcelona"}

	for _, cluster := range clusters {
		ChangeContext(cluster)

		clusterNS := GatherLegacyNamespaceInfo()

		for _, ns := range clusterNS {

			if !SliceContains(legNS, ns) {
				legNS = append(legNS, ns)
			}
		}
	}

	return legNS

}
