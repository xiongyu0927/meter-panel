package capacity

import (
	"encoding/json"
	"fmt"
	"log"
	"meter-panel/store"
	"meter-panel/tools"
	"strings"

	"github.com/robfig/cron"
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
	TotalClusterCount int
	TotalNodeCount    int
	TotalPodCount     int
	TotalSubnetCount  int
	// key 1 is cluster, key 2 is namespace and node
	Detail map[string]map[string]*resource
}

type resource struct {
	Node   int
	Pod    int
	Subnet int
}

func NewCebCacpityReport() *CebCacpityReport {
	tmp := &CebCacpityReport{
		Result: make(map[string]*result),
	}
	return tmp
}

func Start() {
	c := cron.New()
	spec := "0 0 18 * * *"
	c.AddFunc(spec, func() {
		start()
	})
	c.Start()
}

func start() {
	data := getData()
	transferData(*data)
}

func getData() *CebCacpityReport {
	report := NewCebCacpityReport()
	cns, err := getClusterAndAllNamespaces()
	if err != nil {
		log.Println(err)
		return nil
	}
	if err := report.getCacpity(cns); err != nil {
		log.Println(err)
		return nil
	}

	return report
}

func (r *CebCacpityReport) getCacpity(cns map[string][]string) error {
	for k, v := range cns {
		site := getSiteByClusterName(k)
		r.Result[site].TotalClusterCount++
		nc, err := getNodeCount(k)
		if err != nil {
			return err
		}
		r.Result[site].TotalNodeCount = r.Result[site].TotalNodeCount + nc

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

		r.Result[site].Detail[k]["Nodes"] = &resource{
			Node: nc,
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

func getNodeCount(cluster string) (int, error) {
	tmp, err := store.AllLister.NodeLister[cluster].List(labels.Everything())
	if err != nil {
		return -1, err
	}

	return len(tmp), nil
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
