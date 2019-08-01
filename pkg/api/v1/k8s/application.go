package k8s

import (
	"encoding/json"
	"log"
	"meter-panel/configs"
)

var (
	applist                          SingleClusterApplicationsList
	_HumanAllClusterApplicationsList HumanAllClusterApplicationsList = make(map[string]HumanSingleClusterApplicationsList)
	NilHumanSingleClusterAppList     HumanSingleClusterApplicationsList
)

func ListAllClusterApplications(k8sconfigs configs.HumanAllK8SConfigs, allpod HumanAllClusterPodsList) (HumanAllClusterApplicationsList, error) {
	for k, v1 := range k8sconfigs {
		tmp, err := ListSingleClusterApplications(v1, allpod[k])
		if err != nil {
			return nil, err
		}
		_HumanAllClusterApplicationsList[k] = tmp
	}
	return _HumanAllClusterApplicationsList, nil
}

func ListSingleClusterApplications(k8sconfig configs.HumanSingleK8sConfigs, oneclusterpods HumanSingleClusterPodsList) (HumanSingleClusterApplicationsList, error) {
	K8sRequest.Host = k8sconfig.EndPoint
	K8sRequest.BearToken = k8sconfig.Token
	K8sRequest.Path = "/apis/app.k8s.io/v1beta1/applications?limit=500"
	data, err := K8sRequest.Get()
	if err != nil {
		return NilHumanSingleClusterAppList, err
	}
	err = json.Unmarshal(data, &applist)
	if err != nil {
		return NilHumanSingleClusterAppList, err
	}
	log.Println(applist.Items)
	// 处理数据得到一个人类的Application
	healthyappstatus := make(map[string]Pod, 40)
	unhealthyappstatus := make(map[string]Pod, 40)
	AppsDetail(applist.Items, healthyappstatus, unhealthyappstatus, oneclusterpods)

	if applist.Metadata.Continue != "" {
		K8sRequest.Path = "/apis/app.k8s.io/v1beta1/applications?limit=500&continue=" + applist.Metadata.Continue
		data, err = K8sRequest.Get()
		if err != nil {
			return NilHumanSingleClusterAppList, err
		}
		err = json.Unmarshal(data, &applist)
		if err != nil {
			return NilHumanSingleClusterAppList, err
		}
		AppsDetail(applist.Items, healthyappstatus, unhealthyappstatus, oneclusterpods)
	}

	x := len(healthyappstatus)
	y := len(unhealthyappstatus)
	z := x + y
	var tmp2 = HumanSingleClusterApplicationsList{
		SingleClusterHealthyApp: healthyapp{
			AppStatus: healthyappstatus,
			Number:    &x,
		},
		SingleClusterUnHealthyApp: unhealthyapp{
			AppStatus: unhealthyappstatus,
			Number:    &y,
		},
		SingleClusterAppNumber: &z,
	}

	return tmp2, nil
}

func AppsDetail(items []app, healthyappstatus, unhealthyappstatus map[string]Pod, oneclusterpods HumanSingleClusterPodsList) {
LABEL1:
	for _, v := range items {
		for _, v1 := range oneclusterpods.SingleClusterUnHealthyPods.PodStatus {
			if v.Spec.Selector.MatchLabels.Apps == v1.Apps || v.Spec.Selector.MatchLabels.Service_name == v1.Service_name {
				unhealthyappstatus[v.Metadata.Name] = Pod{
					Status:       "Processing",
					Apps:         v.Spec.Selector.MatchLabels.Apps,
					Service_name: v.Spec.Selector.MatchLabels.Service_name,
				}
				continue LABEL1
			}
		}

		for _, v2 := range oneclusterpods.SingleClusterHealthyPods.PodStatus {
			if v.Spec.Selector.MatchLabels.Apps == v2.Apps || v.Spec.Selector.MatchLabels.Service_name == v2.Service_name {
				healthyappstatus[v.Metadata.Name] = Pod{
					Status:       "Running",
					Apps:         v.Spec.Selector.MatchLabels.Apps,
					Service_name: v.Spec.Selector.MatchLabels.Service_name,
				}
				log.Println("i'm here")
				continue LABEL1
			}
			healthyappstatus[v.Metadata.Name] = Pod{
				Status:       "Stop",
				Apps:         v.Spec.Selector.MatchLabels.Apps,
				Service_name: v.Spec.Selector.MatchLabels.Service_name,
			}
			continue LABEL1
		}
	}
}
