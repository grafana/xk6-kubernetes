// Package persistentvolumeclaims provides implementation of PersistentVolumeClaim resources for Kubernetes
//
// Deprecated: Use the resources package instead.
package persistentvolumeclaims

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
func New(ctx context.Context, client kubernetes.Interface, metaOptions metav1.ListOptions) *PersistentVolumeClaims {
	return &PersistentVolumeClaims{
		client,
		metaOptions,
		ctx,
	}
}

// PersistentVolumeClaims provides API for manipulating PersistentVolumeClaim resources within a Kubernetes cluster
//
// Deprecated: No longer used in favor of generic resources.
type PersistentVolumeClaims struct {
	client      kubernetes.Interface
	metaOptions metav1.ListOptions
	ctx         context.Context
}

// Apply creates the Kubernetes resource given the supplied YAML configuration
//
// Deprecated: Use resources.Apply instead.
func (obj *PersistentVolumeClaims) Apply(yaml string, namespace string) (k8sTypes.PersistentVolumeClaim, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	yamlobj, _, err := decode([]byte(yaml), nil, nil)
	persistentvolumeclaim := k8sTypes.PersistentVolumeClaim{}

	if err != nil {
		return persistentvolumeclaim, err
	}

	if pvc, ok := yamlobj.(*k8sTypes.PersistentVolumeClaim); ok {
		persistentvolumeclaim = *pvc
	} else {
		return persistentvolumeclaim, errors.New("YAML was not a PersistentVolumeClaim")
	}

	pvc, err := obj.client.CoreV1().PersistentVolumeClaims(namespace).Create(
		obj.ctx, &persistentvolumeclaim, metav1.CreateOptions{})
	if err != nil {
		return k8sTypes.PersistentVolumeClaim{}, err
	}
	return *pvc, nil
}

// Create creates the Kubernetes resource given the supplied object
//
// Deprecated: Use resources.Create instead.
func (obj *PersistentVolumeClaims) Create(
	persistentvolumeclaim k8sTypes.PersistentVolumeClaim,
	namespace string,
	opts metav1.CreateOptions,
) (k8sTypes.PersistentVolumeClaim, error) {
	pvc, err := obj.client.CoreV1().PersistentVolumeClaims(namespace).Create(obj.ctx, &persistentvolumeclaim, opts)
	if err != nil {
		return k8sTypes.PersistentVolumeClaim{}, err
	}
	return *pvc, nil
}

// List returns a collection of PersistentVolumeClaims available within the namespace
//
// Deprecated: Use resources.List instead.
func (obj *PersistentVolumeClaims) List(namespace string) ([]k8sTypes.PersistentVolumeClaim, error) {
	pvcs, err := obj.client.CoreV1().PersistentVolumeClaims(namespace).List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.PersistentVolumeClaim{}, err
	}
	return pvcs.Items, nil
}

// Delete removes the named PersistentVolumeClaims from the namespace
//
// Deprecated: Use resources.Delete instead.
func (obj *PersistentVolumeClaims) Delete(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.client.CoreV1().PersistentVolumeClaims(namespace).Delete(obj.ctx, name, opts)
}

// Get returns the named PersistentVolumeClaims instance within the namespace if available
//
// Deprecated: Use resources.Get instead.
func (obj *PersistentVolumeClaims) Get(
	name, namespace string, opts metav1.GetOptions,
) (k8sTypes.PersistentVolumeClaim, error) {
	pvc, err := obj.client.CoreV1().PersistentVolumeClaims(namespace).Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.PersistentVolumeClaim{}, err
	}

	return *pvc, nil
}
