/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package storage

import (
	"reflect"
	"strings"
	"testing"

	v1 "github.com/mneverov/cluster-cidr-controller/pkg/apis/clustercidr/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistrytest "k8s.io/apiserver/pkg/registry/generic/testing"
	"k8s.io/apiserver/pkg/server/options"
	"k8s.io/apiserver/pkg/server/resourceconfig"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	etcd3testing "k8s.io/apiserver/pkg/storage/etcd3/testing"
	"k8s.io/apiserver/pkg/storage/storagebackend"
)

var (
	namespace = metav1.NamespaceNone
	name      = "foo-clustercidr"
)

func TestCreate(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()

	test := genericregistrytest.New(t, storage.Store)
	test = test.ClusterScope()
	validCC := validClusterCIDR()
	noCIDRCC := validClusterCIDR()
	noCIDRCC.Spec.IPv4 = ""
	noCIDRCC.Spec.IPv6 = ""
	invalidCCPerNodeHostBits := validClusterCIDR()
	invalidCCPerNodeHostBits.Spec.PerNodeHostBits = 100
	invalidCCCIDR := validClusterCIDR()
	invalidCCCIDR.Spec.IPv6 = "10.1.0.0/16"

	test.TestCreate(
		// valid
		validCC,
		// invalid
		noCIDRCC,
		invalidCCPerNodeHostBits,
		invalidCCCIDR,
	)
}

func TestUpdate(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store)
	test = test.ClusterScope()
	test.TestUpdate(
		// valid
		validClusterCIDR(),
		// updateFunc
		func(obj runtime.Object) runtime.Object {
			object := obj.(*v1.ClusterCIDR)
			object.Finalizers = []string{"test.k8s.io/test-finalizer"}
			return object
		},
		// invalid updateFunc: ObjectMeta is not to be tampered with.
		func(obj runtime.Object) runtime.Object {
			object := obj.(*v1.ClusterCIDR)
			object.Name = ""
			return object
		},
	)
}

func TestDelete(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store)
	test = test.ClusterScope()
	test.TestDelete(validClusterCIDR())
}

func TestGet(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store)
	test = test.ClusterScope()
	test.TestGet(validClusterCIDR())
}

func TestList(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store)
	test = test.ClusterScope()
	test.TestList(validClusterCIDR())
}

func TestWatch(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store)
	test = test.ClusterScope()
	test.TestWatch(
		validClusterCIDR(),
		// matching labels
		[]labels.Set{},
		// not matching labels
		[]labels.Set{
			{"a": "c"},
			{"foo": "bar"},
		},
		// matching fields
		[]fields.Set{
			{"metadata.name": name},
		},
		// not matching fields
		[]fields.Set{
			{"metadata.name": "bar"},
			{"name": name},
		},
	)
}

func TestShortNames(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	expected := []string{"cc"}

	actual := storage.ShortNames()
	ok := reflect.DeepEqual(actual, expected)
	if !ok {
		t.Errorf("short names not equal. expected = %v actual = %v", expected, actual)
	}
}

func newStorage(t *testing.T) (*REST, *etcd3testing.EtcdTestServer) {
	scheme := runtime.NewScheme()
	utilruntime.Must(v1.AddToScheme(scheme))
	etcdStorage, server := newEtcdStorageForResource(t, scheme, v1.Resource("clustercidrs"))
	restOptions := generic.RESTOptions{
		StorageConfig:           etcdStorage,
		Decorator:               generic.UndecoratedStorage,
		DeleteCollectionWorkers: 1,
		ResourcePrefix:          "clustercidrs",
	}

	clusterCIDRStorage, err := NewREST(scheme, restOptions)
	if err != nil {
		t.Fatalf("unexpected error from REST storage: %v", err)
	}
	return clusterCIDRStorage, server
}

func newEtcdStorageForResource(t *testing.T, scheme *runtime.Scheme, resource schema.GroupResource) (*storagebackend.ConfigForResource, *etcd3testing.EtcdTestServer) {
	t.Helper()

	server, config := etcd3testing.NewUnsecuredEtcd3TestClientServer(t)
	etcdOptions := options.NewEtcdOptions(config)
	factory, err := newStorageFactory(scheme, etcdOptions)
	if err != nil {
		t.Fatalf("Error while making storage factory: %v", err)
	}
	resourceConfig, err := factory.NewConfig(resource)
	if err != nil {
		t.Fatalf("Error while finding storage destination: %v", err)
	}
	return resourceConfig, server
}

// new returns a new storage factory created from the completed storage factory configuration.
func newStorageFactory(scheme *runtime.Scheme, opts *options.EtcdOptions) (*serverstorage.DefaultStorageFactory, error) {
	codecs := serializer.NewCodecFactory(scheme)
	defaultResourceEncoding := serverstorage.NewDefaultResourceEncodingConfig(scheme)
	resourceEncodingOverrides := []schema.GroupVersionResource{v1.Resource("clustercidrs").WithVersion("v1")}
	resourceEncodingConfig := resourceconfig.MergeResourceEncodingConfigs(defaultResourceEncoding, resourceEncodingOverrides)
	storageFactory := serverstorage.NewDefaultStorageFactory(
		opts.StorageConfig,
		opts.DefaultStorageMediaType,
		codecs,
		resourceEncodingConfig,
		serverstorage.NewResourceConfig(),
		nil) // SpecialDefaultResourcePrefixes

	for _, override := range opts.EtcdServersOverrides {
		tokens := strings.Split(override, "#")
		apiresource := strings.Split(tokens[0], "/")

		group := apiresource[0]
		resource := apiresource[1]
		groupResource := schema.GroupResource{Group: group, Resource: resource}

		servers := strings.Split(tokens[1], ";")
		storageFactory.SetEtcdLocation(groupResource, servers)
	}
	return storageFactory, nil
}

func validClusterCIDR() *v1.ClusterCIDR {
	return newClusterCIDR()
}

func newClusterCIDR() *v1.ClusterCIDR {
	return &v1.ClusterCIDR{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1.ClusterCIDRSpec{
			PerNodeHostBits: int32(8),
			IPv4:            "10.1.0.0/16",
			IPv6:            "fd00:1:1::/64",
			NodeSelector: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{
					{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      "foo",
								Operator: corev1.NodeSelectorOpIn,
								Values:   []string{"bar"},
							},
						},
					},
				},
			},
		},
	}
}
