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
	"context"

	v1 "github.com/mneverov/cluster-cidr-controller/pkg/apis/clustercidr/v1"
	"github.com/mneverov/cluster-cidr-controller/pkg/apis/clustercidr/v1/validation"
	genscheme "github.com/mneverov/cluster-cidr-controller/pkg/client/clientset/versioned/scheme"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"
)

// clusterCIDRStrategy implements verification logic for ClusterCIDRs.
type clusterCIDRStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating clusterCIDR objects.
var Strategy = clusterCIDRStrategy{genscheme.Scheme, names.SimpleNameGenerator}

// NamespaceScoped returns false because all clusterCIDRs do not need to be within a namespace.
func (clusterCIDRStrategy) NamespaceScoped() bool {
	return false
}

func (clusterCIDRStrategy) PrepareForCreate(_ context.Context, _ runtime.Object) {}

func (clusterCIDRStrategy) PrepareForUpdate(_ context.Context, _, _ runtime.Object) {}

// Validate validates a new ClusterCIDR.
func (clusterCIDRStrategy) Validate(_ context.Context, obj runtime.Object) field.ErrorList {
	clusterCIDR := obj.(*v1.ClusterCIDR)
	return validation.ValidateClusterCIDR(clusterCIDR)
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (clusterCIDRStrategy) WarningsOnCreate(_ context.Context, _ runtime.Object) []string {
	return nil
}

// Canonicalize normalizes the object after validation.
func (clusterCIDRStrategy) Canonicalize(_ runtime.Object) {}

// AllowCreateOnUpdate is false for ClusterCIDR; this means POST is needed to create one.
func (clusterCIDRStrategy) AllowCreateOnUpdate() bool {
	return false
}

// ValidateUpdate is the default update validation for an end user.
func (clusterCIDRStrategy) ValidateUpdate(_ context.Context, obj, old runtime.Object) field.ErrorList {
	validationErrorList := validation.ValidateClusterCIDR(obj.(*v1.ClusterCIDR))
	updateErrorList := validation.ValidateClusterCIDRUpdate(obj.(*v1.ClusterCIDR), old.(*v1.ClusterCIDR))
	return append(validationErrorList, updateErrorList...)
}

// WarningsOnUpdate returns warnings for the given update.
func (clusterCIDRStrategy) WarningsOnUpdate(_ context.Context, _, _ runtime.Object) []string {
	return nil
}

// AllowUnconditionalUpdate is the default update policy for ClusterCIDR objects.
func (clusterCIDRStrategy) AllowUnconditionalUpdate() bool {
	return true
}
