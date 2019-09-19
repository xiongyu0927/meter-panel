package constome

import (
	"k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
)

// OrganizeData is a interface to save some Organize data's
type OrganizeData interface {
	OrganizePodsList(cluster string, pl []*v1.Pod) interface{}
	OrganizeNodesList(cluster string, nl []*v1.Node) interface{}
	OrganizePvList(cluster string, pvl []*v1.PersistentVolume) interface{}
	OrganizeSvcList(cluster string, sl []*v1.Service) interface{}
	// OrganizeEpList(cluster string, epl []*v1.Endpoints, sl []*v1.Service) interface{}
	OrganizeIngressList(cluster string, il []*v1beta1.Ingress) interface{}
	// OrganzieDeploymentList(cluster string, dl []*av1.Deployment) interface{}
	// OrganzieDaemonSetList(cluster string, dsl []*av1.DaemonSet) interface{}
	// OrganzieStatefulSetList(cluster string, sfl []*av1.StatefulSet) interface{}
	// OrganzieEventList()
	OrganzieApplicationList(cluster string, al []interface{}) interface{}
}

// Orgnazi call your OrgnaziData's function
// func Orgnazi(this OrganizeData, Kind interface{}, Cluster string) interface{} {
// 	switch t := Kind.(type) {
// 	case []*v1.Pod:
// 		return this.OrganizePodsList(Cluster, t)
// 	case []*v1.Node:
// 		return this.OrganizeNodesList(Cluster, t)
// 	case []*v1.PersistentVolume:
// 		return this.OrganizePvList(Cluster, t)
// 	case []*v1.Service:
// 		return this.OrganizeSvcList(Cluster, t)
// 	// case []*v1.Endpoints:
// 	// 	return this.OrganizeEpList(Cluster, t)
// 	case []*v1beta1.Ingress:
// 		return this.OrganizeIngressList(Cluster, t)
// 		// case []*av1.Deployment:
// 		// 	return this.OrganzieDeploymentList(Cluster, t)
// 		// case []*av1.DaemonSet:
// 		// 	return this.OrganzieDaemonSetList(Cluster, t)
// 		// case []*av1.StatefulSet:
// 		// 	return this.OrganzieStatefulSetList(Cluster, t)
// 	case []interface{}:
// 		return this.OrganzieApplicationList(Cluster, t)
// 	default:
// 		log.Println("Unknow Kind")
// 		return nil
// 	}
// }
