package store

import (
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
	nodedetail := make(map[string]string)
	go func() {
		for {
			tmp := <-k8s.K8sRequest.Chan
			for k, v := range tmp {
				// k was cluster name
				// 处理序列化后的数据
				switch x := v.(type) {
				case k8s.NodeEvents:
					nodename := x.Object.Metadata.Name
					for _, v2 := range x.Object.Status.Conditions {
						if v2.Type == "Ready" {
							nodedetail[nodename] = v2.Status
						}
					}
					Modifyed(k, nodedetail, nodename)
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
