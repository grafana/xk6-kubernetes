package ingresses

import (
	"context"

	k8sTypes "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func New(client *kubernetes.Clientset, metaOptions metav1.ListOptions, ctx context.Context) *Ingresses {
	return &Ingresses{
		client,
		metaOptions,
		ctx,
	}
}

type Ingresses struct {
	client      *kubernetes.Clientset
	metaOptions metav1.ListOptions
	ctx         context.Context
}

func (obj *Ingresses) Create(
	ingress k8sTypes.Ingress,
	namespace string,
	opts metav1.CreateOptions,
) (k8sTypes.Ingress, error) {
	ing, err := obj.client.NetworkingV1().Ingresses(namespace).Create(obj.ctx, &ingress, opts)
	if err != nil {
		return k8sTypes.Ingress{}, err
	}
	return *ing, nil
}

func (obj *Ingresses) List(namespace string) ([]k8sTypes.Ingress, error) {
	ings, err := obj.client.NetworkingV1().Ingresses(namespace).List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.Ingress{}, err
	}
	return ings.Items, nil
}

func (obj *Ingresses) Kill(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.client.NetworkingV1().Ingresses(namespace).Delete(obj.ctx, name, opts)
}

func (obj *Ingresses) Get(name, namespace string, opts metav1.GetOptions) (k8sTypes.Ingress, error) {
	ing, err := obj.client.NetworkingV1().Ingresses(namespace).Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.Ingress{}, err
	}

	return *ing, nil
}
