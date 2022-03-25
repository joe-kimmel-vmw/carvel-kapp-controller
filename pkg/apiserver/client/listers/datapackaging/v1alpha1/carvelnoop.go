// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CarvelNoopLister helps list CarvelNoops.
// All objects returned here must be treated as read-only.
type CarvelNoopLister interface {
	// List lists all CarvelNoops in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.CarvelNoop, err error)
	// CarvelNoops returns an object that can list and get CarvelNoops.
	CarvelNoops(namespace string) CarvelNoopNamespaceLister
	CarvelNoopListerExpansion
}

// foo_CarvelNoopLister implements the CarvelNoopLister interface.
type foo_CarvelNoopLister struct {
	indexer cache.Indexer
}

// NewCarvelNoopLister returns a new CarvelNoopLister.
func NewCarvelNoopLister(indexer cache.Indexer) CarvelNoopLister {
	return &foo_CarvelNoopLister{indexer: indexer}
}

// List lists all CarvelNoops in the indexer.
func (s *foo_CarvelNoopLister) List(selector labels.Selector) (ret []*v1alpha1.CarvelNoop, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.CarvelNoop))
	})
	return ret, err
}

// CarvelNoops returns an object that can list and get CarvelNoops.
func (s *foo_CarvelNoopLister) CarvelNoops(namespace string) CarvelNoopNamespaceLister {
	return foo_CarvelNoopNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CarvelNoopNamespaceLister helps list and get CarvelNoops.
// All objects returned here must be treated as read-only.
type CarvelNoopNamespaceLister interface {
	// List lists all CarvelNoops in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.CarvelNoop, err error)
	// Get retrieves the CarvelNoop from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.CarvelNoop, error)
	CarvelNoopNamespaceListerExpansion
}

// foo_CarvelNoopNamespaceLister implements the CarvelNoopNamespaceLister
// interface.
type foo_CarvelNoopNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all CarvelNoops in the indexer for a given namespace.
func (s foo_CarvelNoopNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.CarvelNoop, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.CarvelNoop))
	})
	return ret, err
}

// Get retrieves the CarvelNoop from the indexer for a given namespace and name.
func (s foo_CarvelNoopNamespaceLister) Get(name string) (*v1alpha1.CarvelNoop, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("carvelnoop"), name)
	}
	return obj.(*v1alpha1.CarvelNoop), nil
}
