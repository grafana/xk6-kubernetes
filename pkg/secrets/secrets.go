// Package secrets provides implementation of Secret resources for Kubernetes
//
// Deprecated: Use the resources package instead.
package secrets

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
func New(ctx context.Context, client kubernetes.Interface, metaOptions metav1.ListOptions) *Secrets {
	return &Secrets{
		client,
		metaOptions,
		ctx,
	}
}

// Secrets provides API for manipulating Secret resources within a Kubernetes cluster
//
// Deprecated: No longer used in favor of generic resources.
type Secrets struct {
	client      kubernetes.Interface
	metaOptions metav1.ListOptions
	ctx         context.Context
}

// Apply creates the Kubernetes resource given the supplied YAML configuration
//
// Deprecated: Use resources.Apply instead.
func (obj *Secrets) Apply(yaml string, namespace string) (k8sTypes.Secret, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	yamlobj, _, err := decode([]byte(yaml), nil, nil)
	secret := k8sTypes.Secret{}

	if err != nil {
		return secret, err
	}

	if scrt, ok := yamlobj.(*k8sTypes.Secret); ok {
		secret = *scrt
	} else {
		return secret, errors.New("YAML was not a Secret")
	}

	scrt, err := obj.client.CoreV1().Secrets(namespace).Create(obj.ctx, &secret, metav1.CreateOptions{})
	if err != nil {
		return k8sTypes.Secret{}, err
	}
	return *scrt, nil
}

// Create creates the Kubernetes resource given the supplied object
//
// Deprecated: Use resources.Create instead.
func (obj *Secrets) Create(
	secret k8sTypes.Secret,
	namespace string,
	opts metav1.CreateOptions,
) (k8sTypes.Secret, error) {
	scrt, err := obj.client.CoreV1().Secrets(namespace).Create(obj.ctx, &secret, opts)
	if err != nil {
		return k8sTypes.Secret{}, err
	}
	return *scrt, nil
}

// List returns a collection of Secrets available within the namespace
//
// Deprecated: Use resources.List instead.
func (obj *Secrets) List(namespace string) ([]k8sTypes.Secret, error) {
	scrts, err := obj.client.CoreV1().Secrets(namespace).List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.Secret{}, err
	}
	return scrts.Items, nil
}

// Delete removes the named Secret from the namespace
//
// Deprecated: Use resources.Delete instead.
func (obj *Secrets) Delete(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.client.CoreV1().Secrets(namespace).Delete(obj.ctx, name, opts)
}

// Kill removes the named Secret from the namespace
//
// Deprecated: Use resources.Delete instead.
func (obj *Secrets) Kill(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.Delete(name, namespace, opts)
}

// Get returns the named Secrets instance within the namespace if available
//
// Deprecated: Use resources.Get instead.
func (obj *Secrets) Get(name, namespace string, opts metav1.GetOptions) (k8sTypes.Secret, error) {
	scrt, err := obj.client.CoreV1().Secrets(namespace).Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.Secret{}, err
	}

	return *scrt, nil
}
