// Package configmaps provides implementation of ConfigMap resources for Kubernetes
package configmaps

import (
	"context"
	"errors"

	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

// New creates a new instance backed by the provided client
func New(ctx context.Context, client kubernetes.Interface, metaOptions metav1.ListOptions) *ConfigMaps {
	return &ConfigMaps{
		client,
		metaOptions,
		ctx,
	}
}

// ConfigMaps provides API for manipulating ConfigMap resources within a Kubernetes cluster
type ConfigMaps struct {
	client      kubernetes.Interface
	metaOptions metav1.ListOptions
	ctx         context.Context
}

// Apply creates the Kubernetes resource given the supplied YAML configuration
func (obj *ConfigMaps) Apply(yaml string, namespace string) (k8sTypes.ConfigMap, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	yamlobj, _, err := decode([]byte(yaml), nil, nil)
	configmap := k8sTypes.ConfigMap{}

	if err != nil {
		return configmap, err
	}

	if cm, ok := yamlobj.(*k8sTypes.ConfigMap); ok {
		configmap = *cm
	} else {
		return k8sTypes.ConfigMap{}, errors.New("YAML was not a ConfigMap")
	}

	cm, err := obj.client.CoreV1().ConfigMaps(namespace).Create(obj.ctx, &configmap, metav1.CreateOptions{})
	if err != nil {
		return k8sTypes.ConfigMap{}, err
	}
	return *cm, nil
}

// Create creates the Kubernetes resource given the supplied object
func (obj *ConfigMaps) Create(
	configMap k8sTypes.ConfigMap,
	namespace string,
	opts metav1.CreateOptions,
) (k8sTypes.ConfigMap, error) {
	cm, err := obj.client.CoreV1().ConfigMaps(namespace).Create(obj.ctx, &configMap, opts)
	if err != nil {
		return k8sTypes.ConfigMap{}, err
	}
	return *cm, nil
}

// List returns a collection of ConfigMaps available within the namespace
func (obj *ConfigMaps) List(namespace string) ([]k8sTypes.ConfigMap, error) {
	cms, err := obj.client.CoreV1().ConfigMaps(namespace).List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.ConfigMap{}, err
	}
	return cms.Items, nil
}

// Delete removes the named ConfigMap from the namespace
func (obj *ConfigMaps) Delete(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.client.CoreV1().ConfigMaps(namespace).Delete(obj.ctx, name, opts)
}

// Kill removes the named ConfigMap from the namespace
// Deprecated: Use Delete instead.
func (obj *ConfigMaps) Kill(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.Delete(name, namespace, opts)
}

// Get returns the named ConfigMaps instance within the namespace if available
func (obj *ConfigMaps) Get(name, namespace string, opts metav1.GetOptions) (k8sTypes.ConfigMap, error) {
	cm, err := obj.client.CoreV1().ConfigMaps(namespace).Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.ConfigMap{}, err
	}
	return *cm, nil
}
