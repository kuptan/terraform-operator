package v1alpha1

import (
	// "context"

	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type terraformInterface interface {
	List(opts metav1.ListOptions) (*TerraformList, error)
	Get(name string, options metav1.GetOptions) (*Terraform, error)
	Create(*Terraform) (*Terraform, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type terraformV1Alpha1InterfaceMie interface {
	Terraforms(namespace string) terraformInterface
}

type terraformV1Alpha1Client struct {
	restClient rest.Interface
}

type runClient struct {
	restClient rest.Interface
	ns         string
}

var KubeClient terraformV1Alpha1Client

func NewForConfig(c *rest.Config) (*terraformV1Alpha1Client, error) {
	AddToScheme(scheme.Scheme)

	config := *c
	config.ContentConfig.GroupVersion = &GroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)

	if err != nil {
		return nil, err
	}

	KubeClient = terraformV1Alpha1Client{restClient: client}

	return &KubeClient, nil
}

func (c *terraformV1Alpha1Client) Runs(namespace string) terraformInterface {
	return &runClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *runClient) List(opts metav1.ListOptions) (*TerraformList, error) {
	result := TerraformList{}

	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("terraform").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *runClient) Get(name string, opts metav1.GetOptions) (*Terraform, error) {
	result := Terraform{}

	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("terraform").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *runClient) Create(run *Terraform) (*Terraform, error) {
	result := Terraform{}

	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("terraform").
		Body(run).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *runClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true

	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("terraform").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.Background())
}
