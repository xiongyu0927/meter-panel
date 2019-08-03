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
	// StoreAllClusterPodslist is used save all cluster pod list in the memory
	StoreAllClusterPodList k8s.HumanAllClusterPodsList
	// StoreAllClusterAppList is used save all cluster app list in the memory
	StoreAllClusterAppList k8s.HumanAllClusterApplicationsList
	// NilSingleClusterNodeList is used return nil value of HumanSingleClusterNodeList
	NilSingleClusterNodeList k8s.HumanSingleClusterNodeList
	NilSlingeClusterPodList  k8s.HumanSingleClusterApplicationsList
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

	StoreAllClusterPodList, err = k8s.ListAllClusterPods(StoreAllK8SConfigs)
	if err != nil {
		log.Println(err)
	}

	StoreAllClusterAppList, err = k8s.ListAllClusterApplications(StoreAllK8SConfigs, StoreAllClusterPodList)

	k8s.WatchAllClusterResource(StoreAllK8SConfigs, "nodes")
	k8s.WatchAllClusterResource(StoreAllK8SConfigs, "pods")
}

func init() {
	go func() {
		for {
			tmp := <-k8s.K8SChan
			for k, v := range tmp {
				switch x := v.(type) {
				case k8s.NodeEvents:
					nodedetail := make(map[string]string)
					nodename := x.Object.Metadata.Name
					for _, v2 := range x.Object.Status.Conditions {
						if v2.Type == "Ready" {
							nodedetail[nodename] = v2.Status
						}
					}
					NodeModifyed(k, nodedetail, nodename)
				case k8s.PodEvents:
					poddetail := make(map[string]k8s.Pod)
					eventtype := x.Type
					podname := x.Object.Metadata.Name
					poddetail[podname] = k8s.Pod{
						Status:       x.Object.Status.Phase,
						Service_name: x.Object.Metadata.Labels.Service_name,
						Apps:         x.Object.Metadata.Labels.Apps,
						PodIp:        x.Object.Status.PodIp,
					}
					PodModifyed(k, poddetail, podname, eventtype)
					// AppModifyed(k, poddetail, podname, eventtype)
				}
			}
		}
	}()
}
