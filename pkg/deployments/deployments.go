package deployments

import (
	"context"
	"errors"

	k8sTypes "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

func New(client *kubernetes.Clientset, metaOptions metav1.ListOptions, ctx context.Context) *Deployments {
	return &Deployments{
		client,
		metaOptions,
		ctx,
	}
}

type Deployments struct {
	client      *kubernetes.Clientset
	metaOptions metav1.ListOptions
	ctx         context.Context
}

func (obj *Deployments) Apply(yaml string, namespace string) (k8sTypes.Deployment, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	yamlobj, _, err := decode([]byte(yaml), nil, nil)
	deployment := k8sTypes.Deployment{}

	if err != nil {
		return deployment, err
	}

	switch yamlobj.(type) {
	case *k8sTypes.Deployment:
		deployment = *yamlobj.(*k8sTypes.Deployment)
	default:
		return deployment, errors.New("Yaml was not a Deployment")
	}

	dep, err := obj.client.AppsV1().Deployments(namespace).Create(obj.ctx, &deployment, metav1.CreateOptions{})
	return *dep, err
}

func (obj *Deployments) Create(
	deployment k8sTypes.Deployment,
	namespace string,
	opts metav1.CreateOptions,
) (k8sTypes.Deployment, error) {
	dep, err := obj.client.AppsV1().Deployments(namespace).Create(obj.ctx, &deployment, opts)
	if err != nil {
		return k8sTypes.Deployment{}, err
	}
	return *dep, nil
}

func (obj *Deployments) List(namespace string) ([]k8sTypes.Deployment, error) {
	dps, err := obj.client.AppsV1().Deployments(namespace).List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.Deployment{}, err
	}
	return dps.Items, nil
}

func (obj *Deployments) Delete(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.client.AppsV1().Deployments(namespace).Delete(obj.ctx, name, opts)
}

// Deprecated: Use Delete instead.
func (obj *Deployments) Kill(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.Delete(name, namespace, opts)
}

func (obj *Deployments) Get(name, namespace string, opts metav1.GetOptions) (k8sTypes.Deployment, error) {
	dp, err := obj.client.AppsV1().Deployments(namespace).Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.Deployment{}, err
	}
	return *dp, nil
}
