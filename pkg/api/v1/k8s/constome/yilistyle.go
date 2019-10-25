package constome

import (
	"encoding/json"
	"fmt"
	"log"
	"meter-panel/pkg/api/v1/k8s/crd/application"
	"meter-panel/pkg/api/v1/k8s/crd/auth"
	"meter-panel/pkg/api/v1/k8s/crd/devops"
	"meter-panel/store"
	"meter-panel/tools"
	"strings"

	av1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"
)

type YiLiUnhealthyPodslis struct {
	Name      string
	NameSpace string
}

// YiLiClusterPodsList is YiLi style pod list
type YiLiClusterPodsList struct {
	UnHealthyPodList []YiLiUnhealthyPodslis
	UnHealthyNumber  int
	HealthyNumber    int
	Total            int
}

// YiLiClusterNodesList is YiLi style node list
type YiLiClusterNodesList struct {
	UnHealthyNodeList      []v1.Node
	ResourcePressure       []map[string][]string
	ResourcePressureNumber int
	UnHealthyNumber        int
	HealthyNumber          int
	Total                  int
}

// YiLiClusterPvList is a string don't have type

// YiLiClusterSvcList is a YiLi style service list
type YiLiClusterSvcList struct {
	LbName string
	Type   string
	Access string
}

//
type YiLiClusterIngressList struct {
	LbName string
	Type   string
	Access []v1beta1.IngressRule
}

type YiLiLB struct {
	L4 []YiLiClusterSvcList
	L7 []YiLiClusterIngressList
}

// YiLiDeploymentList is those unhealthy app's deployment
type YiLiDeploymentList struct {
	Name      string
	NameSpace string
}

// YiLiDaemonSetList is those unhealthy app's ds
type YiLiDaemonSetList struct {
	Name      string
	NameSpace string
}

type YiLiStatefulSetList struct {
	Name      string
	NameSpace string
}

// YiLiStyle achieve OrganizeData interface
type YiLiStyle struct {
	// PodList     YiLiClusterPodsList
	// NodeList    YiLiClusterNodesList
	// SvcList     YiLiClusterSvcList
	// IngressList YiLiClusterIngressList
	// AppDeployTotal  int
	// AppDeployStop   int
	// AppDsTotal      int
	// AppStaefulTotal int
	// AppStaefulStop  int
}

func NewYiLiStyle() *YiLiStyle {
	Style := &YiLiStyle{}
	return Style
}

// OrganizePodsList is a functiona under the OrganizeData interface that used to return your style of pod list
func (its *YiLiStyle) OrganizePodsList(cluster string, Pl []*v1.Pod) interface{} {
	var tmp []YiLiUnhealthyPodslis
	for _, v := range Pl {
		if v.Status.Phase != "Running" {
			if v.Status.Phase == "Succeeded" {
				continue
			}
			a := YiLiUnhealthyPodslis{
				Name:      v.Name,
				NameSpace: v.Namespace,
			}
			tmp = append(tmp, a)
			continue
		}
		if its.IsReady(v.Status.Conditions) {
			a := YiLiUnhealthyPodslis{
				Name:      v.Name,
				NameSpace: v.Namespace,
			}
			tmp = append(tmp, a)
		}
	}
	total := len(Pl)
	uhn := len(tmp)
	hn := total - uhn
	list := YiLiClusterPodsList{
		UnHealthyPodList: tmp,
		UnHealthyNumber:  uhn,
		HealthyNumber:    hn,
		Total:            total,
	}
	return list
}

// OrganizeNodesList is a functiona under the OrganizeData interface that used to return your style of node list
func (its *YiLiStyle) OrganizeNodesList(Cluster string, Nl []*v1.Node) interface{} {
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
	list := YiLiClusterNodesList{
		UnHealthyNodeList:      tmp,
		ResourcePressure:       tmp2,
		ResourcePressureNumber: len(tmp2),
		UnHealthyNumber:        uhn,
		HealthyNumber:          hn,
		Total:                  total,
	}
	return list
}

// OrganizePvList a functiona under the OrganizeData interface that used to return your style of pv list
func (its *YiLiStyle) OrganizePvList(cluster string, Pvl []*v1.PersistentVolume) interface{} {
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
func (its *YiLiStyle) OrganizeSvcList(cluster string, sl []*v1.Service) interface{} {
	var tmp []YiLiClusterSvcList
	var access string
	for _, v := range sl {
		if v.Spec.Type == "LoadBalancer" {
			if v.Status.LoadBalancer.Ingress == nil {
				access = "Pending"
			} else {
				access = v.Status.LoadBalancer.Ingress[0].IP
			}
			var x = YiLiClusterSvcList{
				LbName: v.ObjectMeta.Name,
				Type:   "LoadBalancer",
				Access: access,
			}
			tmp = append(tmp, x)
		}
	}
	return tmp
}

func (its *YiLiStyle) OrganizeIngressList(cluster string, il []*v1beta1.Ingress) interface{} {
	var tmp []YiLiClusterIngressList
	for _, v := range il {
		var x = YiLiClusterIngressList{
			LbName: v.ObjectMeta.Name,
			Type:   "Ingress",
			Access: v.Spec.Rules,
		}
		tmp = append(tmp, x)
	}
	return tmp
}

// 0 is Stop, 1 is Processing, 2 is Running
func (its *YiLiStyle) OrganizeDeploymentList(dl []*av1.Deployment) []int {
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
func (its *YiLiStyle) OrganizeDaemonSetList(dsl []*av1.DaemonSet) []int {
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
func (its *YiLiStyle) OrganizeStatefulSetList(sfl []*av1.StatefulSet) []int {
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
func (its *YiLiStyle) OrganizeApplicationList(cluster string, al []interface{}) interface{} {
	var papp []YiLiProcessingApp
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
			tmp := YiLiProcessingApp{
				Name:      t.Name,
				NameSpace: t.Namespace,
			}
			papp = append(papp, tmp)
		case Stop:
			Snum++
		}
	}

	var tmp4 = YiLiApp{
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

func (its *YiLiStyle) judgeSingleAppStatus(cluster string, t *application.Application) string {
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
			tmp := its.OrganizeDeploymentList(dl)
			status = append(status, tmp...)
		case "DaemonSet":
			lister := store.AllLister.DaemonSetLister[cluster]
			dsl, err := lister.List(labelset)
			if err != nil {
				log.Println(err)
			}
			tmp := its.OrganizeDaemonSetList(dsl)
			status = append(status, tmp...)
		case "StatefulSet":
			lister := store.AllLister.StatefulSetLister[cluster]
			sl, err := lister.List(labelset)
			if err != nil {
				log.Println(err)
			}
			tmp := its.OrganizeStatefulSetList(sl)
			status = append(status, tmp...)
		default:
		}
	}

	target := targetStatus(status)
	return target
}

// IsReady is used to judge pods that is abled to recive request
func (its *YiLiStyle) IsReady(con interface{}) bool {
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

func (its *YiLiStyle) NodeIsEasy(NodeName string, con interface{}) map[string][]string {
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

type YiLiApp struct {
	TotalNum      int
	RunningNum    int
	StopNum       int
	ProcessingNum int
	PApp          []YiLiProcessingApp
}

type YiLiProcessingApp struct {
	Name      string
	NameSpace string
}

type YiLiPipeline struct {
	Successed int
	Failed    int
}

func (its *YiLiStyle) OrganizePipelineList(project []string, pipel []interface{}) interface{} {
	var Rnum, Fnum, Pnum int
	for _, v := range pipel {
		t, ok := v.(*devops.Pipeline)
		if !ok {
			log.Println("wrong type")
			break
		}

		for _, v1 := range project {
			pj := strings.SplitN(v1, ".", -1)
			if t.Namespace != "" && t.Namespace == pj[len(pj)-1] {
				target := judgeSinglePipelineStatus(t)
				switch target {
				case "Successed":
					Rnum++
				case "Failed":
					Fnum++
				default:
					Pnum++
				}
			}
		}
	}

	tmp := YiLiPipeline{
		Successed: Rnum,
		Failed:    Fnum,
	}
	return tmp
}

func judgeSinglePipelineStatus(t *devops.Pipeline) string {
	switch t.Status.Phase {
	case devops.PipelinePhaseComplete:
		return "Successed"
	case devops.PipelinePhaseFailed, devops.PipelinePhaseError,
		devops.PipelinePhaseCancelled:
		return "Failed"
	default:
		return ""
	}
}

type Statistics struct {
	Succ    int `json:"succ"`
	Failed  int `json:"failed"`
	Warning int `json:"warn"`
	Total   int `json:"total"`
}

const path string = "/api/v1/statistics/codequality/"
const tokenFile string = "/var/run/secrets/kubernetes.io/serviceaccount/token"

var cqbreq = tools.Request{
	Methoud: "GET",
	Host:    "devops-api:8080",
	Path:    "",
	// IsHttps shuold be https or http
	IsHTTPS: "http",
}

func (its *YiLiStyle) OrganizeCqbList(project []string) interface{} {
	// tmp := viper.GetString("ACP_API_HOST")
	// address := strings.SplitN(tmp, "://", 2)
	// if token, err := ioutil.ReadFile(tokenFile); err == nil {
	// 	cqbreq.BearToken = string(token)
	// }
	// cqbreq.Host = address[1]
	// if address[0] == "https" {
	// 	cqbreq.IsHTTPS = "https"
	// } else {
	// 	cqbreq.IsHTTPS = "http"
	// }

	var result Statistics
	for _, v := range project {
		var tmp Statistics
		pj := strings.SplitN(v, ".", -1)
		opj := pj[len(pj)-1]
		cqbreq.Path = path + opj
		data, err := cqbreq.Get()
		if err != nil {
			log.Println(err)
			return nil
		}
		json.Unmarshal(data, &tmp)
		result.Failed = result.Failed + tmp.Failed
		result.Succ = result.Succ + tmp.Succ
		result.Warning = result.Warning + tmp.Warning
		result.Total = result.Total + tmp.Total
	}
	return result
	// var Rnum, Fnum, Wnum int
	// for _, v := range cqbl {
	// 	t, ok := v.(*devops.CodeQualityProject)
	// 	if !ok {
	// 		log.Println("wrong type")
	// 		break
	// 	}
	//
	// 	for _, v1 := range project {
	// 		pj := strings.SplitN(v1, ".", -1)
	// 		if t.Namespace != "" && t.Namespace == pj[len(pj)-1] {
	// 			target := judgeSingleCqbStatus(t)
	// 			switch target {
	// 			case "Successed":
	// 				Rnum++
	// 			case "Failed":
	// 				Fnum++
	// 			default:
	// 				Wnum++
	// 			}
	// 		}
	// 	}
	// }
	// tmp := YiLicqb{
	// 	Successed: Rnum,
	// 	Failed:    Fnum,
	// 	Warning:   Wnum,
	// }
	// return tmp
}

// func judgeSingleCqbStatus(t *devops.CodeQualityProject) string {
// switch t.Status.Phase {
// case devops.PipelinePhaseComplete:
// 	return "Successed"
// case devops.PipelinePhaseFailed, devops.PipelinePhaseError,
// 	devops.PipelinePhaseCancelled:
// 	return "Failed"
// default:
// 	return ""
// }
// 	return "waiting for data"
// }

type YiLiProject struct {
	Running    int
	Processing int
}

func (its *YiLiStyle) OrganizeProjectList(cluster string, prjl []interface{}) interface{} {
	var Rnum, Pnum int
	for _, v := range prjl {
		t, ok := v.(*auth.Project)
		if !ok {
			log.Println("wrong type")
			break
		}
		target := judgeSingleProjectStatus(cluster, t)
		switch target {
		case "Running":
			Rnum++
		case "Processing":
			Pnum++
		default:
		}
	}
	tmp := YiLiProject{
		Running:    Rnum,
		Processing: Pnum,
	}
	return tmp
}

func judgeSingleProjectStatus(cluster string, t *auth.Project) string {
	for _, v := range t.Spec.Clusters {
		if v.Name != "" && v.Name == cluster {
			switch t.Status.Phase {
			case "Active":
				return "Running"
			default:
				return "Processing"
			}
		}
	}
	return ""
}

func (its *YiLiStyle) OrganizeLbList(lb []interface{}) interface{} {
	return len(lb)
}
