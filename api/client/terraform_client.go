package client

import (
	"context"

	v1alpha1 "github.com/kube-champ/terraform-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type TerraformInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.TerraformList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.Terraform, error)
	Create(*v1alpha1.Terraform) (*v1alpha1.Terraform, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type TerraformV1Alpha1InterfaceMie interface {
	Terraforms(namespace string) TerraformInterface
}

type TerraformV1Alpha1Client struct {
	restClient rest.Interface
}

type TerraformClient struct {
	restClient rest.Interface
	ns         string
}

const k8sResourceName string = "terraform"

var KubeClient TerraformV1Alpha1Client

func NewForConfig(c *rest.Config) (*TerraformV1Alpha1Client, error) {
	v1alpha1.AddToScheme(scheme.Scheme)

	config := *c
	config.ContentConfig.GroupVersion = &v1alpha1.GroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)

	if err != nil {
		return nil, err
	}

	KubeClient = TerraformV1Alpha1Client{restClient: client}

	return &KubeClient, nil
}

func (c *TerraformV1Alpha1Client) Terraforms(namespace string) TerraformInterface {
	return &TerraformClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *TerraformClient) List(opts metav1.ListOptions) (*v1alpha1.TerraformList, error) {
	result := v1alpha1.TerraformList{}

	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(k8sResourceName).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *TerraformClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.Terraform, error) {
	result := v1alpha1.Terraform{}

	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(k8sResourceName).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *TerraformClient) Create(run *v1alpha1.Terraform) (*v1alpha1.Terraform, error) {
	result := v1alpha1.Terraform{}

	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource(k8sResourceName).
		Body(run).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *TerraformClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true

	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource(k8sResourceName).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.Background())
}
