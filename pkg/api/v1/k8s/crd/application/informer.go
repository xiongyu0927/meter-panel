package application

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type AppInterface interface {
	List(opts metav1.ListOptions) (*ApplicationList, error)
	Get(name string, options metav1.GetOptions) (*Application, error)
	Create(*Application) (*Application, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type AppClient struct {
	RestClient rest.Interface
	Ns         string
}

func (c *AppClient) List(opts metav1.ListOptions) (*ApplicationList, error) {
	result := ApplicationList{}
	err := c.RestClient.
		Get().
		Namespace(c.Ns).
		Resource("applications").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *AppClient) Get(name string, opts metav1.GetOptions) (*Application, error) {
	result := Application{}
	err := c.RestClient.
		Get().
		Namespace(c.Ns).
		Resource("applications").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *AppClient) Create(app *Application) (*Application, error) {
	result := Application{}
	err := c.RestClient.
		Post().
		Namespace(c.Ns).
		Resource("applications").
		Body(app).
		Do().
		Into(&result)

	return &result, err
}

func (c *AppClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Ns).
		Resource("applications").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

func WatchResources(clientSet AppInterface) cache.Store {
	appStore, appController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.Watch(lo)
			},
		},
		&Application{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)
	stop := make(chan struct{})
	go appController.Run(stop)
	return appStore
}
