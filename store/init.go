package store

import (
	"fmt"
	"log"
	"meter-panel/configs"
	"meter-panel/pkg/api/v1/k8s/core"
	"meter-panel/pkg/api/v1/k8s/crd"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var (
	// StoreAllK8SConfigs is used save all cluster configs in the memory
	StoreAllK8SConfigs configs.AllK8SConfigs
	AllLister          *core.AllLister
	AllStore           *crd.AllStore
	ProUseEnv          bool
	ProCfg             map[string]string
	EsClient           *core.EsClient
	err                error
)

func init() {
	// StoreAllK8SConfigs = configs.InitK8SCoinfg()
	StoreAllK8SConfigs = configs.LocalTest()
	AllLister = core.NewAllLister(StoreAllK8SConfigs)
	AllStore = crd.NewAllStore(StoreAllK8SConfigs)
	ProCfg = GetProAddressFromEnv(StoreAllK8SConfigs)
	// EsClient, err = core.NewEsClient(StoreAllK8SConfigs)
	if err != nil {
		log.Println(err)
	}
	// EsClient.Loop()
	// log.Println(EsClient.Data)
	log.Println("all resource of k8s, prometheus and es were init successed")
}

func AddNewClusterResource(cluster string) {
	cf := StoreAllK8SConfigs.GetSingleK8SConig(cluster)
	cs := kubernetes.NewForConfigOrDie(cf)
	stopper := make(chan struct{})
	factory := informers.NewSharedInformerFactory(cs, 30)
	ifs := AllLister.RegisterInformorAndLister(factory, cluster)
	defer runtime.HandleCrash()
	go factory.Start(stopper)

	for _, v := range ifs {
		if !cache.WaitForCacheSync(stopper, v.HasSynced) {
			runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
			return
		}
	}
}
