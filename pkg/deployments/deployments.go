// Package deployments provides implementation of Deployment resources for Kubernetes
package deployments

import (
	"context"
	"errors"

	k8sTypes "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

// New creates a new instance backed by the provided client
func New(ctx context.Context, client kubernetes.Interface, metaOptions metav1.ListOptions) *Deployments {
	return &Deployments{
		client,
		metaOptions,
		ctx,
	}
}

// Deployments provides API for manipulating Deployment resources within a Kubernetes cluster
type Deployments struct {
	client      kubernetes.Interface
	metaOptions metav1.ListOptions
	ctx         context.Context
}

// Apply creates the Kubernetes resource given the supplied YAML configuration
func (obj *Deployments) Apply(yaml string, namespace string) (k8sTypes.Deployment, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	yamlobj, _, err := decode([]byte(yaml), nil, nil)
	deployment := k8sTypes.Deployment{}

	if err != nil {
		return deployment, err
	}

	if dep, ok := yamlobj.(*k8sTypes.Deployment); ok {
		deployment = *dep
	} else {
		return deployment, errors.New("YAML was not a Deployment")
	}

	dep, err := obj.client.AppsV1().Deployments(namespace).Create(obj.ctx, &deployment, metav1.CreateOptions{})
	if err != nil {
		return k8sTypes.Deployment{}, err
	}
	return *dep, nil
}

// Create creates the Kubernetes resource given the supplied object
func (obj *Deployments) Create(
	deployment k8sTypes.Deployment,
	namespace string,
	opts metav1.CreateOptions,
) (k8sTypes.Deployment, error) {
	dep, err := obj.client.AppsV1().Deployments(namespace).Create(obj.ctx, &deployment, opts)
	if err != nil {
		return k8sTypes.Deployment{}, err
	}
	return *dep, nil
}

// List returns a collection of Deployments available within the namespace
func (obj *Deployments) List(namespace string) ([]k8sTypes.Deployment, error) {
	dps, err := obj.client.AppsV1().Deployments(namespace).List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.Deployment{}, err
	}
	return dps.Items, nil
}

// Delete removes the named Deployment from the namespace
func (obj *Deployments) Delete(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.client.AppsV1().Deployments(namespace).Delete(obj.ctx, name, opts)
}

// Kill removes the named Deployment from the namespace
// Deprecated: Use Delete instead.
func (obj *Deployments) Kill(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.Delete(name, namespace, opts)
}

// Get returns the named Deployments instance within the namespace if available
func (obj *Deployments) Get(name, namespace string, opts metav1.GetOptions) (k8sTypes.Deployment, error) {
	dp, err := obj.client.AppsV1().Deployments(namespace).Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.Deployment{}, err
	}
	return *dp, nil
}
