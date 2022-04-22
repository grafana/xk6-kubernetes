package jobs

import (
	"context"
	"errors"

	v1 "k8s.io/api/batch/v1"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

func New(client *kubernetes.Clientset, metaOptions metav1.ListOptions, ctx context.Context) *Jobs {
	return &Jobs{
		client,
		metaOptions,
		ctx,
	}
}

type Jobs struct {
	client      *kubernetes.Clientset
	metaOptions metav1.ListOptions
	ctx         context.Context
}

type JobOptions struct {
	Namespace     string
	Name          string
	NodeName      string
	Image         string
	Command       []string
	RestartPolicy coreV1.RestartPolicy
}

func (obj *Jobs) List(namespace string) ([]v1.Job, error) {
	result, err := obj.client.BatchV1().Jobs(namespace).List(obj.ctx, obj.metaOptions)
	return result.Items, err
}

func (obj *Jobs) Get(name, namespace string) (v1.Job, error) {
	result, err := obj.client.BatchV1().Jobs(namespace).Get(obj.ctx, name, metav1.GetOptions{})
	return *result, err
}

func (obj *Jobs) Delete(name, namespace string) error {
	return obj.client.BatchV1().Jobs(namespace).Delete(obj.ctx, name, metav1.DeleteOptions{})
}

// Deprecated: Use Delete instead.
func (obj *Jobs) Kill(name, namespace string) error {
	return obj.Delete(name, namespace)
}

func (obj *Jobs) Apply(yaml string, namespace string) (v1.Job, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	yamlobj, _, err := decode([]byte(yaml), nil, nil)
	job := v1.Job{}

	if err != nil {
		return job, err
	}

	switch yamlobj.(type) {
	case *v1.Job:
		job = *yamlobj.(*v1.Job)
	default:
		return job, errors.New("Yaml was not a Job")
	}

	jb, err := obj.client.BatchV1().Jobs(namespace).Create(obj.ctx, &job, metav1.CreateOptions{})
	return *jb, err
}

func (obj *Jobs) Create(options JobOptions) (v1.Job, error) {
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
					NodeName:      options.NodeName,
					Containers:    containers,
					RestartPolicy: restartPolicy,
				},
			},
		},
	}

	job, err := obj.client.BatchV1().Jobs(options.Namespace).Create(obj.ctx, &newJob, metav1.CreateOptions{})
	return *job, err
}
