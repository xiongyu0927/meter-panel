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
		Host:            "https://94.191.86.161:6443",
		BearerToken:     "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4tamJsdDciLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjI4N2M4Yjg1LTc3YjItMTFlOS05YmU1LTUyNTQwMDU3NTRhOCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.e2r2hPC1P8NyAN28VhpthYCmZ1sOwGBLc-8j-szr0ISmfEWudEQ-kTuP0w4BmfS9DpE_6GjERzOzKgbncNpQc3CrFgqHokLXzB-YbZnrxWZG4DdEsrio6Jvx_yeLTzTnAEzKRKDpcRqPCpx0kQ4GtzoJWKoptBXrLuuYHi6fJ7KY4DlfYFxHB_zh2mW9OQKfvfBJgYmHzSSrlIpf6GZdtiMjKpC1R9OYic9YHaUGiLf7q-TiWyTtudN6m-U--bDzbyzFeoXyMdbvJX0sbDNE16nRtSP0feU9vsKqwjbENCmDm1IX3c5JeRjDEMgecAyCTw6RCaXL2Y0yYVlLvCSQeg",
		Timeout:         time.Duration(10) * time.Second,
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
	}
	var tmp = make(map[string]*rest.Config)
	tmp["ace"] = cf
	return tmp
}
