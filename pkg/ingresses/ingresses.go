// Package ingresses provides implementation of Ingress resources for Kubernetes
package ingresses

import (
	"context"
	"errors"

	k8sTypes "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

// New creates a new instance backed by the provided client
func New(ctx context.Context, client kubernetes.Interface, metaOptions metav1.ListOptions) *Ingresses {
	return &Ingresses{
		client,
		metaOptions,
		ctx,
	}
}

// Ingresses provides API for manipulating Ingress resources within a Kubernetes cluster
type Ingresses struct {
	client      kubernetes.Interface
	metaOptions metav1.ListOptions
	ctx         context.Context
}

// Apply creates the Kubernetes resource given the supplied YAML configuration
func (obj *Ingresses) Apply(yaml string, namespace string) (k8sTypes.Ingress, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	yamlobj, _, err := decode([]byte(yaml), nil, nil)
	ingress := k8sTypes.Ingress{}

	if err != nil {
		return ingress, err
	}

	if ing, ok := yamlobj.(*k8sTypes.Ingress); ok {
		ingress = *ing
	} else {
		return ingress, errors.New("YAML was not an Ingress")
	}

	ing, err := obj.client.NetworkingV1().Ingresses(namespace).Create(obj.ctx, &ingress, metav1.CreateOptions{})
	if err != nil {
		return k8sTypes.Ingress{}, err
	}
	return *ing, nil
}

// Create creates the Kubernetes resource given the supplied object
func (obj *Ingresses) Create(
	ingress k8sTypes.Ingress,
	namespace string,
	opts metav1.CreateOptions,
) (k8sTypes.Ingress, error) {
	ing, err := obj.client.NetworkingV1().Ingresses(namespace).Create(obj.ctx, &ingress, opts)
	if err != nil {
		return k8sTypes.Ingress{}, err
	}
	return *ing, nil
}

// List returns a collection of Ingresses available within the namespace
func (obj *Ingresses) List(namespace string) ([]k8sTypes.Ingress, error) {
	ings, err := obj.client.NetworkingV1().Ingresses(namespace).List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.Ingress{}, err
	}
	return ings.Items, nil
}

// Delete removes the named Ingress from the namespace
func (obj *Ingresses) Delete(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.client.NetworkingV1().Ingresses(namespace).Delete(obj.ctx, name, opts)
}

// Kill removes the named Ingress from the namespace
// Deprecated: Use Delete instead.
func (obj *Ingresses) Kill(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.Delete(name, namespace, opts)
}

// Get returns the named Ingresses instance within the namespace if available
func (obj *Ingresses) Get(name, namespace string, opts metav1.GetOptions) (k8sTypes.Ingress, error) {
	ing, err := obj.client.NetworkingV1().Ingresses(namespace).Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.Ingress{}, err
	}

	return *ing, nil
}
