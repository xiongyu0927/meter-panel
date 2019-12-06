package application

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
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
	Update(*Application) (*Application, error)
	Delete(name string, options *metav1.DeleteOptions) error
}

type AppClient struct {
	restClient rest.Interface
	ns         string
}

type ApplicationInterface interface {
	App(namespace string) AppInterface
}

type ApplicationClient struct {
	restClient rest.Interface
}

func (c *ApplicationClient) App(namespace string) AppInterface {
	return &AppClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *AppClient) List(opts metav1.ListOptions) (*ApplicationList, error) {
	result := ApplicationList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("applications").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *AppClient) Get(name string, opts metav1.GetOptions) (*Application, error) {
	result := Application{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("applications").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *AppClient) Create(app *Application) (*Application, error) {
	result := Application{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("applications").
		Body(app).
		Do().
		Into(&result)

	return &result, err
}

func (c *AppClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("applications").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

func (c *AppClient) Update(app *Application) (result *Application, err error) {
	result = &Application{}
	err = c.restClient.Put().
		Namespace(c.ns).
		Resource("applications").
		Name(app.Name).
		Body(app).
		Do().
		Into(result)
	return
}

func (c *AppClient) Delete(name string, options *metav1.DeleteOptions) error {
	return c.restClient.Delete().
		Namespace(c.ns).
		Resource("applications").
		Name(name).
		Body(options).
		Do().
		Error()
}

func WatchResources(clientSet ApplicationInterface) cache.Store {
	appStore, appController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.App("").List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.App("").Watch(lo)
			},
		},
		&Application{},
		5*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)
	stop := make(chan struct{})
	go appController.Run(stop)
	return appStore
}

func NewAppClientForConfig(c *rest.Config) (*ApplicationClient, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: GroupName, Version: GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &ApplicationClient{restClient: client}, nil
}
