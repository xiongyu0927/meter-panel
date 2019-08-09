package store

import (
	"log"
	"meter-panel/configs"
	"meter-panel/pkg/api/v1/k8s"
	"reflect"
)

var NilK8SPod k8s.Pod

func GetSingleClusterPodsList(cluster string) k8s.HumanSingleClusterPodsList {
	data := StoreAllClusterPodList[cluster]
	if !reflect.DeepEqual(data, NilSlingeClusterPodList) {
		return data
	}
	StoreAllK8SConfigs, _ = configs.GetK8SCoinfg()
	log.Println("maybe have new cluster join")
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
	if eventtype == "DELETED" {
		if StoreAllClusterPodList[cluster].SingleClusterHealthyPods.PodStatus[podname] == NilK8SPod {
			delete(StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.PodStatus, podname)
			*StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.Number--
			*StoreAllClusterPodList[cluster].SingleClusterPodsNumber--
			log.Println("delete a unhealthy pods " + podname + " in " + cluster)
			return
		}
		delete(StoreAllClusterPodList[cluster].SingleClusterHealthyPods.PodStatus, podname)
		*StoreAllClusterPodList[cluster].SingleClusterHealthyPods.Number--
		*StoreAllClusterPodList[cluster].SingleClusterPodsNumber--
		log.Println("delete a healthy pods " + podname + " in " + cluster)
		return
	}

	// add pods
	if StoreAllClusterPodList[cluster].SingleClusterHealthyPods.PodStatus[podname] == NilK8SPod &&
		StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.PodStatus[podname] == NilK8SPod {
		if poddetail[podname].Status == "Running" || poddetail[podname].Status == "Succeeded" {
			StoreAllClusterPodList[cluster].SingleClusterHealthyPods.PodStatus[podname] = poddetail[podname]
			*StoreAllClusterPodList[cluster].SingleClusterHealthyPods.Number++
			*StoreAllClusterPodList[cluster].SingleClusterPodsNumber++
			log.Println("add a healthy pod " + podname + "in " + cluster)
			return
		}
		StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.PodStatus[podname] = poddetail[podname]
		*StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.Number++
		*StoreAllClusterPodList[cluster].SingleClusterPodsNumber++
		return
	}

	// change pod
	if StoreAllClusterPodList[cluster].SingleClusterHealthyPods.PodStatus[podname] == NilK8SPod {
		if poddetail[podname].Status == "Running" || poddetail[podname].Status == "Succeeded" {
			// log.Println("the Pod status of cluster " + cluster + " doesn't has changed")
			delete(StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.PodStatus, podname)
			StoreAllClusterPodList[cluster].SingleClusterHealthyPods.PodStatus[podname] = poddetail[podname]
			*StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.Number--
			*StoreAllClusterPodList[cluster].SingleClusterHealthyPods.Number++
			log.Println("change a unhealthy pod " + podname + "to healthy in " + cluster)
			return
		}
	}

	if StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.PodStatus[podname] == NilK8SPod {
		if poddetail[podname].Status != "Running" && poddetail[podname].Status != "Succeeded" {
			// log.Println("the Pod status of cluster " + cluster + " doesn't has changed")
			delete(StoreAllClusterPodList[cluster].SingleClusterHealthyPods.PodStatus, podname)
			StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.PodStatus[podname] = poddetail[podname]
			*StoreAllClusterPodList[cluster].SingleClusterHealthyPods.Number--
			*StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.Number++
			log.Println("change a healthy pod " + podname + "to unhealthy in " + cluster)
			return
		}
	}
}

// ProAddrModified is used for watch every prometheus's address in every cluster
func ProAddrModified(cluster string, poddetail map[string]k8s.Pod, podname string) {
	if podname == "prometheus-kube-prometheus-0" {
		if StoreAllProConfigs[cluster] != poddetail[podname].PodIp {
			StoreAllProConfigs[cluster] = poddetail[podname].PodIp
		}
	}
}
