package namespaces

import (
	"context"
	"errors"

	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

func New(client kubernetes.Interface, metaOptions metav1.ListOptions, ctx context.Context) *Namespaces {
	return &Namespaces{
		client,
		metaOptions,
		ctx,
	}
}

type Namespaces struct {
	client      kubernetes.Interface
	metaOptions metav1.ListOptions
	ctx         context.Context
}

func (obj *Namespaces) Apply(yaml string) (k8sTypes.Namespace, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	yamlobj, _, err := decode([]byte(yaml), nil, nil)
	namespace := k8sTypes.Namespace{}

	if err != nil {
		return namespace, err
	}

	switch yamlobj.(type) {
	case *k8sTypes.Namespace:
		namespace = *yamlobj.(*k8sTypes.Namespace)
	default:
		return namespace, errors.New("Yaml was not a Namespace")
	}

	ns, err := obj.client.CoreV1().Namespaces().Create(obj.ctx, &namespace, metav1.CreateOptions{})
	return *ns, err
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

func (obj *Namespaces) Delete(name string, opts metav1.DeleteOptions) error {
	return obj.client.CoreV1().Namespaces().Delete(obj.ctx, name, opts)
}

// Deprecated: Use Delete instead.
func (obj *Namespaces) Kill(name string, opts metav1.DeleteOptions) error {
	return obj.Delete(name, opts)
}

func (obj *Namespaces) Get(name string, opts metav1.GetOptions) (k8sTypes.Namespace, error) {
	ns, err := obj.client.CoreV1().Namespaces().Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.Namespace{}, err
	}

	return *ns, nil
}
