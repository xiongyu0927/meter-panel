package store

import (
	"meter-panel/pkg/api/v1/k8s"
)

func GetSingleClusterLBList(cluster string) k8s.HumanSingleClusterLbsList {
	singlek8sconfig := StoreAllK8SConfigs.GetSingleK8SConfigs(cluster)
	tmp := k8s.ListSingleClusterLB(singlek8sconfig)
	return tmp
}
