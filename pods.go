package kubernetes

import (
	"context"

	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodsNamespace struct {
	Client      *kubernetes.Clientset
	metaOptions metav1.ListOptions
	ctx         context.Context
}

func (obj *PodsNamespace) List(namespace string) ([]k8sTypes.Pod, error) {
	pods, err := obj.Client.CoreV1().Pods(namespace).List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.Pod{}, err
	}
	return pods.Items, nil
}

func (obj *PodsNamespace) Kill(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.Client.CoreV1().Pods(namespace).Delete(obj.ctx, name, opts)
}
