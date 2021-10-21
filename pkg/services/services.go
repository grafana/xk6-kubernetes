package services

import (
	"context"

	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func New(client *kubernetes.Clientset, metaOptions metav1.ListOptions, ctx context.Context) *Services {
	return &Services{
		client,
		metaOptions,
		ctx,
	}
}

type Services struct {
	client      *kubernetes.Clientset
	metaOptions metav1.ListOptions
	ctx         context.Context
}

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

func (obj *Services) List(namespace string) ([]k8sTypes.Service, error) {
	svcs, err := obj.client.CoreV1().Services(namespace).List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.Service{}, err
	}
	return svcs.Items, nil
}

func (obj *Services) Kill(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.client.CoreV1().Services(namespace).Delete(obj.ctx, name, opts)
}

func (obj *Services) Get(name, namespace string, opts metav1.GetOptions) (k8sTypes.Service, error) {
	svc, err := obj.client.CoreV1().Services(namespace).Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.Service{}, err
	}

	return *svc, nil
}
