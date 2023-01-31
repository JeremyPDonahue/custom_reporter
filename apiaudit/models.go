package apiaudit

import "time"

type APIAudit struct {
	Cluster       string
	ResourceName  string
	Namespace     string
	Type          string
	AppID         string
	ProjectID     string
	ApiVersion    string
	SpRelease     string
	ArchiveStatus string
}

type IngressJSON struct {
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

type RoleJSON struct {
	Cluster    string
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Annotations struct {
			KubectlKubernetesIoLastAppliedConfiguration string `json:"kubectl.kubernetes.io/last-applied-configuration"`
		} `json:"annotations"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
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
	Rules []struct {
		APIGroups []string `json:"apiGroups"`
		Resources []string `json:"resources"`
		Verbs     []string `json:"verbs"`
	} `json:"rules"`
}

type RoleBindingJSON struct {
	Cluster    string
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Annotations struct {
			KubectlKubernetesIoLastAppliedConfiguration string `json:"kubectl.kubernetes.io/last-applied-configuration"`
		} `json:"annotations"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
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
	RoleRef struct {
		APIGroup string `json:"apiGroup"`
		Kind     string `json:"kind"`
		Name     string `json:"name"`
	} `json:"roleRef"`
	Subjects []struct {
		Kind string `json:"kind"`
		Name string `json:"name"`
	} `json:"subjects"`
}
