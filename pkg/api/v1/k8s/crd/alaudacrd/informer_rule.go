package alaudacrd

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type RuleInterface interface {
	List(opts metav1.ListOptions) (*RuleList, error)
	Get(name string, options metav1.GetOptions) (*Rule, error)
	Create(rule *Rule) (*Rule, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type RuleClient struct {
	RestClient rest.Interface
	Ns         string
}

func (c *RuleClient) List(opts metav1.ListOptions) (*RuleList, error) {
	result := RuleList{}
	err := c.RestClient.
		Get().
		Namespace(c.Ns).
		Resource("rules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *RuleClient) Get(name string, opts metav1.GetOptions) (*Rule, error) {
	result := Rule{}
	err := c.RestClient.
		Get().
		Namespace(c.Ns).
		Resource("rules").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *RuleClient) Create(rule *Rule) (*Rule, error) {
	result := Rule{}
	err := c.RestClient.
		Post().
		Namespace(c.Ns).
		Resource("rules").
		Body(rule).
		Do().
		Into(&result)

	return &result, err
}

func (c *RuleClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Ns).
		Resource("rules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

func RuleWatchResources(clientSet RuleInterface) cache.Store {
	ruleStore, ruleController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.Watch(lo)
			},
		},
		&Rule{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)
	stop := make(chan struct{})
	go ruleController.Run(stop)
	return ruleStore
}
