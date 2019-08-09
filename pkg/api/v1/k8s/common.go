package k8s

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"meter-panel/configs"
	"meter-panel/tools"
	"net/http"
)

// WatchAllClusterNodes is used watch all cluster k8s node events
func WatchAllClusterResource(k8sconfigs configs.HumanAllK8SConfigs, resourece string) {
	for k, v1 := range k8sconfigs {
		WatchSingleClusterResource(k, v1, resourece)
	}
}

// WatchSingleClusterNodes is used watch single cluster k8s node evens
func WatchSingleClusterResource(cluster string, k8sconfig configs.HumanSingleK8sConfigs, resource string) {
	K8sRequest.Host = k8sconfig.EndPoint
	K8sRequest.BearToken = k8sconfig.Token
	switch resource {
	case "nodes":
		K8sRequest.Path = "/api/v1/watch/nodes"
	case "pods":
		K8sRequest.Path = "/api/v1/watch/pods"
	}

	go Watch(cluster, K8sRequest, resource)
}

// Watch is used to watch k8s resource
func Watch(cluster string, its tools.Request, resource string) {
	request, err := http.NewRequest("GET", its.IsHTTPS+"://"+its.Host+its.Path, nil)
	if err != nil {
		log.Println(err)
	}

	var tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// add BearToken auth
	if its.BearToken != "" {
		request.Header.Add("Authorization", "Bearer "+its.BearToken)
	}

	client := http.Client{}
	// add InsecureSkipVerify
	if its.IsHTTPS == "https" {
		client.Transport = tr
	}

	// execute this request
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}

	// get,read and return
	defer resp.Body.Close()

	DecodeJson(cluster, resp, resource)
}

var nilpodevent PodEvents

func DecodeJson(cluster string, resp *http.Response, resource string) {
	switch resource {
	case "nodes":
		for {
			var a NodeEvents
			chandata := make(map[string]interface{})
			err := json.NewDecoder(resp.Body).Decode(&a)
			if err != nil {
				log.Println(err)
			}
			chandata[cluster] = a
			K8SChan <- chandata
		}
	case "pods":
		for {
			var b PodEvents
			chandata := make(map[string]interface{})
			_ = json.NewDecoder(resp.Body).Decode(&b)
			// if err != nil {
			// 	log.Println(err)
			// }
			if b != nilpodevent {
				chandata[cluster] = b
				K8SChan <- chandata
			}
		}
	}
}

func ListAllClusterProCfg(allp HumanAllClusterPodsList, allc configs.HumanAllK8SConfigs) (map[string]string, error) {
	tmp := make(map[string]string, 4)

	tmp = configs.GetProConfig4ENV(allc)
	if tmp != nil {
		return tmp, nil
	}

	tmp = GetProCfg4LocalPodCache(allp)
	if tmp != nil {
		return tmp, nil
	}

	return nil, errors.New("can't load promethues cfg from anyway")
}

// GetProCfg4LocalPodCache is used when the env doesn't set
func GetProCfg4LocalPodCache(allp HumanAllClusterPodsList) map[string]string {
	ProCfg := make(map[string]string, 4)
	for k1, v1 := range allp {
		for k2, v2 := range v1.SingleClusterHealthyPods.PodStatus {
			if k2 == "prometheus-kube-prometheus-0" {
				ProCfg[k1] = v2.PodIp
			}
		}
	}
	return ProCfg
}
