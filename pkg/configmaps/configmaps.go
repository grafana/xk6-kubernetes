package configmaps

import (
	"context"
	"errors"

	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

func New(client kubernetes.Interface, metaOptions metav1.ListOptions, ctx context.Context) *ConfigMaps {
	return &ConfigMaps{
		client,
		metaOptions,
		ctx,
	}
}

type ConfigMaps struct {
	client      kubernetes.Interface
	metaOptions metav1.ListOptions
	ctx         context.Context
}

func (obj *ConfigMaps) Apply(yaml string, namespace string) (k8sTypes.ConfigMap, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	yamlobj, _, err := decode([]byte(yaml), nil, nil)
	configmap := k8sTypes.ConfigMap{}

	if err != nil {
		return configmap, err
	}

	switch yamlobj.(type) {
	case *k8sTypes.ConfigMap:
		configmap = *yamlobj.(*k8sTypes.ConfigMap)
	default:
		return configmap, errors.New("Yaml was not a ConfigMap")
	}

	cm, err := obj.client.CoreV1().ConfigMaps(namespace).Create(obj.ctx, &configmap, metav1.CreateOptions{})
	if err != nil {
		return k8sTypes.ConfigMap{}, err
	}
	return *cm, nil
}

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

func (obj *ConfigMaps) List(namespace string) ([]k8sTypes.ConfigMap, error) {
	cms, err := obj.client.CoreV1().ConfigMaps(namespace).List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.ConfigMap{}, err
	}
	return cms.Items, nil
}

func (obj *ConfigMaps) Delete(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.client.CoreV1().ConfigMaps(namespace).Delete(obj.ctx, name, opts)
}

// Deprecated: Use Delete instead.
func (obj *ConfigMaps) Kill(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.Delete(name, namespace, opts)
}

func (obj *ConfigMaps) Get(name, namespace string, opts metav1.GetOptions) (k8sTypes.ConfigMap, error) {
	cm, err := obj.client.CoreV1().ConfigMaps(namespace).Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.ConfigMap{}, err
	}
	return *cm, nil
}
