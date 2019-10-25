package auth

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type ProjectInterface interface {
	List(opts metav1.ListOptions) (*ProjectList, error)
	Get(name string, options metav1.GetOptions) (*Project, error)
	Create(proj *Project) (*Project, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type ProjectClient struct {
	RestClient rest.Interface
	Ns         string
}

func (c *ProjectClient) List(opts metav1.ListOptions) (*ProjectList, error) {
	result := ProjectList{}
	err := c.RestClient.
		Get().
		Namespace(c.Ns).
		Resource("projects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)
	// fmt.Println(result)
	return &result, err
}

func (c *ProjectClient) Get(name string, opts metav1.GetOptions) (*Project, error) {
	result := Project{}
	err := c.RestClient.
		Get().
		Namespace(c.Ns).
		Resource("projects").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *ProjectClient) Create(proj *Project) (*Project, error) {
	result := Project{}
	err := c.RestClient.
		Post().
		Namespace(c.Ns).
		Resource("projects").
		Body(proj).
		Do().
		Into(&result)

	return &result, err
}

func (c *ProjectClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Ns).
		Resource("projects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

func ProjectWatchResources(clientSet ProjectInterface) cache.Store {
	pjStore, pjController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.Watch(lo)
			},
		},
		&Project{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)
	stop := make(chan struct{})
	go pjController.Run(stop)
	return pjStore
}
