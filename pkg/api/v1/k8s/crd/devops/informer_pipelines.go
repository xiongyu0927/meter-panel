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

type PipelineInterface interface {
	List(opts metav1.ListOptions) (*PipelineList, error)
	Get(name string, options metav1.GetOptions) (*Pipeline, error)
	Create(pipe *Pipeline) (*Pipeline, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type PipelineClient struct {
	RestClient rest.Interface
	Ns         string
}

func (c *PipelineClient) List(opts metav1.ListOptions) (*PipelineList, error) {
	result := PipelineList{}
	err := c.RestClient.
		Get().
		Namespace(c.Ns).
		Resource("pipelines").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)
	// fmt.Println(result)
	return &result, err
}

func (c *PipelineClient) Get(name string, opts metav1.GetOptions) (*Pipeline, error) {
	result := Pipeline{}
	err := c.RestClient.
		Get().
		Namespace(c.Ns).
		Resource("pipelines").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *PipelineClient) Create(pipe *Pipeline) (*Pipeline, error) {
	result := Pipeline{}
	err := c.RestClient.
		Post().
		Namespace(c.Ns).
		Resource("pipelines").
		Body(pipe).
		Do().
		Into(&result)

	return &result, err
}

func (c *PipelineClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Ns).
		Resource("pipelines").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

func PipeWatchResources(clientSet PipelineInterface) cache.Store {
	pipelineStore, pipelineController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.Watch(lo)
			},
		},
		&Pipeline{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)
	stop := make(chan struct{})
	go pipelineController.Run(stop)
	return pipelineStore
}
