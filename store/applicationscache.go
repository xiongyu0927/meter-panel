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

func AppModifyed(cluster string, poddetail map[string]k8s.Pod, podname string) {
	for _, v := range poddetail {
		if v.Appredis != "" {
			thisL := v.Appredis
			log.Println("Redis " + thisL)
			Appplace, place, Appname := GetAppPlace(cluster, thisL, "R")
			// log.Println(Appplace)
			if Appplace == NilK8SPod {
				// 新增带label的pod，Application也无法获取只能重新list一遍
				StoreAllClusterAppList, err = k8s.ListAllClusterApplications(StoreAllK8SConfigs, StoreAllClusterPodList)
				log.Println("List a app, maybe have new app add")
				return
			}
			modifyapp(Appplace, cluster, place, Appname)
			return
		}

		if v.Apps != "" {
			thisL := v.Apps
			log.Println("hi " + thisL)
			Appplace, place, Appname := GetAppPlace(cluster, thisL, "A")
			// log.Println(Appplace)
			if Appplace == NilK8SPod {
				// 新增带label的pod，Application也无法获取只能重新list一遍
				StoreAllClusterAppList, err = k8s.ListAllClusterApplications(StoreAllK8SConfigs, StoreAllClusterPodList)
				log.Println("List a app, maybe have new app add")
				return
			}
			modifyapp(Appplace, cluster, place, Appname)
			return
		}

		if v.Service_name != "" {
			thisL := poddetail[podname].Service_name
			log.Println("hey " + thisL)
			Appplace, place, Appname := GetAppPlace(cluster, thisL, "S")
			log.Println(Appplace)
			if Appplace == NilK8SPod {
				// 新增带label的pod，Application也无法获取只能重新list一遍
				StoreAllClusterAppList, err = k8s.ListAllClusterApplications(StoreAllK8SConfigs, StoreAllClusterPodList)
				log.Println("List a app, maybe have new app add")
				return
			}
			modifyapp(Appplace, cluster, place, Appname)
			return
		}
	}
}

func GetAppPlace(cluster, label, labelT string) (tmp k8s.Pod, place string, appname string) {
	if labelT == "R" {
		for k, v := range StoreAllClusterAppList[cluster].SingleClusterHealthyApp.AppStatus {
			if label == v.Appredis {
				return v, "H", k
			}
		}

		for k, v := range StoreAllClusterAppList[cluster].SingleClusterUnHealthyApp.AppStatus {
			if label == v.Appredis {
				return v, "UH", k
			}
		}
	}

	if labelT == "A" {
		for k, v := range StoreAllClusterAppList[cluster].SingleClusterHealthyApp.AppStatus {
			if label == v.Apps {
				return v, "H", k
			}
		}

		for k, v := range StoreAllClusterAppList[cluster].SingleClusterUnHealthyApp.AppStatus {
			if label == v.Apps {
				return v, "UH", k
			}
		}
	}

	if labelT == "S" {
		for k, v := range StoreAllClusterAppList[cluster].SingleClusterHealthyApp.AppStatus {
			if label == v.Service_name {
				return v, "H", k
			}
		}

		for k, v := range StoreAllClusterAppList[cluster].SingleClusterUnHealthyApp.AppStatus {
			if label == v.Service_name {
				return v, "UH", k
			}
		}
	}
	return NilK8SPod, "", ""
}

func modifyapp(Appplace k8s.Pod, cluster string, place string, Appname string) {
	for _, v := range StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.PodStatus {
		if Appplace.Appredis == v.Appredis {
			if place == "UH" {
			} else {
				delete(StoreAllClusterAppList[cluster].SingleClusterHealthyApp.AppStatus, Appname)
				*StoreAllClusterAppList[cluster].SingleClusterHealthyApp.Number--
				StoreAllClusterAppList[cluster].SingleClusterUnHealthyApp.AppStatus[Appname] = k8s.Pod{
					Status:       "Processing",
					Apps:         Appplace.Apps,
					Service_name: Appplace.Service_name,
					Appredis:     Appplace.Appredis,
				}
				*StoreAllClusterAppList[cluster].SingleClusterUnHealthyApp.Number++
				log.Println(Appname + "become unhealthy")
				return
			}
		}
	}

	for _, v := range StoreAllClusterPodList[cluster].SingleClusterUnHealthyPods.PodStatus {
		if Appplace.Apps == v.Apps {
			if place == "UH" {
			} else {
				delete(StoreAllClusterAppList[cluster].SingleClusterHealthyApp.AppStatus, Appname)
				*StoreAllClusterAppList[cluster].SingleClusterHealthyApp.Number--
				StoreAllClusterAppList[cluster].SingleClusterUnHealthyApp.AppStatus[Appname] = k8s.Pod{
					Status:       "Processing",
					Apps:         Appplace.Apps,
					Service_name: Appplace.Service_name,
					Appredis:     Appplace.Appredis,
				}
				*StoreAllClusterAppList[cluster].SingleClusterUnHealthyApp.Number++
				log.Println(Appname + "become unhealthy")
				return
			}
		}
	}

	for _, v := range StoreAllClusterPodList[cluster].SingleClusterHealthyPods.PodStatus {
		if Appplace.Service_name == v.Service_name {
			if place == "UH" {
				delete(StoreAllClusterAppList[cluster].SingleClusterUnHealthyApp.AppStatus, Appname)
				*StoreAllClusterAppList[cluster].SingleClusterUnHealthyApp.Number--
				StoreAllClusterAppList[cluster].SingleClusterHealthyApp.AppStatus[Appname] = k8s.Pod{
					Status:       "Running",
					Apps:         Appplace.Apps,
					Service_name: Appplace.Service_name,
					Appredis:     Appplace.Appredis,
				}
				*StoreAllClusterAppList[cluster].SingleClusterHealthyApp.Number++
				log.Println(Appname + "become healthy")
				return
			}
		}
	}

	// 根据pod的状态来判断application的状态，当pod为0时，无法知道application是停止了还是删除了，所以重新list一遍
	StoreAllClusterAppList, err = k8s.ListAllClusterApplications(StoreAllK8SConfigs, StoreAllClusterPodList)
	log.Println("List a app, maybe have new app delete")
	return
}
