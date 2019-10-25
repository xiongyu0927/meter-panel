package configs

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"meter-panel/tools"
	"time"

	cav1 "acp/cluster-registry/api/v1alpha1"

	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

type AllK8SConfigs map[string]*rest.Config

const tokenFile string = "/var/run/secrets/kubernetes.io/serviceaccount/token"

var furionrequest = tools.Request{
	Methoud: "GET",
	Host:    "",
	Path:    "/v1/regions",
	//IsHttps shuold be https or http
	IsHTTPS:   "http",
	BearToken: "",
}

var kubernetesRequest = tools.Request{
	Methoud:   "GET",
	Host:      "10.96.0.1:443",
	Path:      "/apis/clusterregistry.k8s.io/v1alpha1/clusters?limit=500",
	IsHTTPS:   "https",
	BearToken: "",
}

var GlobalName string

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
	platform := viper.GetString("ALAUDA")
	GlobalName = viper.GetString("GLOBAL_CLUSTER_NAME")

	if platform == "ACE" {
		furionrequest.Host = viper.GetString("FurionHost")
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

	// ACP init k8s config
	if platform == "ACP" {
		if token, err := ioutil.ReadFile(tokenFile); err == nil {
			kubernetesRequest.BearToken = string(token)
		}
		data, err := kubernetesRequest.Get()
		if err != nil {
			log.Println(err)
		}

		var cfg cav1.ClusterList
		err = json.Unmarshal(data, &cfg)
		if err != nil {
			log.Println(err)
		}

		var tmp = make(map[string]*rest.Config)
		var c attr
		for _, v := range cfg.Items {
			if a := v.Annotations["legacy.cluster.alauda.io/attr"]; a != "" {
				json.Unmarshal([]byte(a), &c)
				endpoint := c.Kubernetes.Endpoint
				token := c.Kubernetes.Token
				tmp[v.ObjectMeta.Name] = GenerateRestConfig(endpoint, token)
			} else {
				endpoint := "https://10.96.0.1:443"
				token := kubernetesRequest.BearToken
				tmp[v.ObjectMeta.Name] = GenerateRestConfig(endpoint, token)
			}
		}
		return tmp
	}

	log.Fatal("please check env \"ALAUDA\" it must be ACP or ACE")
	return nil
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
		Host:            "https://94.191.112.216:6443",
		BearerToken:     "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJhbGF1ZGEtc3lzdGVtIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4tcnE5NGMiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjRjYjUzYmU1LWNmYTItMTFlOS04ZjZkLTUyNTQwMDEzMTJjZiIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDphbGF1ZGEtc3lzdGVtOmRlZmF1bHQifQ.ILtyR5NJrHcFk0P7pN4xymup236MtSCQfH9fwdvpdC75B-1eXR4Z_TsiXM6J5pXNM8I-YoYGq8G1Cx0n8_HGifxBm5BVo-vV5mOL6jzGX6lZ32XO5y7yJ571mEW4ibpUN3rJiGvNPVmXyEtRyVe0NICHmG7waHhvR9HxqB-YQ5UROohuIX5wOPaRwDMiPNRhJEVIE72P4HE0FsC0Ab9vzrJLxuUHf1UznxHlGWD64LfhfUDEBxVGgiQYDcs8D_F9ImGHUXkDhSxK-f8ovSluhHg1MpZOlWumd_ObJ4mzLAae1cgOfsiCbTJX9Nj1IyLa4X0sRuYTVb_bITCTXyDiLQ",
		Timeout:         time.Duration(10) * time.Second,
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
	}
	var tmp = make(map[string]*rest.Config)
	tmp["ace"] = cf
	return tmp
}
