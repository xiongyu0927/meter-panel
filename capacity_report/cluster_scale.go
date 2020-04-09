package capacity

import (
	"encoding/json"
	"fmt"
	"log"
	"meter-panel/store"
	"meter-panel/tools"

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

// clusterScalData is your report's data
type clusterScalData struct {
	TotalClusterCount int `json:"totalClusterCount"`
	TotalNodeCount    int `json:"totalNodeCount"`
	TotalPodCount     int `json:"totalPodCount"`
	TotalSubnetCount  int `json:"totalSubnetCount"`
}

func clusterScalReportStart() {
	data := GetClusterScalData()
	transferData(*data)
}

func GetClusterScalData() *CebCapacityReport {
	report := NewCebCapacityReport("capacityData")
	cns, err := getClusterAndAllNamespaces()
	if err != nil {
		log.Println(err)
		return nil
	}
	if err := report.getClusterScalCapacity(cns); err != nil {
		log.Println(err)
		return nil
	}

	return report
}

func (r *CebCapacityReport) getClusterScalCapacity(cns map[string][]string) error {
	var data clusterScalData
	for k, v := range cns {
		if r.isYourSiteByClusterName(k) {
			data.TotalClusterCount++

			nc, err := getNodeCount(k)
			if err != nil {
				return err
			}
			data.TotalNodeCount = data.TotalNodeCount + nc

			for _, v1 := range v {
				pc, err := getPodCount(k, v1)
				if err != nil {
					return err
				}
				data.TotalPodCount = data.TotalPodCount + pc

				sc, err := getSubnetCount(k, v1)
				if err != nil {
					return err
				}
				data.TotalSubnetCount = data.TotalSubnetCount + sc
			}
		}

	}
	r.Data = data
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

func getNodeCount(cluster string) (int, error) {
	tmp, err := store.AllLister.NodeLister[cluster].List(labels.Everything())
	if err != nil {
		return -1, err
	}

	return len(tmp), nil
}
