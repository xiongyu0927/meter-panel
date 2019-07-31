package store

import (
	"log"
	"meter-panel/pkg/api/v1/k8s"
	"reflect"
)

var NilK8SPod k8s.Pod

func GetSingleClusterPodsList(cluster string) k8s.HumanSingleClusterPodsList {
	data := StoreAllClusterPodList[cluster]
	if !reflect.DeepEqual(data, NilSlingeClusterPodList) {
		return data
	}
	singlek8sconfig := StoreAllK8SConfigs.GetSingleK8SConfigs(cluster)
	tmp, err := k8s.ListSingleClusterPods(singlek8sconfig)
	if err != nil {
		log.Println(err)
	}
	StoreAllClusterPodList[cluster] = tmp
	log.Println(cluster + " is seted in podcache")
	return tmp
}

func PodModifyed(cluster string, poddetail map[string]k8s.Pod, podname string, eventtype string) {
	// delete pods
	if eventtype == "DELETE" {
		if StoreAllClusterPodList[cluster].SingleClusterHealthyPods.PodStatus[podname] == NilK8SPod {
			delete(StoreAllClusterPodList[cluster].SingleClusterHealthyPods.PodStatus, podname)
			*StoreAllClusterPodList[cluster].SingleClusterHealthyPods.Number--
			return
		}
		delete(StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.PodStatus, podname)
		*StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.Number--
		return
	}

	// add pods
	if StoreAllClusterPodList[cluster].SingleClusterHealthyPods.PodStatus[podname] == NilK8SPod &&
		StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.PodStatus[podname] == NilK8SPod {
		if poddetail[podname].Status == "Running" || poddetail[podname].Status == "Succeeded" {
			StoreAllClusterPodList[cluster].SingleClusterHealthyPods.PodStatus[podname] = poddetail[podname]
			*StoreAllClusterPodList[cluster].SingleClusterHealthyPods.Number++
			return
		}
		StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.PodStatus[podname] = poddetail[podname]
		*StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.Number++
		return
	}

	// change pod
	if StoreAllClusterPodList[cluster].SingleClusterHealthyPods.PodStatus[podname] == NilK8SPod {
		if StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.PodStatus[podname] == poddetail[podname] {
			// log.Println("the Pod status of cluster " + cluster + " doesn't has changed")
		} else {
			delete(StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.PodStatus, podname)
			StoreAllClusterPodList[cluster].SingleClusterHealthyPods.PodStatus[podname] = poddetail[podname]
			*StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.Number--
			*StoreAllClusterPodList[cluster].SingleClusterHealthyPods.Number++
			log.Println("the Pod status cluster " + cluster + "was changed")
			return
		}
	} else {
		if StoreAllClusterPodList[cluster].SingleClusterHealthyPods.PodStatus[podname] == poddetail[podname] {
			// log.Println("the Pod status of cluster " + cluster + " doesn't has changed")
		} else {
			delete(StoreAllClusterPodList[cluster].SingleClusterHealthyPods.PodStatus, podname)
			StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.PodStatus[podname] = poddetail[podname]
			*StoreAllClusterPodList[cluster].SingleClusterHealthyPods.Number--
			*StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.Number++
			log.Println("the Pod status cluster " + cluster + "was changed")
			return
		}
	}
}
