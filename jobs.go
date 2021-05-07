package kubernetes

import (
	"context"

	v1 "k8s.io/api/batch/v1"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type JobsNamespace struct {
	client      *kubernetes.Clientset
	metaOptions metav1.ListOptions
	ctx         context.Context
}

type JobOptions struct {
	Namespace     string
	Name          string
	Image         string
	Command       []string
	RestartPolicy coreV1.RestartPolicy
}

func (obj *JobsNamespace) List(namespace string) ([]v1.Job, error) {
	result, err := obj.client.BatchV1().Jobs(namespace).List(obj.ctx, obj.metaOptions)
	return result.Items, err
}

func (obj *JobsNamespace) Get(name, namespace string) (v1.Job, error) {
	result, err := obj.client.BatchV1().Jobs(namespace).Get(obj.ctx, name, metav1.GetOptions{})
	return *result, err
}

func (obj *JobsNamespace) Kill(name, namespace string) error {
	err := obj.client.BatchV1().Jobs(namespace).Delete(obj.ctx, name, metav1.DeleteOptions{})
	return err
}

func (obj *JobsNamespace) Create(options JobOptions) (v1.Job, error) {
	container := coreV1.Container{
		Name:    options.Name,
		Image:   options.Image,
		Command: options.Command,
	}

	containers := []coreV1.Container{
		container,
	}

	var restartPolicy coreV1.RestartPolicy = "Never"

	if options.RestartPolicy != "" {
		restartPolicy = options.RestartPolicy
	}

	newJob := v1.Job{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Name: options.Name},
		Spec: v1.JobSpec{
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: coreV1.PodSpec{
					Containers:    containers,
					RestartPolicy: restartPolicy,
				},
			},
		},
	}

	job, err := obj.client.BatchV1().Jobs(options.Namespace).Create(obj.ctx, &newJob, metav1.CreateOptions{})
	return *job, err
}
