// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/client/clientset/versioned/scheme"
	rest "k8s.io/client-go/rest"
)

type DataV1alpha1Interface interface {
	RESTClient() rest.Interface
	CarvelNoopsGetter
	PackagesGetter
	PackageMetadatasGetter
}

// DataV1alpha1Client is used to interact with features provided by the data.packaging.carvel.dev group.
type DataV1alpha1Client struct {
	restClient rest.Interface
}

func (c *DataV1alpha1Client) CarvelNoops(namespace string) CarvelNoopInterface {
	return newCarvelNoops(c, namespace)
}

func (c *DataV1alpha1Client) Packages(namespace string) PackageInterface {
	return newPackages(c, namespace)
}

func (c *DataV1alpha1Client) PackageMetadatas(namespace string) PackageMetadataInterface {
	return newPackageMetadatas(c, namespace)
}

// NewForConfig creates a new DataV1alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*DataV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &DataV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new DataV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *DataV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new DataV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *DataV1alpha1Client {
	return &DataV1alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *DataV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
