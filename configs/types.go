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
