package cluster

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type ClusterInterface interface {
	List(opts metav1.ListOptions) (*ClusterList, error)
	Get(name string, options metav1.GetOptions) (*Cluster, error)
	Create(proj *Cluster) (*Cluster, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type ClusterClient struct {
	RestClient rest.Interface
	Ns         string
}

func (c *ClusterClient) List(opts metav1.ListOptions) (*ClusterList, error) {
	result := ClusterList{}
	err := c.RestClient.
		Get().
		Namespace(c.Ns).
		Resource("clusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)
	// fmt.Println(result)
	return &result, err
}

func (c *ClusterClient) Get(name string, opts metav1.GetOptions) (*Cluster, error) {
	result := Cluster{}
	err := c.RestClient.
		Get().
		Namespace(c.Ns).
		Resource("clusters").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *ClusterClient) Create(proj *Cluster) (*Cluster, error) {
	result := Cluster{}
	err := c.RestClient.
		Post().
		Namespace(c.Ns).
		Resource("clusters").
		Body(proj).
		Do().
		Into(&result)

	return &result, err
}

func (c *ClusterClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Ns).
		Resource("clusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

func ClusterWatchResources(clientSet ClusterInterface) cache.Store {
	clusterStore, clusterController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.Watch(lo)
			},
		},
		&Cluster{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)
	stop := make(chan struct{})
	go clusterController.Run(stop)
	return clusterStore
}
