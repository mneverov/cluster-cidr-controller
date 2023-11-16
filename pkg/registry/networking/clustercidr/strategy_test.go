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

package clustercidr

import (
	"testing"

	v1 "github.com/mneverov/cluster-cidr-controller/pkg/apis/clustercidr/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

func newClusterCIDR() v1.ClusterCIDR {
	return v1.ClusterCIDR{
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo",
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

func TestClusterCIDRStrategy(t *testing.T) {
	ctx := genericapirequest.NewDefaultContext()
	apiRequest := genericapirequest.RequestInfo{
		APIGroup:   "networking.x-k8s.io",
		APIVersion: "v1",
		Resource:   "clustercidrs",
	}
	ctx = genericapirequest.WithRequestInfo(ctx, &apiRequest)
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1.AddToScheme(scheme))
	strategy := NewStrategy(scheme)

	if strategy.NamespaceScoped() {
		t.Errorf("ClusterCIDRs must be cluster scoped")
	}
	if strategy.AllowCreateOnUpdate() {
		t.Errorf("ClusterCIDRs should not allow create on update")
	}

	ccc := newClusterCIDR()
	strategy.PrepareForCreate(ctx, &ccc)

	errs := strategy.Validate(ctx, &ccc)
	if len(errs) != 0 {
		t.Errorf("Unexpected error validating %v", errs)
	}
	invalidCCC := newClusterCIDR()
	invalidCCC.ResourceVersion = "4"
	invalidCCC.Spec = v1.ClusterCIDRSpec{}
	strategy.PrepareForUpdate(ctx, &invalidCCC, &ccc)
	errs = strategy.ValidateUpdate(ctx, &invalidCCC, &ccc)
	if len(errs) == 0 {
		t.Errorf("Expected a validation error")
	}
	if invalidCCC.ResourceVersion != "4" {
		t.Errorf("Incoming resource version on update should not be mutated")
	}
}
