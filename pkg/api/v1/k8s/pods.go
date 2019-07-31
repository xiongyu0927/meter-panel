package k8s

import (
	"encoding/json"
	"meter-panel/configs"
)

var (
	podslist                      SingleClusterPodsList
	_HumanAllClusterPodsList      HumanAllClusterPodsList = make(map[string]HumanSingleClusterPodsList)
	NilHumanSingleClusterPodsList HumanSingleClusterPodsList
)

func ListAllClusterPods(k8sconfigs configs.HumanAllK8SConfigs) (HumanAllClusterPodsList, error) {
	for k, v1 := range k8sconfigs {
		tmp, err := ListSingleClusterPods(v1)
		if err != nil {
			return nil, err
		}
		_HumanAllClusterPodsList[k] = tmp
	}
	return _HumanAllClusterPodsList, nil
}

func ListSingleClusterPods(k8sconfig configs.HumanSingleK8sConfigs) (HumanSingleClusterPodsList, error) {
	K8sRequest.Host = k8sconfig.EndPoint
	K8sRequest.BearToken = k8sconfig.Token
	K8sRequest.Path = "/api/v1/pods"
	data, err := K8sRequest.Get()
	if err != nil {
		return NilHumanSingleClusterPodsList, err
	}
	err = json.Unmarshal(data, &podslist)

	if err != nil {
		return NilHumanSingleClusterPodsList, err
	}
	healthypodstatus := make(map[string]Pod, 80)
	unhealthypodstatus := make(map[string]Pod, 80)
	PodsDetail(podslist.Items, healthypodstatus, unhealthypodstatus)

	if podslist.Metadata.Continue != "" {
		K8sRequest.Path = "/api/v1/node?limit=500&continue=" + podslist.Metadata.Continue
		data, err = K8sRequest.Get()
		if err != nil {
			return NilHumanSingleClusterPodsList, err
		}
		err = json.Unmarshal(data, &podslist)
		if err != nil {
			return NilHumanSingleClusterPodsList, err
		}

		PodsDetail(podslist.Items, healthypodstatus, unhealthypodstatus)
	}

	x := len(healthypodstatus)
	y := len(unhealthypodstatus)
	z := x + y
	var tmp2 = HumanSingleClusterPodsList{
		SingleClusterHealthyPods: healthypods{
			PodStatus: healthypodstatus,
			Number:    &x,
		},
		SingleClusterUnHealthyPods: unhealthypods{
			PodStatus: unhealthypodstatus,
			Number:    &y,
		},
		SingleClusterPodsNumber: &z,
	}

	return tmp2, nil
}

func PodsDetail(items []pod, healthypodstatus, unhealthypodstatus map[string]Pod) {
	for _, v2 := range items {
		if v2.Status.Phase != "Running" && v2.Status.Phase != "Succeeded" {
			unhealthypodstatus[v2.Metadata.Name] = Pod{
				Status:       v2.Status.Phase,
				Service_name: v2.Metadata.Labels.Service_name,
				Apps:         v2.Metadata.Labels.Apps,
			}
		} else {
			healthypodstatus[v2.Metadata.Name] = Pod{
				Status:       v2.Status.Phase,
				Service_name: v2.Metadata.Labels.Service_name,
				Apps:         v2.Metadata.Labels.Apps,
			}
		}
	}
}
