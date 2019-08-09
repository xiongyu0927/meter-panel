package configs

// K8sconfigs is used get k8s and xuliehua json info from furion
type K8sconfigs struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Attr attr   `json:"attr"`
}

type attr struct {
	Kubernetes kubernetes `json:"kubernetes"`
}

type kubernetes struct {
	Token    string `json:"token"`
	Endpoint string `json:"endpoint"`
}

// HumanSingleK8sConfigs protect your eyes
type HumanSingleK8sConfigs struct {
	EndPoint string
	Token    string
}

// HumanAllK8SConfigs save all cluster's k8s configs
type HumanAllK8SConfigs map[string]HumanSingleK8sConfigs
