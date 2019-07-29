package store

import (
	"log"
	"meter-panel/configs"
	"meter-panel/pkg/api/v1/k8s"
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
			tmp := <-k8s.K8SChan
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
					NodeModifyed(k, nodedetail, nodename)
				}
			}
		}
	}()
}
