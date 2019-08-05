package store

import (
	"log"
	"meter-panel/pkg/api/v1/k8s"
)

func GetPv(cluster string) k8s.HumanSingleClusterPvsList {
	singlek8sconfig := StoreAllK8SConfigs.GetSingleK8SConfigs(cluster)
	tmp, err := k8s.ListSingleClusterPvs(singlek8sconfig)
	if err != nil {
		log.Println(err)
	}
	return tmp
}
