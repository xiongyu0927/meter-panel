package k8s

import (
	"encoding/json"
	"log"
	"meter-panel/configs"
)

var (
	svcslist                     SingleClusterSvcList
	ingresslist                  SingleClusterIngressList
	NilHumanSingleClusterLbsList HumanSingleClusterLbsList
)

func ListSingleClusterLB(k8sconfig configs.HumanSingleK8sConfigs) HumanSingleClusterLbsList {
	var lbstatus []lb
	ListSingleClusterSvcs(k8sconfig, &lbstatus)
	ListSingleClusterIngress(k8sconfig, &lbstatus)
	var tmp = HumanSingleClusterLbsList{
		LoadBalancer: lbstatus,
	}
	return tmp
}

func ListSingleClusterSvcs(k8sconfig configs.HumanSingleK8sConfigs, lbstatus *[]lb) {
	K8sRequest.Host = k8sconfig.EndPoint
	K8sRequest.BearToken = k8sconfig.Token
	K8sRequest.Path = "/api/v1/services"
	data, err := K8sRequest.Get()
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(data, &svcslist)
	if err != nil {
		log.Println(err)
	}

	SvcDetail(svcslist.Items, lbstatus)

	for {
		if svcslist.Metadata.Continue != "" {
			K8sRequest.Path = "/api/v1/services?limit=500&continue=" + svcslist.Metadata.Continue
			data, err = K8sRequest.Get()
			if err != nil {
				log.Println(err)
			}

			err = json.Unmarshal(data, &svcslist)
			if err != nil {
				log.Println(err)
			}
			SvcDetail(svcslist.Items, lbstatus)
		} else {
			break
		}
	}
}

func SvcDetail(item []svc, lbstatus *[]lb) {
	for _, v2 := range item {
		if v2.Spec.Type == "LoadBalancer" {
			var tmpip []string
			for _, v3 := range v2.Status.LoadBalancer.Ingress {
				tmpip = append(tmpip, v3.Ip)
			}

			tmp := lb{
				Name: v2.Metadata.Name,
				Type: v2.Spec.Type,
				Ip:   tmpip,
				Port: v2.Spec.Ports,
			}
			*lbstatus = append(*lbstatus, tmp)
		}
	}
}

func ListSingleClusterIngress(k8sconfig configs.HumanSingleK8sConfigs, lbstatus *[]lb) {
	K8sRequest.Host = k8sconfig.EndPoint
	K8sRequest.BearToken = k8sconfig.Token
	K8sRequest.Path = "/apis/extensions/v1beta1/ingresses"
	data, err := K8sRequest.Get()
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(data, &ingresslist)
	if err != nil {
		log.Println(err)
	}

	IngressDetail(ingresslist.Items, lbstatus)

	for {
		if ingresslist.Metadata.Continue != "" {
			K8sRequest.Path = "/apis/extensions/v1beta1/ingresses?limit=500&continue=" + applist.Metadata.Continue
			data, err = K8sRequest.Get()
			if err != nil {
				log.Println(err)
			}
			err = json.Unmarshal(data, &ingresslist)
			if err != nil {
				log.Println(err)
			}
			IngressDetail(ingresslist.Items, lbstatus)
		} else {
			break
		}
	}
}

func IngressDetail(items []Realingress, lbstatus *[]lb) {
	for _, v2 := range items {
		var tmphost []string
		for _, v3 := range v2.Spec.Rules {
			tmphost = append(tmphost, v3.Host)
		}

		tmp := lb{
			Name: v2.Metadata.Name,
			Type: "ingress",
			Ip:   tmphost,
			Port: nil,
		}
		*lbstatus = append(*lbstatus, tmp)
	}
}
