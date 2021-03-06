/*
Copyright The Kubernetes Authors.

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

package v1

import (
	"context"
	"time"

	v1 "github.com/Shanghai-Lunara/helixsaga-operator/pkg/apis/helixsaga/v1"
	scheme "github.com/Shanghai-Lunara/helixsaga-operator/pkg/generated/helixsaga/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// HelixSagasGetter has a method to return a HelixSagaInterface.
// A group's client should implement this interface.
type HelixSagasGetter interface {
	HelixSagas(namespace string) HelixSagaInterface
}

// HelixSagaInterface has methods to work with HelixSaga resources.
type HelixSagaInterface interface {
	Create(ctx context.Context, helixSaga *v1.HelixSaga, opts metav1.CreateOptions) (*v1.HelixSaga, error)
	Update(ctx context.Context, helixSaga *v1.HelixSaga, opts metav1.UpdateOptions) (*v1.HelixSaga, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.HelixSaga, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.HelixSagaList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.HelixSaga, err error)
	HelixSagaExpansion
}

// helixSagas implements HelixSagaInterface
type helixSagas struct {
	client rest.Interface
	ns     string
}

// newHelixSagas returns a HelixSagas
func newHelixSagas(c *NevercaseV1Client, namespace string) *helixSagas {
	return &helixSagas{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the helixSaga, and returns the corresponding helixSaga object, and an error if there is any.
func (c *helixSagas) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.HelixSaga, err error) {
	result = &v1.HelixSaga{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("helixsagas").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of HelixSagas that match those selectors.
func (c *helixSagas) List(ctx context.Context, opts metav1.ListOptions) (result *v1.HelixSagaList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.HelixSagaList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("helixsagas").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested helixSagas.
func (c *helixSagas) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("helixsagas").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a helixSaga and creates it.  Returns the server's representation of the helixSaga, and an error, if there is any.
func (c *helixSagas) Create(ctx context.Context, helixSaga *v1.HelixSaga, opts metav1.CreateOptions) (result *v1.HelixSaga, err error) {
	result = &v1.HelixSaga{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("helixsagas").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(helixSaga).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a helixSaga and updates it. Returns the server's representation of the helixSaga, and an error, if there is any.
func (c *helixSagas) Update(ctx context.Context, helixSaga *v1.HelixSaga, opts metav1.UpdateOptions) (result *v1.HelixSaga, err error) {
	result = &v1.HelixSaga{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("helixsagas").
		Name(helixSaga.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(helixSaga).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the helixSaga and deletes it. Returns an error if one occurs.
func (c *helixSagas) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("helixsagas").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *helixSagas) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("helixsagas").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched helixSaga.
func (c *helixSagas) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.HelixSaga, err error) {
	result = &v1.HelixSaga{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("helixsagas").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
