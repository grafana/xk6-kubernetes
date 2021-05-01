package kubernetes

import (
	"context"
	"fmt"

	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type JobsNamespace struct {
	client      *kubernetes.Clientset
	metaOptions metav1.ListOptions
	ctx         context.Context
}

func (obj *JobsNamespace) List(namespace string) ([]v1.Job, error) {
	result, err := obj.client.BatchV1().Jobs(namespace).List(obj.ctx, obj.metaOptions)
	return result.Items, err
}

func (obj *JobsNamespace) Get(name, namespace string) (v1.Job, error) {
	result, err := obj.client.BatchV1().Jobs(namespace).Get(obj.ctx, name, metav1.GetOptions{})
	return *result, err
}

func (obj *JobsNamespace) Create(namespace string, newJob v1.Job) (v1.Job, error) {
	fmt.Printf("%+v\n", newJob)
	job, err := obj.client.BatchV1().Jobs(namespace).Create(obj.ctx, &newJob, metav1.CreateOptions{})
	return *job, err
}
