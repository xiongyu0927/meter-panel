package constome

import (
	"fmt"
	"log"
	"meter-panel/pkg/api/v1/k8s/crd/application"
	"meter-panel/store"

	av1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	label      string = "app.alauda.io/name"
	label2     string = "service_name"
	Running    string = "Running"
	Stop       string = "Stop"
	Processing string = "Processing"
)

type CebUnhealthyPodslis struct {
	Name      string
	NameSpace string
}

// CebClusterPodsList is ceb style pod list
type CebClusterPodsList struct {
	UnHealthyPodList []CebUnhealthyPodslis
	UnHealthyNumber  int
	HealthyNumber    int
	Total            int
}

// CebClusterNodesList is ceb style node list
type CebClusterNodesList struct {
	UnHealthyNodeList []v1.Node
	ResourcePressure  []map[string][]string
	UnHealthyNumber   int
	HealthyNumber     int
	Total             int
}

// CebClusterPvList is a string don't have type

// CebClusterSvcList is a ceb style service list
type CebClusterSvcList struct {
	LbName string
	Type   string
	Access string
}

//
type CebClusterIngressList struct {
	LbName string
	Type   string
	Access []v1beta1.IngressRule
}

type CebLB struct {
	L4 []CebClusterSvcList
	L7 []CebClusterIngressList
}

// CebDeploymentList is those unhealthy app's deployment
type CebDeploymentList struct {
	Name      string
	NameSpace string
}

// CebDaemonSetList is those unhealthy app's ds
type CebDaemonSetList struct {
	Name      string
	NameSpace string
}

type CebStatefulSetList struct {
	Name      string
	NameSpace string
}

// CebStyle achieve OrganizeData interface
type CebStyle struct {
	// PodList     CebClusterPodsList
	// NodeList    CebClusterNodesList
	// SvcList     CebClusterSvcList
	// IngressList CebClusterIngressList
	// AppDeployTotal  int
	// AppDeployStop   int
	// AppDsTotal      int
	// AppStaefulTotal int
	// AppStaefulStop  int
}

func NewCebStyle() *CebStyle {
	Style := &CebStyle{}
	return Style
}

// OrganizePodsList is a functiona under the OrganizeData interface that used to return your style of pod list
func (its *CebStyle) OrganizePodsList(cluster string, Pl []*v1.Pod) interface{} {
	var tmp []CebUnhealthyPodslis
	for _, v := range Pl {
		if v.Status.Phase != "Running" {
			if v.Status.Phase == "Succeeded" {
				continue
			}
			a := CebUnhealthyPodslis{
				Name:      v.Name,
				NameSpace: v.Namespace,
			}
			tmp = append(tmp, a)
			continue
		}
		if its.IsReady(v.Status.Conditions) {
			a := CebUnhealthyPodslis{
				Name:      v.Name,
				NameSpace: v.Namespace,
			}
			tmp = append(tmp, a)
		}
	}
	total := len(Pl)
	uhn := len(tmp)
	hn := total - uhn
	list := CebClusterPodsList{
		UnHealthyPodList: tmp,
		UnHealthyNumber:  uhn,
		HealthyNumber:    hn,
		Total:            total,
	}
	return list
}

// OrganizeNodesList is a functiona under the OrganizeData interface that used to return your style of node list
func (its *CebStyle) OrganizeNodesList(Cluster string, Nl []*v1.Node) interface{} {
	var tmp []v1.Node
	var tmp2 []map[string][]string
	for _, v := range Nl {
		if its.IsReady(v.Status.Conditions) {
			tmp = append(tmp, *v)
		}
		a := its.NodeIsEasy(v.Name, v.Status.Conditions)
		if a != nil {
			tmp2 = append(tmp2, a)
		}
	}
	total := len(Nl)
	uhn := len(tmp)
	hn := total - uhn
	list := CebClusterNodesList{
		UnHealthyNodeList: tmp,
		ResourcePressure:  tmp2,
		UnHealthyNumber:   uhn,
		HealthyNumber:     hn,
		Total:             total,
	}
	return list
}

// OrganizePvList a functiona under the OrganizeData interface that used to return your style of pv list
func (its *CebStyle) OrganizePvList(cluster string, Pvl []*v1.PersistentVolume) interface{} {
	var s int64
	var suffix string
	for _, v := range Pvl {
		quantity := new(resource.Quantity)
		*quantity = v.Spec.Capacity["storage"]
		size := quantity.Value()
		s = s + size
	}
	b := s
	if s > (1 << 30) {
		suffix = "Gi"
		b = s / (1 << 30)
	} else if s > (1 << 20) {
		suffix = "Mi"
		b = s / (1 << 20)
	} else if s > (1 << 10) {
		suffix = "Ki"
		b = s / (1 << 10)
	}
	return fmt.Sprintf("%v%v", b, suffix)
}

// OrganizeSvcList is .....
func (its *CebStyle) OrganizeSvcList(cluster string, sl []*v1.Service) interface{} {
	var tmp []CebClusterSvcList
	var access string
	for _, v := range sl {
		if v.Spec.Type == "LoadBalancer" {
			if v.Status.LoadBalancer.Ingress == nil {
				access = "Pending"
			} else {
				access = v.Status.LoadBalancer.Ingress[0].IP
			}
			var x = CebClusterSvcList{
				LbName: v.ObjectMeta.Name,
				Type:   "LoadBalancer",
				Access: access,
			}
			tmp = append(tmp, x)
		}
	}
	return tmp
}

func (its *CebStyle) OrganizeIngressList(cluster string, il []*v1beta1.Ingress) interface{} {
	var tmp []CebClusterIngressList
	for _, v := range il {
		var x = CebClusterIngressList{
			LbName: v.ObjectMeta.Name,
			Type:   "Ingress",
			Access: v.Spec.Rules,
		}
		tmp = append(tmp, x)
	}
	return tmp
}

// 0 is Stop, 1 is Processing, 2 is Running
func (its *CebStyle) OrganzieDeploymentList(dl []*av1.Deployment) []int {
	status := make([]int, 0)
	for _, v := range dl {
		if *v.Spec.Replicas == 0 {
			status = append(status, 0)
			continue
		}

		if v.Status.UnavailableReplicas != 0 {
			status = append(status, 1)
			continue
		}

		status = append(status, 2)
	}
	return status
}

// 0 is Stop, 1 is Processing, 2 is Running
func (its *CebStyle) OrganzieDaemonSetList(dsl []*av1.DaemonSet) []int {
	status := make([]int, 0)
	for _, v := range dsl {
		if v.Status.DesiredNumberScheduled != v.Status.NumberAvailable {
			status = append(status, 1)
		} else {
			status = append(status, 2)
		}
	}
	return status
}

// 0 is Stop, 1 is Processing, 2 is Running
func (its *CebStyle) OrganzieStatefulSetList(sfl []*av1.StatefulSet) []int {
	status := make([]int, 0)
	for _, v := range sfl {
		if *v.Spec.Replicas == 0 {
			status = append(status, 0)
			continue
		}
		if *v.Spec.Replicas != v.Status.ReadyReplicas || *v.Spec.Replicas != v.Status.CurrentReplicas {
			status = append(status, 1)
			continue
		}
		status = append(status, 2)
	}
	return status
}

// this is not OrgnaziData interface's function bool is
func (its *CebStyle) OrganzieApplicationList(cluster string, al []interface{}) interface{} {
	var papp []CebProcessingApp
	var Tnum, Rnum, Pnum, Snum int

	for _, v := range al {
		t, ok := v.(*application.Application)
		if !ok {
			log.Println("wrong type")
			break
		}
		Tnum++
		target := its.judgeSingleAppStatus(cluster, t)
		switch target {
		case Running:
			Rnum++
		case Processing:
			Pnum++
			tmp := CebProcessingApp{
				Name:      t.Name,
				NameSpace: t.Namespace,
			}
			papp = append(papp, tmp)
		case Stop:
			Snum++
		}
	}

	var tmp4 = CebApp{
		TotalNum:      Tnum,
		RunningNum:    Rnum,
		ProcessingNum: Pnum,
		StopNum:       Snum,
		PApp:          papp,
	}
	return tmp4
}

type singleapp struct {
	Name      string
	NameSpace string
	Status    string
}

func (its *CebStyle) judgeSingleAppStatus(cluster string, t *application.Application) string {
	var status []int
	labelset := labels.Set(t.Spec.Selector.MatchLabels).AsSelector()
	kind := t.Spec.ComponentGroupKinds
	for _, v := range kind {
		switch v.Kind {
		case "Deployment":
			lister := store.AllLister.DeploymentLister[cluster]
			dl, err := lister.List(labelset)
			if err != nil {
				log.Println(err)
			}
			tmp := its.OrganzieDeploymentList(dl)
			status = append(status, tmp...)
		case "DaemonSet":
			lister := store.AllLister.DaemonSetLister[cluster]
			dsl, err := lister.List(labelset)
			if err != nil {
				log.Println(err)
			}
			tmp := its.OrganzieDaemonSetList(dsl)
			status = append(status, tmp...)
		case "StatefulSet":
			lister := store.AllLister.StatefulSetLister[cluster]
			sl, err := lister.List(labelset)
			if err != nil {
				log.Println(err)
			}
			tmp := its.OrganzieStatefulSetList(sl)
			status = append(status, tmp...)
		default:
		}
	}

	target := targetStatus(status)
	return target
}

func targetStatus(status []int) string {
	n := len(status)
	score := 0
	for _, v := range status {
		score = score + v
	}

	if score == 2*n {
		return Running
	}

	if score == 0 {
		return Stop
	}

	return Processing
}

// IsReady is used to judge pods that is abled to recive request
func (its *CebStyle) IsReady(con interface{}) bool {
	switch t := con.(type) {
	case []v1.PodCondition:
		for _, v := range t {
			if v.Type == "Ready" && v.Status != "True" {
				return true
			}
		}
	case []v1.NodeCondition:
		for _, v := range t {
			if v.Type == "Ready" && v.Status != "True" {
				return true
			}
		}
	}
	return false
}

func (its *CebStyle) NodeIsEasy(NodeName string, con interface{}) map[string][]string {
	tmp := make(map[string][]string)
	t, ok := con.([]v1.NodeCondition)
	if ok {
		for _, v := range t {
			if v.Type == "MemoryPressure" && v.Status == "True" {
				x := "MemoryPressure"
				tmp[NodeName] = append(tmp[NodeName], x)
			}
			if v.Type == "DiskPressure" && v.Status == "True" {
				x := "DiskPressure"
				tmp[NodeName] = append(tmp[NodeName], x)
			}
			if v.Type == "PIDPressure" && v.Status == "True" {
				x := "PIDPressure"
				tmp[NodeName] = append(tmp[NodeName], x)
			}
		}
		if len(tmp[NodeName]) == 0 {
			return nil
		}
	}
	return tmp
}

type CebApp struct {
	TotalNum      int
	RunningNum    int
	StopNum       int
	ProcessingNum int
	PApp          []CebProcessingApp
}

type CebProcessingApp struct {
	Name      string
	NameSpace string
}
