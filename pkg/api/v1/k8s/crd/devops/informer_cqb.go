package devops

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type CodeQualityProjectInterface interface {
	List(opts metav1.ListOptions) (*CodeQualityProjectList, error)
	Get(name string, options metav1.GetOptions) (*CodeQualityProject, error)
	Create(pipe *Pipeline) (*CodeQualityProject, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type CodeQualityProjectClient struct {
	RestClient rest.Interface
	Ns         string
}

func (c *CodeQualityProjectClient) List(opts metav1.ListOptions) (*CodeQualityProjectList, error) {
	result := CodeQualityProjectList{}
	err := c.RestClient.
		Get().
		Namespace(c.Ns).
		Resource("codequalityprojects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)
	// fmt.Println(result)
	return &result, err
}

func (c *CodeQualityProjectClient) Get(name string, opts metav1.GetOptions) (*CodeQualityProject, error) {
	result := CodeQualityProject{}
	err := c.RestClient.
		Get().
		Namespace(c.Ns).
		Resource("codequalityprojects").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *CodeQualityProjectClient) Create(pipe *Pipeline) (*CodeQualityProject, error) {
	result := CodeQualityProject{}
	err := c.RestClient.
		Post().
		Namespace(c.Ns).
		Resource("codequalityprojects").
		Body(pipe).
		Do().
		Into(&result)

	return &result, err
}

func (c *CodeQualityProjectClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Ns).
		Resource("codequalityprojects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

func CqbWatchResources(clientSet CodeQualityProjectInterface) cache.Store {
	cqpStore, cqpController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.Watch(lo)
			},
		},
		&CodeQualityProject{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)
	stop := make(chan struct{})
	go cqpController.Run(stop)
	return cqpStore
}
