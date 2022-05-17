package jobs

import (
	"context"
	"errors"
	"fmt"
	"time"

	v1 "k8s.io/api/batch/v1"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

func New(client kubernetes.Interface, metaOptions metav1.ListOptions, ctx context.Context) *Jobs {
	return &Jobs{
		client,
		metaOptions,
		ctx,
	}
}

type Jobs struct {
	client      kubernetes.Interface
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
	Wait          string
}

func (obj *Jobs) List(namespace string) ([]v1.Job, error) {
	result, err := obj.client.BatchV1().Jobs(namespace).List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []v1.Job{}, err
	}
	return result.Items, nil
}

func (obj *Jobs) Get(name, namespace string) (v1.Job, error) {
	result, err := obj.client.BatchV1().Jobs(namespace).Get(obj.ctx, name, metav1.GetOptions{})
	if err != nil {
		return v1.Job{}, err
	}
	return *result, nil
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
	if err != nil {
		return v1.Job{}, err
	}
	return *jb, nil
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
	if err != nil {
		return v1.Job{}, err
	}
	if options.Wait == "" {
		return *job, nil
	}
	waitOpts := WaitOptions{
		Name:      options.Name,
		Namespace: options.Namespace,
		Timeout:   options.Wait,
	}
	status, err := obj.Wait(waitOpts)
	if err != nil {
		return v1.Job{}, err
	}
	if !status {
		return v1.Job{}, errors.New("timeout exceeded waiting for job to complete")
	}
	return obj.Get(options.Name, options.Namespace)
}

// WaitOptions specify the options for waiting for a Job to complete
type WaitOptions struct {
	Name      string
	Namespace string
	Timeout   string
}

// isCompleted returns if the job is completed or not. Returns an error if the job is failed.
func isCompleted(job *v1.Job) (bool, error) {
	for _, condition := range job.Status.Conditions {
		if condition.Type == v1.JobFailed && condition.Status == coreV1.ConditionTrue {
			return false, errors.New("Job failed with reason: " + condition.Reason)
		}
		if condition.Type == v1.JobComplete && condition.Status == coreV1.ConditionTrue {
			return true, nil
		}
	}
	return false, nil
}

// Wait for all pods to complete
func (obj *Jobs) Wait(options WaitOptions) (bool, error) {
	// wait for updates until completion
	timeout, err := time.ParseDuration(options.Timeout)
	if err != nil {
		return false, err
	}
	selector := fields.Set{
		"metadata.name": options.Name,
	}.AsSelector()
	watcher, err := obj.client.BatchV1().Jobs(options.Namespace).Watch(
		obj.ctx,
		metav1.ListOptions{
			FieldSelector: selector.String(),
		},
	)
	if err != nil {
		return false, err
	}
	defer watcher.Stop()

	for {
		select {
		case <-time.After(timeout):
			return false, nil
		case event := <-watcher.ResultChan():
			if event.Type == watch.Error {
				return false, fmt.Errorf("error watching for job: %v", event.Object)
			}
			if event.Type == watch.Modified {
				job, isJob := event.Object.(*v1.Job)
				if !isJob {
					return false, errors.New("received unknown object while watching for jobs")
				}
				completed, err := isCompleted(job)
				if err != nil {
					return false, err
				}
				if completed {
					return true, nil
				}
			}
		}
	}
	return false, nil
}
