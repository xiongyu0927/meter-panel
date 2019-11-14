package store

import (
	"strings"
	"sync"

	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

const (
	key           string = "app.alauda.io/name"
	key2          string = "app.cpaas.io/name"
	appApiVersion string = "app.k8s.io/v1beta1"
	appKind       string = "Application"
	depGroup      string = "extensions"
	depKind       string = "Deployment"
	svcGroup      string = ""
	svcKind       string = "Service"
	sfGroup       string = "app"
	sfKind        string = "StatefulSet"
	appAnotation1 string = "app.alauda.io/create-method"
	appAnotation2 string = "app.alauda.io/display-name"
)

var ehf = cache.ResourceEventHandlerFuncs{
	AddFunc:    OnAdd,
	UpdateFunc: OnUpdate,
	DeleteFunc: OnDelete,
}

type Models struct {
	sync.RWMutex
	IsCreating map[string]bool
}

func NewModels(data map[string][]cache.SharedIndexInformer) *Models {
	tmp := &Models{
		IsCreating: make(map[string]bool),
	}
	ci := strings.SplitN(viper.GetString("MONITOR_CLUSTER"), ",", -1)
	if ci[0] == "" {
		for k, v := range data {
			go registryEventHandler(v, k)
		}
	} else {
		for _, v := range ci {
			go registryEventHandler(data[v], v)
		}
	}
	return tmp
}

func registryEventHandler(data []cache.SharedIndexInformer, cluster string) {
	for _, v := range data {
		v.AddEventHandler(ehf)
	}
	<-AllLister.Stopper[cluster]
}

type meta struct {
	name        string
	namespace   string
	labels      map[string]string
	or          []metav1.OwnerReference
	appname     string
	clustername string
}

func OnAdd(obj interface{}) {
	switch t := obj.(type) {
	case *appsv1.Deployment:
		Model.AddDeploymentFlow(t)
	case *appsv1.StatefulSet:
		Model.AddStatefulSetFlow(t)
	case *v1.Service:
		Model.AddServiceFlow(t)
	}
}

func OnUpdate(oobj interface{}, nobj interface{}) {
	switch new := nobj.(type) {
	case *appsv1.Deployment:
		old, _ := oobj.(*appsv1.Deployment)
		Model.UpdateDeploymentFlow(old, new)
	case *appsv1.StatefulSet:
		old, _ := oobj.(*appsv1.StatefulSet)
		Model.UpdateStatefulSetFlow(old, new)
	case *v1.Service:
		old, _ := oobj.(*v1.Service)
		Model.UpdateServiceFlow(old, new)
	}
}

func OnDelete(obj interface{}) {
	switch t := obj.(type) {
	case *appsv1.Deployment:
		Model.DeleteDeploymentFlow(t)
	case *appsv1.StatefulSet:
		Model.DeleteStatefulSetFlow(t)
	case *v1.Service:
		Model.DeleteServiceFlow(t)
	}
}
