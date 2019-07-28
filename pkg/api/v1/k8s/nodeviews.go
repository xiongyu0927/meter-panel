package k8s

import (
	"encoding/json"
	"meter-panel/configs"
	"meter-panel/tools"
)

var (
	nodelist SingleClusterNodeList
	// HumanAllClusterNodeList is a map, key is cluster name, value is all node of this cluster
	_HumanAllClusterNodeList HumanAllClusterNodeList = make(map[string]HumanSingleClusterNodeList)
	// NilHumanSingleClusterNodeList is used return a nil value that type of HumanSingleClusterNodeList
	NilHumanSingleClusterNodeList HumanSingleClusterNodeList
	// NodeChan is used Transport the k8s events
	NodeChan = make(chan map[string][]byte, 30)
	// ChanData is define the data type in the chan
	ChanData = make(map[string][]byte)
	// K8sRequest is used send the request to the k8s cluster
	K8sRequest = tools.Request{
		Methoud: "GET",
		Host:    "",
		Path:    "",
		//IsHttps shuold be https or http
		IsHTTPS:   "https",
		BearToken: "",
		Chan:      NodeChan,
	}
)

// ListAllClusterNodes is used list all cluster's node this funtiong return a map
func ListAllClusterNodes(k8sconfigs configs.HumanAllK8SConfigs) (HumanAllClusterNodeList, error) {
	for k, v1 := range k8sconfigs {
		tmp, err := ListSingleClusterNodes(v1)
		if err != nil {
			return nil, err
		}
		_HumanAllClusterNodeList[k] = tmp
	}
	return _HumanAllClusterNodeList, nil
}

// ListSingleClusterNodes is used list all node of one cluster
func ListSingleClusterNodes(k8sconfig configs.HumanSingleK8sConfigs) (HumanSingleClusterNodeList, error) {
	K8sRequest.Host = k8sconfig.EndPoint
	K8sRequest.BearToken = k8sconfig.Token
	K8sRequest.Path = "/api/v1/nodes"
	data, err := K8sRequest.Get()
	if err != nil {
		return NilHumanSingleClusterNodeList, err
	}
	err = json.Unmarshal(data, &nodelist)
	if err != nil {
		return NilHumanSingleClusterNodeList, err
	}
	healthynodestatus := make(map[string]string)
	unhealthynodestatus := make(map[string]string)
	var x, y *int
	x = new(int)
	y = new(int)
	NodeDetail(nodelist.Items, healthynodestatus, unhealthynodestatus, x, y)
	// if over imit shuold parse again
	if nodelist.Metadata.Continue != "" {
		K8sRequest.Path = "/api/v1/node?limit=500&continue=" + nodelist.Metadata.Continue
		data, err = K8sRequest.Get()
		if err != nil {
			return NilHumanSingleClusterNodeList, err
		}
		err = json.Unmarshal(data, &nodelist)
		if err != nil {
			return NilHumanSingleClusterNodeList, err
		}

		NodeDetail(nodelist.Items, healthynodestatus, unhealthynodestatus, x, y)
	}
	z := *x + *y
	var tmp2 = HumanSingleClusterNodeList{
		SingleClusterHealthyNode: healthynode{
			NodeStatus: healthynodestatus,
			Number:     x,
		},
		SingleClusterUnHealthyNode: unhealthynode{
			NodeStatus: unhealthynodestatus,
			Number:     y,
		},
		SingleClusterNodeNumber: &z,
	}
	return tmp2, nil
}

// NodeDetail is used for parse k8s api node json data
func NodeDetail(items []node, healthynodestatus, unhealthynodestatus map[string]string, x, y *int) {
	for _, v2 := range items {
		for _, v3 := range v2.Status.Conditions {
			if v3.Type == "Ready" {
				if v3.Status == "true" {
					healthynodestatus[v2.Metadata.Name] = "true"
					*x++
				} else {
					unhealthynodestatus[v2.Metadata.Name] = "false"
					*y++
				}
			}
		}
	}
}

// WatchAllClusterNodes is used watch all cluster k8s node events
func WatchAllClusterNodes(k8sconfigs configs.HumanAllK8SConfigs) {
	for k, v1 := range k8sconfigs {
		go WatchSingleClusterNodes(k, v1, ChanData)
	}
}

// WatchSingleClusterNodes is used watch single cluster k8s node evens
func WatchSingleClusterNodes(cluster string, k8sconfig configs.HumanSingleK8sConfigs, chandata map[string][]byte) {
	K8sRequest.Host = k8sconfig.EndPoint
	K8sRequest.BearToken = k8sconfig.Token
	K8sRequest.Path = "/api/v1/watch/nodes"
	K8sRequest.Watch(cluster, chandata)
}
