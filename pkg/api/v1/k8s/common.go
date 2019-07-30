package k8s

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"meter-panel/configs"
	"meter-panel/tools"
	"net/http"
)

// WatchAllClusterNodes is used watch all cluster k8s node events
func WatchAllClusterResource(k8sconfigs configs.HumanAllK8SConfigs, resourece string) {
	for k, v1 := range k8sconfigs {
		go WatchSingleClusterResource(k, v1, resourece)
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

	Watch(cluster, K8sRequest, resource)
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

	DecodeJson(cluster, resp, ChanData, resource)
}

func DecodeJson(cluster string, resp *http.Response, chandata map[string]interface{}, resource string) {
	switch resource {
	case "nodes":
		var a NodeEvents
		for {
			err := json.NewDecoder(resp.Body).Decode(&a)
			if err != nil {
				log.Println(err)
			}
			chandata[cluster] = a
			K8SChan <- chandata
		}
	case "pods":
		var a PodEvents
		for {
			err := json.NewDecoder(resp.Body).Decode(&a)
			if err != nil {
				log.Println(err)
			}
			chandata[cluster] = a
			K8SChan <- chandata
		}
	}
}