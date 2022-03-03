package services

import (
	"context"
	"errors"

	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
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

func (obj *Services) Apply(yaml string, namespace string) (k8sTypes.Service, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	yamlobj, _, err := decode([]byte(yaml), nil, nil)
	service := k8sTypes.Service{}

	if err != nil {
		return service, err
	}

	switch yamlobj.(type) {
	case *k8sTypes.Service:
		service = *yamlobj.(*k8sTypes.Service)
	default:
		return service, errors.New("Yaml was not a Service")
	}

	svc, err := obj.client.CoreV1().Services(namespace).Create(obj.ctx, &service, metav1.CreateOptions{})
	return *svc, err
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

func (obj *Services) Delete(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.client.CoreV1().Services(namespace).Delete(obj.ctx, name, opts)
}

// Deprecated: Use Delete instead.
func (obj *Services) Kill(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.Delete(name, namespace, opts)
}

func (obj *Services) Get(name, namespace string, opts metav1.GetOptions) (k8sTypes.Service, error) {
	svc, err := obj.client.CoreV1().Services(namespace).Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.Service{}, err
	}

	return *svc, nil
}
