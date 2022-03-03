package ingresses

import (
	"context"
	"errors"

	k8sTypes "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
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

func (obj *Ingresses) Apply(yaml string, namespace string) (k8sTypes.Ingress, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	yamlobj, _, err := decode([]byte(yaml), nil, nil)
	ingress := k8sTypes.Ingress{}

	if err != nil {
		return ingress, err
	}

	switch yamlobj.(type) {
	case *k8sTypes.Ingress:
		ingress = *yamlobj.(*k8sTypes.Ingress)
	default:
		return ingress, errors.New("Yaml was not an Ingress")
	}

	ing, err := obj.client.NetworkingV1().Ingresses(namespace).Create(obj.ctx, &ingress, metav1.CreateOptions{})
	return *ing, err
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

func (obj *Ingresses) Delete(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.client.NetworkingV1().Ingresses(namespace).Delete(obj.ctx, name, opts)
}

// Deprecated: Use Delete instead.
func (obj *Ingresses) Kill(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.Delete(name, namespace, opts)
}

func (obj *Ingresses) Get(name, namespace string, opts metav1.GetOptions) (k8sTypes.Ingress, error) {
	ing, err := obj.client.NetworkingV1().Ingresses(namespace).Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.Ingress{}, err
	}

	return *ing, nil
}
