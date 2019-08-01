package store

import (
	"log"
	"meter-panel/pkg/api/v1/k8s"
	"reflect"
)

var NilK8SApp k8s.HumanSingleClusterApplicationsList

func GetSingleClusterAppsList(cluster string) k8s.HumanSingleClusterApplicationsList {
	data := StoreAllClusterAppList[cluster]
	if !reflect.DeepEqual(data, NilK8SApp) {
		return data
	}
	singlek8sconfig := StoreAllK8SConfigs.GetSingleK8SConfigs(cluster)
	tmp, err := k8s.ListSingleClusterApplications(singlek8sconfig, StoreAllClusterPodList[cluster])
	if err != nil {
		log.Println(err)
	}
	StoreAllClusterAppList[cluster] = tmp
	log.Println(cluster + " is seted in APPcache")
	return tmp
}

func AppModifyed(cluster string, poddetail map[string]k8s.Pod, podname string, eventtype string) {

	for _, v := range StoreAllClusterAppList[cluster].SingleClusterUnHealthyApp.AppStatus {
		if poddetail[podname].Apps == v.Apps || poddetail[podname].Service_name == v.Service_name {
			for _, v1 := range StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.PodStatus {
				if v.Apps == v1.Apps || v.Service_name == v1.Service_name {
					v.Status = "Processing"
					return
				}
			}

			for _, v2 := range StoreAllClusterPodList[cluster].SingleClusterHealthyPods.PodStatus {
				if v.Apps == v2.Apps || v.Service_name == v2.Service_name {
					v.Status = "Running"
					return
				}
			}

			// 根据pod的状态来判断application的状态，当pod为0时，无法知道application是停止了还是删除了，所以重新list一遍
			StoreAllClusterAppList, err = k8s.ListAllClusterApplications(StoreAllK8SConfigs, StoreAllClusterPodList)
			log.Println("List a app, maybe have new app delete or add")
			return
		}

		// 新增带label的pod，Application也无法获取只能重新list一遍
		StoreAllClusterAppList, err = k8s.ListAllClusterApplications(StoreAllK8SConfigs, StoreAllClusterPodList)
		log.Println("List a app, maybe have new app delete or add")
	}

}
