package store

import (
	"log"
	"meter-panel/configs"
	"meter-panel/pkg/api/v1/k8s"
	"reflect"
)

func GetSingleClusterLBList(cluster string) k8s.HumanSingleClusterLbsList {
	singlek8sconfig := StoreAllK8SConfigs.GetSingleK8SConfigs(cluster)
	if reflect.DeepEqual(singlek8sconfig, NilK8sSingleClusterConfigs) {
		StoreAllK8SConfigs, _ = configs.GetK8SCoinfg()
		singlek8sconfig = StoreAllK8SConfigs.GetSingleK8SConfigs(cluster)
		log.Println("maybe have new cluster join")
	}
	tmp := k8s.ListSingleClusterLB(singlek8sconfig)
	return tmp
}
