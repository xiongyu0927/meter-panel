package crd

import (
	"log"
	"meter-panel/configs"
	"meter-panel/pkg/api/v1/k8s/crd/application"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
)

type EveryClusterApplicationStore map[string]cache.Store

type AllStore struct {
	ClientSet map[string]*application.ApplicationClient
	AppStore  EveryClusterApplicationStore
}

func NewAllStore(HAKC configs.AllK8SConfigs) *AllStore {
	var cs = make(map[string]*application.ApplicationClient)
	var err error
	for k, v := range HAKC {
		cs[k], err = application.NewAppClientForConfig(v)
		if err != nil {
			log.Println(err)
		}
	}
	AS := &AllStore{
		ClientSet: cs,
		AppStore:  make(map[string]cache.Store),
	}
	AS.ResuorceLoad()
	return AS
}

func (its *AllStore) ResuorceLoad() {
	application.AddToScheme(scheme.Scheme)
	for k, v := range its.ClientSet {
		its.AppStore[k] = application.WatchResources(v)
	}
}
