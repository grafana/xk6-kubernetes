package namespaces

import (
	"context"

	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func New(client *kubernetes.Clientset, metaOptions metav1.ListOptions, ctx context.Context) *Namespaces {
	return &Namespaces{
		client,
		metaOptions,
		ctx,
	}
}

type Namespaces struct {
	client      *kubernetes.Clientset
	metaOptions metav1.ListOptions
	ctx         context.Context
}

func (obj *Namespaces) Create(
	namespace k8sTypes.Namespace,
	opts metav1.CreateOptions,
) (k8sTypes.Namespace, error) {
	ns, err := obj.client.CoreV1().Namespaces().Create(obj.ctx, &namespace, opts)
	if err != nil {
		return k8sTypes.Namespace{}, err
	}
	return *ns, nil
}

func (obj *Namespaces) List() ([]k8sTypes.Namespace, error) {
	nss, err := obj.client.CoreV1().Namespaces().List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.Namespace{}, err
	}
	return nss.Items, nil
}

func (obj *Namespaces) Kill(name string, opts metav1.DeleteOptions) error {
	return obj.client.CoreV1().Namespaces().Delete(obj.ctx, name, opts)
}

func (obj *Namespaces) Get(name string, opts metav1.GetOptions) (k8sTypes.Namespace, error) {
	ns, err := obj.client.CoreV1().Namespaces().Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.Namespace{}, err
	}

	return *ns, nil
}
