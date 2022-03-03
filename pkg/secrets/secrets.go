package secrets

import (
	"context"
	"errors"

	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

func New(client *kubernetes.Clientset, metaOptions metav1.ListOptions, ctx context.Context) *Secrets {
	return &Secrets{
		client,
		metaOptions,
		ctx,
	}
}

type Secrets struct {
	client      *kubernetes.Clientset
	metaOptions metav1.ListOptions
	ctx         context.Context
}

func (obj *Secrets) Apply(yaml string, namespace string) (k8sTypes.Secret, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	yamlobj, _, err := decode([]byte(yaml), nil, nil)
	secret := k8sTypes.Secret{}

	if err != nil {
		return secret, err
	}

	switch yamlobj.(type) {
	case *k8sTypes.Secret:
		secret = *yamlobj.(*k8sTypes.Secret)
	default:
		return secret, errors.New("Yaml was not a Secret")
	}

	scrt, err := obj.client.CoreV1().Secrets(namespace).Create(obj.ctx, &secret, metav1.CreateOptions{})
	return *scrt, err
}

func (obj *Secrets) Create(
	secret k8sTypes.Secret,
	namespace string,
	opts metav1.CreateOptions,
) (k8sTypes.Secret, error) {
	scrt, err := obj.client.CoreV1().Secrets(namespace).Create(obj.ctx, &secret, opts)
	if err != nil {
		return k8sTypes.Secret{}, err
	}
	return *scrt, nil
}

func (obj *Secrets) List(namespace string) ([]k8sTypes.Secret, error) {
	scrts, err := obj.client.CoreV1().Secrets(namespace).List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.Secret{}, err
	}
	return scrts.Items, nil
}

func (obj *Secrets) Delete(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.client.CoreV1().Secrets(namespace).Delete(obj.ctx, name, opts)
}

// Deprecated: Use Delete instead.
func (obj *Secrets) Kill(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.Delete(name, namespace, opts)
}

func (obj *Secrets) Get(name, namespace string, opts metav1.GetOptions) (k8sTypes.Secret, error) {
	scrt, err := obj.client.CoreV1().Secrets(namespace).Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.Secret{}, err
	}

	return *scrt, nil
}
