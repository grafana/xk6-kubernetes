package persistentvolumeclaims

import (
	"context"
	"errors"

	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

func New(client *kubernetes.Clientset, metaOptions metav1.ListOptions, ctx context.Context) *PersistentVolumeClaims {
	return &PersistentVolumeClaims{
		client,
		metaOptions,
		ctx,
	}
}

type PersistentVolumeClaims struct {
	client      *kubernetes.Clientset
	metaOptions metav1.ListOptions
	ctx         context.Context
}

func (obj *PersistentVolumeClaims) Apply(yaml string, namespace string) (k8sTypes.PersistentVolumeClaim, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	yamlobj, _, err := decode([]byte(yaml), nil, nil)
	persistentvolumeclaim := k8sTypes.PersistentVolumeClaim{}

	if err != nil {
		return persistentvolumeclaim, err
	}

	switch yamlobj.(type) {
	case *k8sTypes.PersistentVolumeClaim:
		persistentvolumeclaim = *yamlobj.(*k8sTypes.PersistentVolumeClaim)
	default:
		return persistentvolumeclaim, errors.New("Yaml was not a PersistentVolumeClaim")
	}

	pvc, err := obj.client.CoreV1().PersistentVolumeClaims(namespace).Create(obj.ctx, &persistentvolumeclaim, metav1.CreateOptions{})
	return *pvc, err
}

func (obj *PersistentVolumeClaims) Create(
	persistentvolumeclaim k8sTypes.PersistentVolumeClaim,
	namespace string,
	opts metav1.CreateOptions,
) (k8sTypes.PersistentVolumeClaim, error) {
	pvc, err := obj.client.CoreV1().PersistentVolumeClaims(namespace).Create(obj.ctx, &persistentvolumeclaim, opts)
	if err != nil {
		return k8sTypes.PersistentVolumeClaim{}, err
	}
	return *pvc, nil
}

func (obj *PersistentVolumeClaims) List(namespace string) ([]k8sTypes.PersistentVolumeClaim, error) {
	pvcs, err := obj.client.CoreV1().PersistentVolumeClaims(namespace).List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.PersistentVolumeClaim{}, err
	}
	return pvcs.Items, nil
}

func (obj *PersistentVolumeClaims) Delete(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.client.CoreV1().PersistentVolumeClaims(namespace).Delete(obj.ctx, name, opts)
}

func (obj *PersistentVolumeClaims) Get(name, namespace string, opts metav1.GetOptions) (k8sTypes.PersistentVolumeClaim, error) {
	pvc, err := obj.client.CoreV1().PersistentVolumeClaims(namespace).Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.PersistentVolumeClaim{}, err
	}

	return *pvc, nil
}
