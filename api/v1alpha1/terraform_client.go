package v1alpha1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type terraformInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*TerraformList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*Terraform, error)
	Create(ctx context.Context, run *Terraform, options metav1.CreateOptions) (*Terraform, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type terraformV1Alpha1Interface interface {
	Terraforms(namespace string) terraformInterface
}

type terraformV1Alpha1Client struct {
	restClient rest.Interface
}

type terraformClient struct {
	restClient rest.Interface
	ns         string
}

const k8sResourceName string = "terraforms"

var terraformKubeClient terraformV1Alpha1Client

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

	terraformKubeClient = terraformV1Alpha1Client{restClient: client}

	return &terraformKubeClient, nil
}

func (c *terraformV1Alpha1Client) Terraforms(namespace string) terraformInterface {
	return &terraformClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *terraformClient) List(ctx context.Context, opts metav1.ListOptions) (*TerraformList, error) {
	result := TerraformList{}

	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(k8sResourceName).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *terraformClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*Terraform, error) {
	result := Terraform{}

	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(k8sResourceName).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *terraformClient) Create(ctx context.Context, run *Terraform, opts metav1.CreateOptions) (*Terraform, error) {
	result := Terraform{}

	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource(k8sResourceName).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(run).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *terraformClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true

	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource(k8sResourceName).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}
