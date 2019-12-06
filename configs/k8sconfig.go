package configs

import (
	"encoding/json"
	"log"
	"meter-panel/tools"
	"time"

	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

type AllK8SConfigs map[string]*rest.Config

var furionrequest = tools.Request{
	Methoud: "GET",
	Host:    "furion:8080",
	Path:    "/v1/regions",
	//IsHttps shuold be https or http
	IsHTTPS:   "http",
	BearToken: "",
}

// GetSingleK8SConig is used return a k8s clien-go rest config
func (its *AllK8SConfigs) GetSingleK8SConig(cluster string) *rest.Config {
	if (*its)[cluster] != nil {
		return (*its)[cluster]
	}

	its.Update(cluster)
	log.Println("this cluster" + cluster + "didn't exit in the cache maybe a new cluster")
	return (*its)[cluster]
}

// Update is used update cluster config when has some new cluster added
func (its *AllK8SConfigs) Update(cluster string) {
	tmp := InitK8SCoinfg()
	(*its)[cluster] = tmp[cluster]
}

// InitK8SCoinfg is used get a all k8s cluster config
func InitK8SCoinfg() AllK8SConfigs {
	furionhost := viper.GetString("FURION_HOST")
	if furionhost == "" {
		log.Println("enviroment doesn't set, will use default config furion:8080")
	} else {
		furionrequest.Host = furionhost
	}
	data, err := furionrequest.Get()
	if err != nil {
		log.Fatal(err)
	}
	var k8sconfig []K8sconfigs
	err = json.Unmarshal(data, &k8sconfig)
	if err != nil {
		log.Println(err)
	}

	var tmp = make(map[string]*rest.Config)
	for _, v := range k8sconfig {
		endpoint := v.Attr.Kubernetes.Endpoint
		Token := v.Attr.Kubernetes.Token
		tmp[v.Name] = GenerateRestConfig(endpoint, Token)
	}

	return tmp
}

// GenerateRestConfig return a client-go rest config
func GenerateRestConfig(ep, tk string) *rest.Config {
	to := viper.GetInt("K8S_TIMEOUT")
	cf := &rest.Config{
		Host:            ep,
		BearerToken:     tk,
		Timeout:         time.Duration(to) * time.Second,
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
	}

	if cf.APIPath == "" {
		cf.APIPath = "/api"
	}
	if cf.GroupVersion == nil {
		cf.GroupVersion = &schema.GroupVersion{}
	}
	if cf.NegotiatedSerializer == nil {
		cf.NegotiatedSerializer = serializer.NewCodecFactory(runtime.NewScheme())
	}
	return cf
}

// LocalTest is local test
func LocalTest() AllK8SConfigs {
	cf := &rest.Config{
		Host:            "https://129.28.147.60:6443",
		BearerToken:     "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJjbHVzdGVycm9sZS1hZ2dyZWdhdGlvbi1jb250cm9sbGVyLXRva2VuLXA0ZndmIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImNsdXN0ZXJyb2xlLWFnZ3JlZ2F0aW9uLWNvbnRyb2xsZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiIzNjljMTU1ZC03NTU0LTExZTktYjFiNS01MjU0MDA1YzAyODEiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1zeXN0ZW06Y2x1c3RlcnJvbGUtYWdncmVnYXRpb24tY29udHJvbGxlciJ9.Egx8MTAAwy0n9FGGDeI-36y7KUCsFcjQWI4cLfJbh3AlIhIswsZpbcf12bNUVUhD-3NsQuN93F-bwgDCr3Ft6-Je9t-ofnCfOIcEPDA-8xPJhbpNZcWYqJGSWZRlVgIEXMbbJNbZGJmgrJkBN6tkc30N0B2RjNE3j4qY-Q-vf1gLPBDk2H8FLGnFa9Iy1py0xvlnrkXyFuqzmjEEvgRmfIBCNn-4ImRFLuGlZsFSGqrlo-JzOWf5tKzfyxXfOKUfF4OOXAgnbPxeWqsxySyJN0rigaqQT_3kuxMmyEN3kAsVk5jgj9jEzyfPe4nX0PySV1N7AqqkrI2qQeSm5FopbQ",
		Timeout:         time.Duration(10) * time.Second,
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
	}
	var tmp = make(map[string]*rest.Config)
	tmp["ace"] = cf
	return tmp
}
