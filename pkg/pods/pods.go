package pods

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func New(client kubernetes.Interface, config *rest.Config, metaOptions metav1.ListOptions, ctx context.Context) *Pods {
	return &Pods{
		client,
		config,
		metaOptions,
		ctx,
	}
}

// ExecOptions describe the command to be executed and the target container
type ExecOptions struct {
	Namespace string   // namespace where the pod is running
	Pod       string   // name of the Pod to execute the command in
	Container string   // name of the container to execute the command in
	Command   []string // command to be exectued with its parameters
	Stdin     []byte   // stdin to be supplied to the command
}

// ExecResult contains the output obtained from the execution of a command
type ExecResult struct {
	Stdout []byte
	Stderr []byte
}

// ContainerOptions describes a container to be started in a pod
type ContainerOptions struct {
	Name         string   // name of the container
	Image        string   // image to be attached
	Command      []string // command to be executed by the container
	Capabilities []string // capabilities to be added to the container's security context
}

type Pods struct {
	client      kubernetes.Interface
	config      *rest.Config
	metaOptions metav1.ListOptions
	ctx         context.Context
}

// PodOptions describe a Pod to be executed
type PodOptions struct {
	Namespace     string                 // namespace where the pod will be executed
	Name          string                 // name of the pod
	Image         string                 // image to be executed by the pod's container
	Command       []string               // command to be executed by the pod's container and its arguments
	RestartPolicy k8sTypes.RestartPolicy // policy for restarting containers in the pod. One of One of Always, OnFailure, Never
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

// Create runs a pod specified by the options
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
	if err != nil {
		return k8sTypes.Pod{}, err
	}
	return *pod, nil
}

// Exec executes a non-interactive command described in options and returns the stdout and stderr outputs
func (obj *Pods) Exec(options ExecOptions) (*ExecResult, error) {

	req := obj.client.CoreV1().RESTClient().
		Post().
		Namespace(options.Namespace).
		Resource("pods").
		Name(options.Pod).
		SubResource("exec").
		VersionedParams(&k8sTypes.PodExecOptions{
			Container: options.Container,
			Command:   options.Command,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(obj.config, "POST", req.URL())
	if err != nil {
		return nil, err
	}

	// connect to the command
	var stdout, stderr bytes.Buffer
	stdin := bytes.NewReader(options.Stdin)
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    true,
	})

	if err != nil {
		return nil, err
	}

	result := ExecResult{
		stdout.Bytes(),
		stderr.Bytes(),
	}

	return &result, nil
}

// AddEphemeralContainer adds an ephemeral container to a running pod. The Pod is identified by name and namespace.
// The container is described by options
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
