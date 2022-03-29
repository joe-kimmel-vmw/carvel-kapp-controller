// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package datapackaging

import (
	"context"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging"
	installclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/kubernetes"
)

// CarvelNoopREST is a rest implementation that proxies the rest endpoints provided by
// CRDs. This will allow us to introduce the api server without the
// complexities associated with custom storage options for now.
type CarvelNoopREST struct {
	crdClient       installclient.Interface
	nsClient        kubernetes.Interface
	globalNamespace string
}

var (
	_ rest.StandardStorage    = &CarvelNoopREST{}
	_ rest.ShortNamesProvider = &CarvelNoopREST{}
)

func NewCarvelNoopREST(crdClient installclient.Interface, nsClient kubernetes.Interface, globalNS string) *CarvelNoopREST {
	return &CarvelNoopREST{crdClient, nsClient, globalNS}
}

func (r *CarvelNoopREST) ShortNames() []string {
	return []string{"cnop"}
}

func (r *CarvelNoopREST) NamespaceScoped() bool {
	return true
}

func (r *CarvelNoopREST) New() runtime.Object {
	return &datapackaging.CarvelNoop{}
}

func (r *CarvelNoopREST) NewList() runtime.Object {
	return &datapackaging.CarvelNoopList{}
}

func (r *CarvelNoopREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	return &datapackaging.CarvelNoop{}, nil
}

func (r *CarvelNoopREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return &datapackaging.CarvelNoop{
		ObjectMeta: v1.ObjectMeta{
			Name:        "foo",
			Namespace:   "default",
			Annotations: map[string]string{"kapp.k14s.io/disable-original": ""},
		},
	}, nil
}

func (r *CarvelNoopREST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	return &datapackaging.CarvelNoopList{Items: []datapackaging.CarvelNoop{datapackaging.CarvelNoop{ObjectMeta: v1.ObjectMeta{Name: "foo", Namespace: "default"}}}}, nil
}

func (r *CarvelNoopREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	return &datapackaging.CarvelNoop{
		ObjectMeta: v1.ObjectMeta{
			Name:        "foo",
			Namespace:   "default",
			Annotations: map[string]string{"kapp.k14s.io/disable-original": ""},
		},
	}, true, nil
}

func (r *CarvelNoopREST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	return &datapackaging.CarvelNoop{}, true, nil
}

func (r *CarvelNoopREST) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternalversion.ListOptions) (runtime.Object, error) {
	return &datapackaging.CarvelNoopList{}, nil
}

func (r *CarvelNoopREST) Watch(ctx context.Context, options *internalversion.ListOptions) (watch.Interface, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageStorageClient(r.crdClient, NewPackageTranslator(namespace))

	watcher, err := client.Watch(ctx, namespace, v1.ListOptions{})
	if errors.IsNotFound(err) && namespace != r.globalNamespace {
		watcher, err = client.Watch(ctx, r.globalNamespace, v1.ListOptions{})
	}

	return watcher, err
}

func (r *CarvelNoopREST) ConvertToTable(ctx context.Context, obj runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	var table metav1.Table
	return &table, nil
}
