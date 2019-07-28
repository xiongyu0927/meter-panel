package k8s

// SingleClusterNodeList is used for parse json from single k8s cluster api
type SingleClusterNodeList struct {
	Kind     string
	Items    []node
	Metadata continuee
}

type continuee struct {
	Continue string
}

// NodeEvents is used for parse json from single k8s cluster api
type NodeEvents struct {
	Type   string
	Object node
}

type node struct {
	Metadata metadata
	Status   status
}

type metadata struct {
	Name string
}

type status struct {
	Conditions []condition
}

type condition struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

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

type HumanAllClusterNodeList map[string]HumanSingleClusterNodeList
