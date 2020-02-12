package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"meter-panel/store"
	"meter-panel/tools"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/labels"
)

var NsxtRequest = tools.Request{
	Methoud: "GET",
	Host:    viper.GetString("NSXT_GO_ENDPOINT"),
	Path:    "",
	//IsHttps shuold be https or http
	IsHTTPS:   "http",
	BearToken: "",
}

// CebSubnetData is used to parse response of nsxt-go
type CebSubnetData struct {
	Subnets []NewSubnet2
}

type NewSubnet2 struct {
	SubnetName string
	Cidr       string
	Id         string
}

type CebCacpityReport struct {
	// key is site
	Result map[string]*result
}

type result struct {
	TotalPodCount    int
	TotalSubnetCount int
	// key 1 is cluster, key 2 is namespace
	Detail map[string]map[string]*resource
}

type resource struct {
	Pod    int
	Subnet int
}

func NewCebCacpityReport() *CebCacpityReport {
	tmp := &CebCacpityReport{
		Result: make(map[string]*result),
	}
	return tmp
}

func CapacityViews(w http.ResponseWriter, r *http.Request) {
	report := NewCebCacpityReport()
	cns, err := getClusterAndAllNamespaces()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	if err := report.getCacpity(cns); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	data2, err := json.Marshal(report)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(200)
	w.Write(data2)
}

func (r *CebCacpityReport) getCacpity(cns map[string][]string) error {
	for k, v := range cns {
		site := getSiteByClusterName(k)
		if _, ok := r.Result[site]; !ok {
			r.Result[site] = &result{
				Detail: make(map[string]map[string]*resource),
			}
		}

		for _, v1 := range v {
			pc, err := getPodCount(k, v1)
			r.Result[site].TotalPodCount = r.Result[site].TotalPodCount + pc
			if err != nil {
				return err
			}

			sc, err := getSubnetCount(k, v1)
			r.Result[site].TotalSubnetCount = r.Result[site].TotalSubnetCount + sc
			if err != nil {
				return err
			}

			var rs = &resource{
				Pod:    pc,
				Subnet: sc,
			}
			if _, ok := r.Result[site].Detail[k]; !ok {
				r.Result[site].Detail[k] = make(map[string]*resource)
			}
			r.Result[site].Detail[k][v1] = rs
		}
	}
	return nil
}

func getClusterAndAllNamespaces() (map[string][]string, error) {
	tmp := make(map[string][]string)
	for k, v := range store.AllLister.NameSpaceLister {
		nss, err := v.List(labels.Everything())
		if err != nil {
			return nil, err
		}

		var namespaces []string
		for _, v := range nss {
			namespaces = append(namespaces, v.ObjectMeta.Name)
		}
		tmp[k] = namespaces
	}
	return tmp, nil
}

func getPodCount(cluster, namespace string) (int, error) {
	tmp, err := store.AllLister.PodLister[cluster].Pods(namespace).List(labels.Everything())
	if err != nil {
		return -1, err
	}
	pc := len(tmp)
	return pc, nil
}

func getSubnetCount(cluster, namespace string) (int, error) {
	var subnet CebSubnetData
	path := fmt.Sprintf("/nsxt/namespace/subnetview?cluster=%v&namespace=%v", cluster, namespace)
	NsxtRequest.Path = path
	data, err := NsxtRequest.Get()
	if err != nil {
		return -1, err
	}

	if err := json.Unmarshal(data, &subnet); err != nil {
		return -1, err
	}
	sc := len(subnet.Subnets)
	return sc, nil
}

func getSiteByClusterName(cluster string) string {
	slice := strings.SplitN(cluster, "-", -1)
	for _, v := range slice {
		if v == "as" || v == "ap" {
			return "A"
		}

		if v == "bs" || v == "bp" {
			return "B"
		}
	}

	return "S"
}
