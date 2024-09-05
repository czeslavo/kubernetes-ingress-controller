/*
Copyright 2021 Kong, Inc.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeKongVaults implements KongVaultInterface
type FakeKongVaults struct {
	Fake *FakeConfigurationV1alpha1
}

var kongvaultsResource = v1alpha1.SchemeGroupVersion.WithResource("kongvaults")

var kongvaultsKind = v1alpha1.SchemeGroupVersion.WithKind("KongVault")

// Get takes name of the kongVault, and returns the corresponding kongVault object, and an error if there is any.
func (c *FakeKongVaults) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.KongVault, err error) {
	emptyResult := &v1alpha1.KongVault{}
	obj, err := c.Fake.
		Invokes(testing.NewRootGetActionWithOptions(kongvaultsResource, name, options), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.KongVault), err
}

// List takes label and field selectors, and returns the list of KongVaults that match those selectors.
func (c *FakeKongVaults) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.KongVaultList, err error) {
	emptyResult := &v1alpha1.KongVaultList{}
	obj, err := c.Fake.
		Invokes(testing.NewRootListActionWithOptions(kongvaultsResource, kongvaultsKind, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.KongVaultList{ListMeta: obj.(*v1alpha1.KongVaultList).ListMeta}
	for _, item := range obj.(*v1alpha1.KongVaultList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested kongVaults.
func (c *FakeKongVaults) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchActionWithOptions(kongvaultsResource, opts))
}

// Create takes the representation of a kongVault and creates it.  Returns the server's representation of the kongVault, and an error, if there is any.
func (c *FakeKongVaults) Create(ctx context.Context, kongVault *v1alpha1.KongVault, opts v1.CreateOptions) (result *v1alpha1.KongVault, err error) {
	emptyResult := &v1alpha1.KongVault{}
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateActionWithOptions(kongvaultsResource, kongVault, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.KongVault), err
}

// Update takes the representation of a kongVault and updates it. Returns the server's representation of the kongVault, and an error, if there is any.
func (c *FakeKongVaults) Update(ctx context.Context, kongVault *v1alpha1.KongVault, opts v1.UpdateOptions) (result *v1alpha1.KongVault, err error) {
	emptyResult := &v1alpha1.KongVault{}
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateActionWithOptions(kongvaultsResource, kongVault, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.KongVault), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeKongVaults) UpdateStatus(ctx context.Context, kongVault *v1alpha1.KongVault, opts v1.UpdateOptions) (result *v1alpha1.KongVault, err error) {
	emptyResult := &v1alpha1.KongVault{}
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceActionWithOptions(kongvaultsResource, "status", kongVault, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.KongVault), err
}

// Delete takes name of the kongVault and deletes it. Returns an error if one occurs.
func (c *FakeKongVaults) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(kongvaultsResource, name, opts), &v1alpha1.KongVault{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeKongVaults) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionActionWithOptions(kongvaultsResource, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.KongVaultList{})
	return err
}

// Patch applies the patch and returns the patched kongVault.
func (c *FakeKongVaults) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.KongVault, err error) {
	emptyResult := &v1alpha1.KongVault{}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceActionWithOptions(kongvaultsResource, name, pt, data, opts, subresources...), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.KongVault), err
}