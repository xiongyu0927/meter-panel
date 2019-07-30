package k8s

import (
	"encoding/json"
	"meter-panel/configs"
)

var (
	applist                          SingleClusterApplicationsList
	_HumanAllClusterApplicationsList HumanAllClusterApplicationsList = make(map[string]HumanSingleClusterApplicationsList)
	NilHumanSingleClusterAppList     HumanSingleClusterApplicationsList
)

func ListAllClusterApplications(k8sconfigs configs.HumanAllK8SConfigs) (HumanAllClusterApplicationsList, error) {
	for k, v1 := range k8sconfigs {
		tmp, err := ListSingleClusterApplications(v1)
		if err != nil {
			return nil, err
		}
		_HumanAllClusterApplicationsList[k] = tmp
	}
	return _HumanAllClusterApplicationsList, nil
}

func ListSingleClusterApplications(k8sconfig configs.HumanSingleK8sConfigs) (HumanSingleClusterApplicationsList, error) {
	K8sRequest.Host = k8sconfig.EndPoint
	K8sRequest.BearToken = k8sconfig.Token
	K8sRequest.Path = "/apis/app.k8s.io/v1beta1/applications"
	data, err := K8sRequest.Get()
	if err != nil {
		return NilHumanSingleClusterAppList, err
	}
	err = json.Unmarshal(data, &applist)
	if err != nil {
		return NilHumanSingleClusterAppList, err
	}
	// 处理数据得到一个人类的Application
	return NilHumanSingleClusterAppList, nil
}
