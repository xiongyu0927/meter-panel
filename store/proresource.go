package store

import (
	"log"
	"os"
	"strings"

	"meter-panel/configs"
	"meter-panel/tools"
)

var ProRequest = tools.Request{
	Methoud: "GET",
	Host:    "",
	Path:    "",
	//IsHttps shuold be https or http
	IsHTTPS:   "http",
	BearToken: "",
}

const (
	cpuapi   string = "/api/v1/query?query=cluster_cpu_utilization%20%2F%20count(avg%20by(instance)(node_cpu%7Bjob%3D%22node-exporter%22%2Cmode%3D%22idle%22%7D))%20*%20100"
	memapi   string = "/api/v1/query?query=((sum(node_memory_MemTotal)%20-%20sum(node_memory_MemFree)%20-%20sum(node_memory_Buffers)%20-%20sum(node_memory_Cached))%20%2F%20sum(node_memory_MemTotal))%20*%20100"
	alertapi string = "/api/v1/query?query=ALERTS"
)

func ListSingleClusterCpu(address string) []byte {
	ProRequest.Host = address
	ProRequest.Path = cpuapi
	data, err := ProRequest.Get()
	if err != nil {
		return nil
	}
	return data
}

func ListSingleClusterMem(address string) []byte {
	ProRequest.Host = address
	ProRequest.Path = memapi
	data, err := ProRequest.Get()
	if err != nil {
		return nil
	}
	return data
}

func ListSingleClusterAlerts(address string) []byte {
	ProRequest.Host = address
	ProRequest.Path = alertapi
	data, err := ProRequest.Get()
	if err != nil {
		return nil
	}
	return data
}

func GetProAddressFromEnv(K8sconfigs configs.AllK8SConfigs) map[string]string {
	ProCfg := make(map[string]string)
	for k := range K8sconfigs {
		tmp := ConvertClusterName(k)
		if tmp != "" {
			ProCfg[k] = tmp
		} else {
			log.Println("this cluster " + k + " doesn't have the prometheus config, we will try use the podip to access it")
			return nil
		}
	}
	ProUseEnv = true
	return ProCfg
}

func ConvertClusterName(cluster string) string {
	if strings.Contains(cluster, "-") {
		cluster = strings.Replace(cluster, "-", "_", -1)
	}

	cluster = os.Getenv(strings.ToUpper(cluster) + "PRO")
	return cluster
}
