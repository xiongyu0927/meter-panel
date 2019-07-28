package configs

import (
	"encoding/json"
	"fmt"
	"log"
	"meter-panel/tools"
	"strings"
)

// AllK8SConfigs is a map, key is cluster name, value is a struct that include k8s EndPoint and BearToken
var AllK8SConfigs HumanAllK8SConfigs = make(map[string]HumanSingleK8sConfigs)

// GetK8SCoinfg is used get a all k8s cluster config
func GetK8SCoinfg() (HumanAllK8SConfigs, error) {
	var furionrequest = tools.Request{
		Methoud: "GET",
		Host:    "furion:8080",
		Path:    "/v1/regions",
		//IsHttps shuold be https or http
		IsHTTPS:   "http",
		BearToken: "",
	}

	furionhost := GetSingleEnvConfigs("FurionHost")
	if furionhost == "" {
		log.Println("enviroment doesn't set, will use default config furion:8080")
	} else {
		furionrequest.Host = furionhost
	}

	data, err := furionrequest.Get()
	if err != nil {
		return nil, err
	}

	var k8sconfig []K8sconfigs
	err = json.Unmarshal(data, &k8sconfig)
	if err != nil {
		return nil, err
	}

	for _, v := range k8sconfig {
		endpoint := strings.SplitN(v.Attr.Kubernetes.Endpoint, "//", -1)[1]
		fmt.Println("hey " + endpoint)
		var tmp = HumanSingleK8sConfigs{
			EndPoint: endpoint,
			Token:    v.Attr.Kubernetes.Token,
		}
		AllK8SConfigs[v.Name] = tmp
	}

	return AllK8SConfigs, nil
}

func (self HumanAllK8SConfigs) GetSingleK8SConfigs(cluster string) HumanSingleK8sConfigs {
	return self[cluster]
}
