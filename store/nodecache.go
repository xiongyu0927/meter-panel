package store

import (
	"log"
	k8s "meter-panel/pkg/api/v1/k8s"
	"reflect"
)

// GetSingleClusterNodeList is used query node info of a cluster
// when a new cluster create in furion, this function can refresh the info of cluster
func GetSingleClusterNodeList(cluster string) k8s.HumanSingleClusterNodeList {
	data := StoreAllClusterNodeList[cluster]
	if !reflect.DeepEqual(data, NilSingleClusterNodeList) {
		return data
	}
	singlek8sconfig := StoreAllK8SConfigs.GetSingleK8SConfigs(cluster)
	tmp, err := k8s.ListSingleClusterNodes(singlek8sconfig)
	if err != nil {
		log.Println(err)
	}
	StoreAllClusterNodeList[cluster] = tmp
	log.Println(cluster + " is seted in nodecache")
	return tmp
}

// NodeModifyed is used when node event type is Modifyed
func NodeModifyed(cluster string, nodedetail map[string]string, nodename string) {
	if StoreAllClusterNodeList[cluster].SingleClusterHealthyNode.NodeStatus[nodename] == "" {
		if StoreAllClusterNodeList[cluster].SingleClusterUnHealthyNode.NodeStatus[nodename] == nodedetail[nodename] {
			// log.Println("the node status of cluster " + cluster + " doesn't has changed")
		} else {
			delete(StoreAllClusterNodeList[cluster].SingleClusterUnHealthyNode.NodeStatus, nodename)
			StoreAllClusterNodeList[cluster].SingleClusterHealthyNode.NodeStatus[nodename] = nodedetail[nodename]
			*StoreAllClusterNodeList[cluster].SingleClusterUnHealthyNode.Number--
			log.Println("the node status cluster " + cluster + "was changed")
		}
	} else {
		if StoreAllClusterNodeList[cluster].SingleClusterHealthyNode.NodeStatus[nodename] == nodedetail[nodename] {
			// log.Println("the node status of cluster " + cluster + " doesn't has changed")
		} else {
			delete(StoreAllClusterNodeList[cluster].SingleClusterHealthyNode.NodeStatus, nodename)
			StoreAllClusterNodeList[cluster].SingleClusterUnHealthyNode.NodeStatus[nodename] = nodedetail[nodename]
			*StoreAllClusterNodeList[cluster].SingleClusterUnHealthyNode.Number--
			log.Println("the node status cluster " + cluster + "was changed")
		}
	}
}
