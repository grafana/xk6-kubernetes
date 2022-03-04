package persistentvolumes

import (
	"context"
	"errors"

	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/kubernetes"
)

func New(client *kubernetes.Clientset, metaOptions metav1.ListOptions, ctx context.Context) *PersistentVolumes {
	return &PersistentVolumes{
		client,
		metaOptions,
		ctx,
	}
}

type PersistentVolumes struct {
	client      *kubernetes.Clientset
	metaOptions metav1.ListOptions
	ctx         context.Context
}

func (obj *PersistentVolumes) Apply(yaml string) (k8sTypes.PersistentVolume, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	yamlobj, _, err := decode([]byte(yaml), nil, nil)
	persistentvolume := k8sTypes.PersistentVolume{}

	if err != nil {
		return persistentvolume, err;
	}

	switch yamlobj.(type) {
	case *k8sTypes.PersistentVolume:
		persistentvolume = *yamlobj.(*k8sTypes.PersistentVolume)
	default:
		return persistentvolume, errors.New("Yaml was not a PersistentVolume")
	}

	pv, err := obj.client.CoreV1().PersistentVolumes().Create(obj.ctx, &persistentvolume, metav1.CreateOptions{})
	return *pv, err
}

func (obj *PersistentVolumes) Create(
	persistentvolume k8sTypes.PersistentVolume,
	opts metav1.CreateOptions,
) (k8sTypes.PersistentVolume, error) {
	pv, err := obj.client.CoreV1().PersistentVolumes().Create(obj.ctx, &persistentvolume, opts)
	if err != nil {
		return k8sTypes.PersistentVolume{}, err
	}
	return *pv, nil
}

func (obj *PersistentVolumes) List() ([]k8sTypes.PersistentVolume, error) {
	pvs, err := obj.client.CoreV1().PersistentVolumes().List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.PersistentVolume{}, err
	}
	return pvs.Items, nil
}

func (obj *PersistentVolumes) Delete(name string, opts metav1.DeleteOptions) error {
	return obj.client.CoreV1().PersistentVolumes().Delete(obj.ctx, name, opts)
}

func (obj *PersistentVolumes) Get(name string, opts metav1.GetOptions) (k8sTypes.PersistentVolume, error) {
	pv, err := obj.client.CoreV1().PersistentVolumes().Get(obj.ctx, name, opts)
	if err != nil {
		return k8sTypes.PersistentVolume{}, err
	}

	return *pv, nil
}
