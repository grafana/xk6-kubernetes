package pods

import (
	"context"
	"encoding/json"
	"errors"

	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"
)

func New(client *kubernetes.Clientset, metaOptions metav1.ListOptions, ctx context.Context) *Pods {
	return &Pods{
		client,
		metaOptions,
		ctx,
	}
}

type ContainerOptions struct {
	Name         string
	Image        string
	Command      []string
	Capabilities []string
}

type Pods struct {
	client      *kubernetes.Clientset
	metaOptions metav1.ListOptions
	ctx         context.Context
}

type PodOptions struct {
	Namespace     string
	Name          string
	Image         string
	Command       []string
	RestartPolicy k8sTypes.RestartPolicy
}

func (obj *Pods) List(namespace string) ([]k8sTypes.Pod, error) {
	pods, err := obj.client.CoreV1().Pods(namespace).List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.Pod{}, err
	}
	return pods.Items, nil
}

func (obj *Pods) Delete(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.client.CoreV1().Pods(namespace).Delete(obj.ctx, name, opts)
}

// Deprecated: Use Delete instead.
func (obj *Pods) Kill(name, namespace string, opts metav1.DeleteOptions) error {
	return obj.Delete(name, namespace, opts)
}

func (obj *Pods) Get(name, namespace string) (k8sTypes.Pod, error) {
	pods, err := obj.List(namespace)
	if err != nil {
		return k8sTypes.Pod{}, err
	}
	for _, pod := range pods {
		if pod.Name == name {
			return pod, nil
		}
	}
	return k8sTypes.Pod{}, errors.New(name + " pod not found")
}

func (obj *Pods) IsTerminating(name, namespace string) (bool, error) {
	pod, err := obj.Get(name, namespace)
	if err != nil {
		return false, err
	}
	return (pod.ObjectMeta.DeletionTimestamp != nil), nil
}

func (obj *Pods) Create(options PodOptions) (k8sTypes.Pod, error) {
	container := k8sTypes.Container{
		Name:    options.Name,
		Image:   options.Image,
		Command: options.Command,
	}

	containers := []k8sTypes.Container{
		container,
	}

	var restartPolicy k8sTypes.RestartPolicy = "Never"

	if options.RestartPolicy != "" {
		restartPolicy = options.RestartPolicy
	}

	newPod := k8sTypes.Pod{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Name: options.Name},
		Spec: k8sTypes.PodSpec{
			Containers:    containers,
			RestartPolicy: restartPolicy,
		},
	}

	pod, err := obj.client.CoreV1().Pods(options.Namespace).Create(obj.ctx, &newPod, metav1.CreateOptions{})
	return *pod, err
}

func (obj *Pods) AddEphemeralContainer(name, namespace string, options ContainerOptions) error {
	pod, err := obj.Get(name, namespace)
	if err != nil {
		return err
	}
	podJson, err := json.Marshal(pod)
	if err != nil {
		return err
	}
	container, err := generateEphemeralContainer(options)
	if err != nil {
		return err
	}

	updatedPod := pod.DeepCopy()
	updatedPod.Spec.EphemeralContainers = append(updatedPod.Spec.EphemeralContainers, *container)
	updateJson, err := json.Marshal(updatedPod)
	if err != nil {
		return err
	}

	patch, err := strategicpatch.CreateTwoWayMergePatch(podJson, updateJson, pod)
	if err != nil {
		return err
	}

	_, err = obj.client.CoreV1().Pods(namespace).Patch(obj.ctx, pod.Name, types.StrategicMergePatchType, patch, metav1.PatchOptions{}, "ephemeralcontainers")

	return err

}

func generateEphemeralContainer(o ContainerOptions) (*k8sTypes.EphemeralContainer, error) {
	var capabilities []k8sTypes.Capability
	for _, capability := range o.Capabilities {
		capabilities = append(capabilities, k8sTypes.Capability(capability))
	}

	return &k8sTypes.EphemeralContainer{
		EphemeralContainerCommon: k8sTypes.EphemeralContainerCommon{
			Name:    o.Name,
			Image:   o.Image,
			Command: o.Command,
			SecurityContext: &k8sTypes.SecurityContext{
				Capabilities: &k8sTypes.Capabilities{
					Add: capabilities,
				},
			},
		},
	}, nil
}
