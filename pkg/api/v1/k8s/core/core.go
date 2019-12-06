package core

import (
	"fmt"
	"meter-panel/configs"
	"time"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	lappv1 "k8s.io/client-go/listers/apps/v1"
	lcorev1 "k8s.io/client-go/listers/core/v1"
	v1beta1 "k8s.io/client-go/listers/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"
)

// EveryClusterPodLister is all clueter's pod lister
type EveryClusterPodLister map[string]lcorev1.PodLister

// EveryClusterNodeLister is all cluster's node lister
type EveryClusterNodeLister map[string]lcorev1.NodeLister

// EveryClusterSvcLister is all cluster's service lister
type EveryClusterSvcLister map[string]lcorev1.ServiceLister

// EveryClusterEpLister is all cluster's endpoint lister
// type EveryClusterEpLister map[string]lcorev1.EndpointsLister

// EveryClusterIngressLister is all cluster's Ingress lister
type EveryClusterIngressLister map[string]v1beta1.IngressLister

// EveryClusterPvLister is all cluster's Pv lister
type EveryClusterPvLister map[string]lcorev1.PersistentVolumeLister

// EveryClusterDeploymentLister is all cluster's Deployment lister
type EveryClusterDeploymentLister map[string]lappv1.DeploymentLister

// EveryClusterDaemonSetLister is all cluster's DaemonSet lister
type EveryClusterDaemonSetLister map[string]lappv1.DaemonSetLister

// EveryClusterStatefulSetLister is all cluster's statefulset lister
type EveryClusterStatefulSetLister map[string]lappv1.StatefulSetLister

type EveryClusterInformers map[string][]cache.SharedIndexInformer

type EveryClusterStopper map[string]chan struct{}

// AllLister is used to read k8s resuorce from local cache
type AllLister struct {
	ClientSet  map[string]*kubernetes.Clientset
	PodLister  EveryClusterPodLister
	NodeLister EveryClusterNodeLister
	SvcLister  EveryClusterSvcLister
	// EpLister      EveryClusterEpLister
	IngressLister     EveryClusterIngressLister
	PvLister          EveryClusterPvLister
	DeploymentLister  EveryClusterDeploymentLister
	DaemonSetLister   EveryClusterDaemonSetLister
	StatefulSetLister EveryClusterStatefulSetLister
	TheInformers      EveryClusterInformers
	Stopper           EveryClusterStopper
}

// ResuorceLoad is padding AllLister.*Lister
func (its *AllLister) ResuorceLoad() {
	for k, v := range its.ClientSet {
		stopper := make(chan struct{})
		factory := informers.NewSharedInformerFactory(v, 5*time.Minute)
		ifs := its.RegisterInformorAndLister(factory, k)
		defer runtime.HandleCrash()
		go factory.Start(stopper)
		its.Stopper[k] = stopper

		for _, v := range ifs {
			if !cache.WaitForCacheSync(stopper, v.HasSynced) {
				runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
				return
			}
		}
	}
}

func (its *AllLister) RegisterInformorAndLister(factory informers.SharedInformerFactory, cluster string) []cache.SharedIndexInformer {
	var tmp []cache.SharedIndexInformer
	podInformer := factory.Core().V1().Pods()
	nodeInformer := factory.Core().V1().Nodes()
	serviceInformer := factory.Core().V1().Services()
	// epInformer := factory.Core().V1().Endpoints()
	inInformer := factory.Extensions().V1beta1().Ingresses()
	pvInformer := factory.Core().V1().PersistentVolumes()
	deploymentInformer := factory.Apps().V1().Deployments()
	daemonsetInformer := factory.Apps().V1().DaemonSets()
	statefulsetInforemer := factory.Apps().V1().StatefulSets()
	its.PodLister[cluster] = podInformer.Lister()
	its.NodeLister[cluster] = nodeInformer.Lister()
	its.SvcLister[cluster] = serviceInformer.Lister()
	// its.EpLister[cluster] = epInformer.Lister()
	its.IngressLister[cluster] = inInformer.Lister()
	its.PvLister[cluster] = pvInformer.Lister()
	its.DeploymentLister[cluster] = deploymentInformer.Lister()
	its.DaemonSetLister[cluster] = daemonsetInformer.Lister()
	its.StatefulSetLister[cluster] = statefulsetInforemer.Lister()
	pinformer := podInformer.Informer()
	ninformer := nodeInformer.Informer()
	sinformer := serviceInformer.Informer()
	// einformer := epInformer.Informer()
	iinformer := inInformer.Informer()
	pvinformer := pvInformer.Informer()
	dpinformer := deploymentInformer.Informer()
	dsinformer := daemonsetInformer.Informer()
	ssinformer := statefulsetInforemer.Informer()
	its.TheInformers[cluster] = append(its.TheInformers[cluster], dpinformer, ssinformer, sinformer)
	tmp = append(tmp, pinformer, ninformer, sinformer, iinformer, pvinformer, dpinformer, dsinformer, ssinformer)
	return tmp
}

// NewAllLister will return a Lister taht used to list all k8s resource from local cache
func NewAllLister(HAKC configs.AllK8SConfigs) *AllLister {
	var cs = make(map[string]*kubernetes.Clientset)
	for k, v := range HAKC {
		cs[k] = kubernetes.NewForConfigOrDie(v)
	}
	Al := &AllLister{
		ClientSet:  cs,
		PodLister:  make(map[string]lcorev1.PodLister),
		NodeLister: make(map[string]lcorev1.NodeLister),
		SvcLister:  make(map[string]lcorev1.ServiceLister),
		// EpLister:      make(map[string]lcorev1.EndpointsLister),
		IngressLister:     make(map[string]v1beta1.IngressLister),
		PvLister:          make(map[string]lcorev1.PersistentVolumeLister),
		DeploymentLister:  make(map[string]lappv1.DeploymentLister),
		DaemonSetLister:   make(map[string]lappv1.DaemonSetLister),
		StatefulSetLister: make(map[string]lappv1.StatefulSetLister),
		TheInformers:      make(map[string][]cache.SharedIndexInformer),
		Stopper:           make(map[string]chan struct{}),
	}
	Al.ResuorceLoad()
	return Al
}
