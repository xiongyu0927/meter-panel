package crd

import (
	"log"
	"meter-panel/configs"
	"meter-panel/pkg/api/v1/k8s/crd/alaudacrd"
	"meter-panel/pkg/api/v1/k8s/crd/application"
	"meter-panel/pkg/api/v1/k8s/crd/auth"
	"meter-panel/pkg/api/v1/k8s/crd/cluster"
	"meter-panel/pkg/api/v1/k8s/crd/devops"

	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type Interface interface {
	App(namespace string) application.AppInterface
	Pipeline(namespace string) devops.PipelineInterface
	Cqp(namesapce string) devops.CodeQualityProjectInterface
	Project(namespace string) auth.ProjectInterface
	Cluster(namespace string) cluster.ClusterInterface
	Rule(namespace string) alaudacrd.RuleInterface
}

type CrdClient struct {
	restClient rest.Interface
}

func (c *CrdClient) App(namespace string) application.AppInterface {
	return &application.AppClient{
		RestClient: c.restClient,
		Ns:         namespace,
	}
}

func (c *CrdClient) Pipeline(namespace string) devops.PipelineInterface {
	return &devops.PipelineClient{
		RestClient: c.restClient,
		Ns:         namespace,
	}
}

func (c *CrdClient) Cqb(namespace string) devops.CodeQualityProjectInterface {
	return &devops.CodeQualityProjectClient{
		RestClient: c.restClient,
		Ns:         namespace,
	}
}

func (c *CrdClient) Project(namespace string) auth.ProjectInterface {
	return &auth.ProjectClient{
		RestClient: c.restClient,
		Ns:         namespace,
	}
}

func (c *CrdClient) Cluster(namespace string) cluster.ClusterInterface {
	return &cluster.ClusterClient{
		RestClient: c.restClient,
		Ns:         namespace,
	}
}

func (c *CrdClient) Rule(namespace string) alaudacrd.RuleInterface {
	return &alaudacrd.RuleClient{
		RestClient: c.restClient,
		Ns:         namespace,
	}
}

type AllStore struct {
	crdclientset  map[string]*CrdClient
	AppStore      map[string]cache.Store
	PipelineStore cache.Store
	CqbStore      cache.Store
	ProjectStore  cache.Store
	ClusterStore  cache.Store
	RuleStore     map[string]cache.Store
}

var (
	appGroup     = &schema.GroupVersion{Group: "app.k8s.io", Version: "v1beta1"}
	devopsGroup  = &schema.GroupVersion{Group: "devops.alauda.io", Version: "v1alpha1"}
	authGroup    = &schema.GroupVersion{Group: "auth.alauda.io", Version: "v1"}
	clusterGroup = &schema.GroupVersion{Group: "clusterregistry.k8s.io", Version: "v1alpha1"}
	ruleGroup    = &schema.GroupVersion{Group: "crd.alauda.io", Version: "v1"}
	group        = []*schema.GroupVersion{appGroup, devopsGroup, authGroup, clusterGroup, ruleGroup}
)

func NewAllStore(HAKC configs.AllK8SConfigs) *AllStore {
	AS := &AllStore{
		AppStore:  make(map[string]cache.Store),
		RuleStore: make(map[string]cache.Store),
	}
	alauda := viper.GetString("ALAUDA")
	for k, v := range HAKC {
		cs1, err := NewCrdClientForConfig(v, appGroup)
		if err != nil {
			log.Println(err)
		}
		application.AddToScheme(scheme.Scheme)
		AS.AppStore[k] = application.WatchResources(cs1.App(""))

		if alauda == "ACE" {
			continue
		}

		cs5, err := NewCrdClientForConfig(v, ruleGroup)
		if err != nil {
			log.Println(err)
		}
		alaudacrd.AddToScheme(scheme.Scheme)
		AS.RuleStore[k] = alaudacrd.RuleWatchResources(cs5.Rule(""))

		if k == configs.GlobalName {
			cs2, err := NewCrdClientForConfig(v, devopsGroup)
			if err != nil {
				log.Println(err)
			}
			devops.AddToScheme(scheme.Scheme)
			AS.PipelineStore = devops.PipeWatchResources(cs2.Pipeline(""))
			AS.CqbStore = devops.CqbWatchResources(cs2.Cqb(""))

			cs3, err := NewCrdClientForConfig(v, authGroup)
			if err != nil {
				log.Println(err)
			}
			auth.AddToScheme(scheme.Scheme)
			AS.ProjectStore = auth.ProjectWatchResources(cs3.Project(""))

			cs4, err := NewCrdClientForConfig(v, clusterGroup)
			if err != nil {
				log.Println(err)
			}
			cluster.AddToScheme(scheme.Scheme)
			AS.ClusterStore = cluster.ClusterWatchResources(cs4.Cluster(""))
		}
	}
	return AS
}

func NewCrdClientForConfig(c *rest.Config, s *schema.GroupVersion) (*CrdClient, error) {
	config := *c

	config.ContentConfig.GroupVersion = s
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &CrdClient{restClient: client}, nil
}
