// Package services provides implementation of Service resources for Kubernetes
package services

import (
	"context"
	"errors"

	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

// New creates a new instance backed by the provided client
func New(ctx context.Context, client kubernetes.Interface, metaOptions metav1.ListOptions) *Services {
	return &Services{
		client,
		metaOptions,
		ctx,
	}
}

// Services provides API for manipulating Service resources within a Kubernetes cluster
type Services struct {
	client      kubernetes.Interface
	metaOptions metav1.ListOptions
	ctx         context.Context
}

// Apply creates the Kubernetes resource given the supplied YAML configuration
func (obj *Services) Apply(yaml string, namespace string) (k8sTypes.Service, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	yamlobj, _, err := decode([]byte(yaml), nil, nil)
	service := k8sTypes.Service{}

	if err != nil {
		return service, err
	}

	if svc, ok := yamlobj.(*k8sTypes.Service); ok {
		service = *svc
	} else {
		return service, errors.New("YAML was not a Service")
	}

	svc, err := obj.client.CoreV1().Services(namespace).Create(obj.ctx, &service, metav1.CreateOptions{})
	if err != nil {
		return k8sTypes.Service{}, err
	}
	return *svc, nil
}

// Create creates the Kubernetes resource given the supplied object
func (obj *Services) Create(
	service k8sTypes.Service,
	namespace string,
	opts metav1.CreateOptions,
) (k8sTypes.Service, error) {
	svc, err := obj.client.CoreV1().Services(namespace).Create(obj.ctx, &service, opts)
	if err != nil {
		return k8sTypes.Service{}, err
	}
	return *svc, nil
}

// List returns a collection of Services available within the namespace
func (obj *Services) List(namespace string) ([]k8sTypes.Service, error) {
	svcs, err := obj.client.CoreV1().Services(namespace).List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.Service{}, err
	}
	return svcs.Items, nil
}

// Delete removes the named Service from the namespace
func (obj *Services) Delete(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.client.CoreV1().Services(namespace).Delete(obj.ctx, name, opts)
}

// Kill removes the named Service from the namespace
// Deprecated: Use Delete instead.
func (obj *Services) Kill(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.Delete(name, namespace, opts)
}

// Get returns the named Services instance within the namespace if available
func (obj *Services) Get(name, namespace string, opts metav1.GetOptions) (k8sTypes.Service, error) {
	svc, err := obj.client.CoreV1().Services(namespace).Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.Service{}, err
	}

	return *svc, nil
}
