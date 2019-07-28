package store

import (
	"encoding/json"
	"log"
	"meter-panel/configs"
	k8s "meter-panel/pkg/api/v1/k8s"
	"reflect"
)

var (
	//StoreAllClusterNodeList is used save all cluster node list in the memory
	StoreAllClusterNodeList k8s.HumanAllClusterNodeList
	// StoreAllK8SConfigs is used save all cluster configs in the memory
	StoreAllK8SConfigs configs.HumanAllK8SConfigs
	// NilSingleClusterNodeList is used return nil value of HumanSingleClusterNodeList
	NilSingleClusterNodeList k8s.HumanSingleClusterNodeList
	err                      error
)

func init() {
	StoreAllK8SConfigs, err = configs.GetK8SCoinfg()
	if err != nil {
		log.Println(err)
	}

	StoreAllClusterNodeList, err = k8s.ListAllClusterNodes(StoreAllK8SConfigs)
	if err != nil {
		log.Println(err)
	}

	k8s.WatchAllClusterNodes(StoreAllK8SConfigs)
}

func init() {
	var event k8s.NodeEvents
	go func() {
		for {
			tmp := <-k8s.K8sRequest.Chan
			nodedetail := make(map[string]string)
			for k, v := range tmp {
				err := json.Unmarshal(v, &event)
				if err != nil {
					log.Println("there is a error?")
					log.Println(err)
				}
				switch event.Type {
				case "MODIFIED":
					nodename := event.Object.Metadata.Name
					for _, v := range event.Object.Status.Conditions {
						if v.Type == "Ready" {
							nodedetail[nodename] = v.Status
						}
					}
					Modifyed(k, nodedetail, nodename)
					// case ""：
				}
			}
		}
	}()
}

// GetSingleClusterNodeList is used query node info of a cluster
// when a new cluster create in furion, this function can refresh the info of cluster
func (its *NodeCache) GetSingleClusterNodeList(cluster string) k8s.HumanSingleClusterNodeList {
	singlek8sconfig := StoreAllK8SConfigs.GetSingleK8SConfigs(cluster)
	data := its.StoreAllClusterNodeList[cluster]
	if !reflect.DeepEqual(data, NilSingleClusterNodeList) {
		return data
	}
	tmp, err := k8s.ListSingleClusterNodes(singlek8sconfig)
	if err != nil {
		log.Println(err)
	}
	its.StoreAllClusterNodeList[cluster] = tmp
	log.Println(cluster + " is seted in nodecache")
	return tmp
}

// Modifyed is used when node event type is Modifyed
func Modifyed(cluster string, nodedetail map[string]string, nodename string) {
	if StoreAllClusterNodeList[cluster].SingleClusterHealthyNode.NodeStatus[nodename] == "" {
		if StoreAllClusterNodeList[cluster].SingleClusterUnHealthyNode.NodeStatus[nodename] == nodedetail[nodename] {
			log.Println("the node status of cluster " + cluster + " doesn't has changed")
		} else {
			delete(StoreAllClusterNodeList[cluster].SingleClusterUnHealthyNode.NodeStatus, nodename)
			StoreAllClusterNodeList[cluster].SingleClusterHealthyNode.NodeStatus[nodename] = nodedetail[nodename]
			*StoreAllClusterNodeList[cluster].SingleClusterUnHealthyNode.Number--
			log.Println("the node status cluster " + cluster + "was changed")
		}
	} else {
		if StoreAllClusterNodeList[cluster].SingleClusterHealthyNode.NodeStatus[nodename] == nodedetail[nodename] {
			log.Println("the node status of cluster " + cluster + " doesn't has changed")
		} else {
			delete(StoreAllClusterNodeList[cluster].SingleClusterHealthyNode.NodeStatus, nodename)
			StoreAllClusterNodeList[cluster].SingleClusterUnHealthyNode.NodeStatus[nodename] = nodedetail[nodename]
			*StoreAllClusterNodeList[cluster].SingleClusterUnHealthyNode.Number--
			log.Println("the node status cluster " + cluster + "was changed")
		}
	}
}
