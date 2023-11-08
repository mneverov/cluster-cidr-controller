// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "github.com/mneverov/cluster-cidr-controller/pkg/apis/clustercidr/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeClusterCIDRs implements ClusterCIDRInterface
type FakeClusterCIDRs struct {
	Fake *FakeNetworkingV1
}

var clustercidrsResource = v1.SchemeGroupVersion.WithResource("clustercidrs")

var clustercidrsKind = v1.SchemeGroupVersion.WithKind("ClusterCIDR")

// Get takes name of the clusterCIDR, and returns the corresponding clusterCIDR object, and an error if there is any.
func (c *FakeClusterCIDRs) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.ClusterCIDR, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(clustercidrsResource, name), &v1.ClusterCIDR{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1.ClusterCIDR), err
}

// List takes label and field selectors, and returns the list of ClusterCIDRs that match those selectors.
func (c *FakeClusterCIDRs) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ClusterCIDRList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(clustercidrsResource, clustercidrsKind, opts), &v1.ClusterCIDRList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1.ClusterCIDRList{ListMeta: obj.(*v1.ClusterCIDRList).ListMeta}
	for _, item := range obj.(*v1.ClusterCIDRList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterCIDRs.
func (c *FakeClusterCIDRs) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(clustercidrsResource, opts))
}

// Create takes the representation of a clusterCIDR and creates it.  Returns the server's representation of the clusterCIDR, and an error, if there is any.
func (c *FakeClusterCIDRs) Create(ctx context.Context, clusterCIDR *v1.ClusterCIDR, opts metav1.CreateOptions) (result *v1.ClusterCIDR, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(clustercidrsResource, clusterCIDR), &v1.ClusterCIDR{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1.ClusterCIDR), err
}

// Update takes the representation of a clusterCIDR and updates it. Returns the server's representation of the clusterCIDR, and an error, if there is any.
func (c *FakeClusterCIDRs) Update(ctx context.Context, clusterCIDR *v1.ClusterCIDR, opts metav1.UpdateOptions) (result *v1.ClusterCIDR, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(clustercidrsResource, clusterCIDR), &v1.ClusterCIDR{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1.ClusterCIDR), err
}

// Delete takes name of the clusterCIDR and deletes it. Returns an error if one occurs.
func (c *FakeClusterCIDRs) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(clustercidrsResource, name, opts), &v1.ClusterCIDR{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClusterCIDRs) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(clustercidrsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1.ClusterCIDRList{})
	return err
}

// Patch applies the patch and returns the patched clusterCIDR.
func (c *FakeClusterCIDRs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ClusterCIDR, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(clustercidrsResource, name, pt, data, subresources...), &v1.ClusterCIDR{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1.ClusterCIDR), err
}
