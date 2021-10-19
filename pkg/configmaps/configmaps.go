package configmaps

import (
	"context"

	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func New(client *kubernetes.Clientset, metaOptions metav1.ListOptions, ctx context.Context) *ConfigMaps {
	return &ConfigMaps{
		client,
		metaOptions,
		ctx,
	}
}

type ConfigMaps struct {
	client      *kubernetes.Clientset
	metaOptions metav1.ListOptions
	ctx         context.Context
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

func (obj *ConfigMaps) Kill(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.client.CoreV1().ConfigMaps(namespace).Delete(obj.ctx, name, opts)
}

func (obj *ConfigMaps) Get(name, namespace string, opts metav1.GetOptions) (k8sTypes.ConfigMap, error) {
	cm, err := obj.client.CoreV1().ConfigMaps(namespace).Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.ConfigMap{}, err
	}
	return *cm, nil
}
