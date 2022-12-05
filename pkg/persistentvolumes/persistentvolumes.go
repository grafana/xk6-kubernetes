// Package persistentvolumes provides implementation of PersistentVolume resources for Kubernetes
//
// Deprecated: Use the resources package instead.
package persistentvolumes

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
func New(ctx context.Context, client kubernetes.Interface, metaOptions metav1.ListOptions) *PersistentVolumes {
	return &PersistentVolumes{
		client,
		metaOptions,
		ctx,
	}
}

// PersistentVolumes provides API for manipulating PersistentVolume resources within a Kubernetes cluster
//
// Deprecated: No longer used in favor of generic resources.
type PersistentVolumes struct {
	client      kubernetes.Interface
	metaOptions metav1.ListOptions
	ctx         context.Context
}

// Apply creates the Kubernetes resource given the supplied YAML configuration
//
// Deprecated: Use resources.Apply instead.
func (obj *PersistentVolumes) Apply(yaml string) (k8sTypes.PersistentVolume, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	yamlobj, _, err := decode([]byte(yaml), nil, nil)
	persistentvolume := k8sTypes.PersistentVolume{}

	if err != nil {
		return persistentvolume, err
	}

	if pv, ok := yamlobj.(*k8sTypes.PersistentVolume); ok {
		persistentvolume = *pv
	} else {
		return persistentvolume, errors.New("YAML was not a PersistentVolume")
	}

	pv, err := obj.client.CoreV1().PersistentVolumes().Create(obj.ctx, &persistentvolume, metav1.CreateOptions{})
	if err != nil {
		return k8sTypes.PersistentVolume{}, err
	}
	return *pv, nil
}

// Create creates the Kubernetes resource given the supplied object
//
// Deprecated: Use resources.Create instead.
func (obj *PersistentVolumes) Create(
	persistentvolume k8sTypes.PersistentVolume,
	opts metav1.CreateOptions,
) (k8sTypes.PersistentVolume, error) {
	pv, err := obj.client.CoreV1().PersistentVolumes().Create(obj.ctx, &persistentvolume, opts)
	if err != nil {
		return k8sTypes.PersistentVolume{}, err
	}
	return *pv, nil
}

// List returns a collection of PersistentVolumes available within the cluster
//
// Deprecated: Use resources.List instead.
func (obj *PersistentVolumes) List() ([]k8sTypes.PersistentVolume, error) {
	pvs, err := obj.client.CoreV1().PersistentVolumes().List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.PersistentVolume{}, err
	}
	return pvs.Items, nil
}

// Delete removes the named PersistentVolumes from the cluster
//
// Deprecated: Use resources.Delete instead.
func (obj *PersistentVolumes) Delete(name string, opts metav1.DeleteOptions) error {
	return obj.client.CoreV1().PersistentVolumes().Delete(obj.ctx, name, opts)
}

// Get returns the named PersistentVolumes instance within the cluster if available
//
// Deprecated: Use resources.Get instead.
func (obj *PersistentVolumes) Get(name string, opts metav1.GetOptions) (k8sTypes.PersistentVolume, error) {
	pv, err := obj.client.CoreV1().PersistentVolumes().Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.PersistentVolume{}, err
	}

	return *pv, nil
}
