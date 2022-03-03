package pods

import (
	"context"
	"errors"

	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func New(client *kubernetes.Clientset, metaOptions metav1.ListOptions, ctx context.Context) *Pods {
	return &Pods{
		client,
		metaOptions,
		ctx,
	}
}

type Pods struct {
	client      *kubernetes.Clientset
	metaOptions metav1.ListOptions
	ctx         context.Context
}

func (obj *Pods) List(namespace string) ([]k8sTypes.Pod, error) {
	pods, err := obj.client.CoreV1().Pods(namespace).List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.Pod{}, err
	}
	return pods.Items, nil
}

func (obj *Pods) Delete(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.client.CoreV1().Pods(namespace).Delete(obj.ctx, name, opts)
}

// Deprecated: Use Delete instead.
func (obj *Pods) Kill(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.Delete(name, namespace, opts)
}

func (obj *Pods) Get(name, namespace string) (k8sTypes.Pod, error) {
	pods, err := obj.List(namespace)
	if err != nil {
		return k8sTypes.Pod{}, err
	}
	for _, pod := range pods {
		if pod.Name == name {
			return pod, nil
		}
	}
	return k8sTypes.Pod{}, errors.New(name + " pod not found")
}

func (obj *Pods) IsTerminating(name, namespace string) (bool, error) {
	pod, err := obj.Get(name, namespace)
	if err != nil {
		return false, err
	}
	return (pod.ObjectMeta.DeletionTimestamp != nil), nil
}
