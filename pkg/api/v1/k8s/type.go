package k8s

// SingleClusterNodeList is used for parse json from single k8s node list  api
type SingleClusterNodeList struct {
	Kind     string
	Items    []node
	Metadata continuee
}

type SingleClusterApplicationsList struct {
	Kind     string
	Items    []app
	Metadata continuee
}

type SingleClusterPodsList struct {
	Kind     string
	Items    []pod
	Metadata continuee
}

type SingleClusterPvList struct {
	Metadata continuee
	Items    []pv
}

type pv struct {
	Metadata metadata
	Spec     spec
}

type pod struct {
	Metadata metadata
	Spec     spec
	Status   podstatus
}

type podstatus struct {
	Phase string
	PodIp string `json:"podIp"`
}

type app struct {
	Metadata metadata
	Spec     spec
}

type metadata struct {
	Name      string
	Namespace string
	Labels    labels
}

type spec struct {
	Selector selector
	Capacity capacity
}

type capacity struct {
	Storage string
}

type selector struct {
	MatchLabels labels `json:"matchLabels"`
}

type continuee struct {
	Continue string
}

type labels struct {
	Service_name string `json:"Service_name"`
	Apps         string `json:"app.alauda.io/name"`
	Appredis     string `json:"app_name"`
}

// NodeEvents is used for parse json from single k8s node watch api
type NodeEvents struct {
	Type   string
	Object node
}

type node struct {
	Metadata metadata
	Status   status
}

type status struct {
	Conditions []condition
}

type condition struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

// PodEvents is used for pare json fron single k8s pods watch api
type PodEvents struct {
	Type   string
	Object object
}

type object struct {
	Metadata metadata
	Status   podstatus
}

// human type is real type
// HumanSingleClusterNodeList is used show a beuatiful of  Single k8s Cluster Node List
type HumanSingleClusterNodeList struct {
	SingleClusterHealthyNode   healthynode
	SingleClusterUnHealthyNode unhealthynode
	SingleClusterNodeNumber    *int
}

type healthynode struct {
	NodeStatus map[string]string
	Number     *int
}

type unhealthynode struct {
	NodeStatus map[string]string
	Number     *int
}

type HumanSingleClusterApplicationsList struct {
	SingleClusterHealthyApp   healthyapp
	SingleClusterUnHealthyApp unhealthyapp
	SingleClusterAppNumber    *int
}

type healthyapp struct {
	AppStatus map[string]Pod
	Number    *int
}

type unhealthyapp struct {
	AppStatus map[string]Pod
	Number    *int
}

type HumanSingleClusterPodsList struct {
	SingleClusterHealthyPods   healthypods
	SingleClusterUnHealthyPods unhealthypods
	SingleClusterPodsNumber    *int
}

type healthypods struct {
	PodStatus map[string]Pod
	Number    *int
}

type unhealthypods struct {
	PodStatus map[string]Pod
	Number    *int
}

type Pod struct {
	Status       string
	PodIp        string
	Service_name string `json:"Service_name"`
	Apps         string `json:"app.alauda.io/name"`
	Appredis     string `json:"app_name"`
}

type HumanSingleClusterPvsList struct {
	PvStatus map[string]string
	AllStore string
}

// HumanAllClusterNodeList is used store all cluster's nodes
type HumanAllClusterNodeList map[string]HumanSingleClusterNodeList

type HumanAllClusterApplicationsList map[string]HumanSingleClusterApplicationsList

type HumanAllClusterPodsList map[string]HumanSingleClusterPodsList

type HumanAllClusterPvsList map[string]HumanSingleClusterPvsList
