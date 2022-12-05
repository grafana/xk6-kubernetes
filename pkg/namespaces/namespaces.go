// Package namespaces provides implementation of Namespace resources for Kubernetes
//
// Deprecated: Use the resources package instead.
package namespaces

import (
	"context"
	"errors"

	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

// New creates a new instance backed by the provided client
//
// Deprecated: No longer used.
func New(ctx context.Context, client kubernetes.Interface, metaOptions metav1.ListOptions) *Namespaces {
	return &Namespaces{
		client,
		metaOptions,
		ctx,
	}
}

// Namespaces provides API for manipulating Namespace resources within a Kubernetes cluster
//
// Deprecated: No longer used in favor of generic resources.
type Namespaces struct {
	client      kubernetes.Interface
	metaOptions metav1.ListOptions
	ctx         context.Context
}

// Apply creates the Kubernetes resource given the supplied YAML configuration
//
// Deprecated: Use resources.Apply instead.
func (obj *Namespaces) Apply(yaml string) (k8sTypes.Namespace, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	yamlobj, _, err := decode([]byte(yaml), nil, nil)
	namespace := k8sTypes.Namespace{}

	if err != nil {
		return namespace, err
	}

	if ns, ok := yamlobj.(*k8sTypes.Namespace); ok {
		namespace = *ns
	} else {
		return namespace, errors.New("YAML was not a Namespace")
	}

	ns, err := obj.client.CoreV1().Namespaces().Create(obj.ctx, &namespace, metav1.CreateOptions{})
	if err != nil {
		return k8sTypes.Namespace{}, err
	}
	return *ns, nil
}

// Create creates the Kubernetes resource given the supplied object
//
// Deprecated: Use resources.Create instead.
func (obj *Namespaces) Create(
	namespace k8sTypes.Namespace,
	opts metav1.CreateOptions,
) (k8sTypes.Namespace, error) {
	ns, err := obj.client.CoreV1().Namespaces().Create(obj.ctx, &namespace, opts)
	if err != nil {
		return k8sTypes.Namespace{}, err
	}
	return *ns, nil
}

// List returns a collection of Namespaces available within the cluster
//
// Deprecated: Use resources.List instead.
func (obj *Namespaces) List() ([]k8sTypes.Namespace, error) {
	nss, err := obj.client.CoreV1().Namespaces().List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.Namespace{}, err
	}
	return nss.Items, nil
}

// Delete removes the named Namespace from the cluster
//
// Deprecated: Use resources.Delete instead.
func (obj *Namespaces) Delete(name string, opts metav1.DeleteOptions) error {
	return obj.client.CoreV1().Namespaces().Delete(obj.ctx, name, opts)
}

// Kill removes the named Namespace from the cluster
//
// Deprecated: Use resources.Delete instead.
func (obj *Namespaces) Kill(name string, opts metav1.DeleteOptions) error {
	return obj.Delete(name, opts)
}

// Get returns the named Namespaces instance within the cluster if available
//
// Deprecated: Use resources.Get instead.
func (obj *Namespaces) Get(name string, opts metav1.GetOptions) (k8sTypes.Namespace, error) {
	ns, err := obj.client.CoreV1().Namespaces().Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.Namespace{}, err
	}

	return *ns, nil
}
