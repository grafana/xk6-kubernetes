package helpers

import (
	"errors"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// PodHelper defines helper functions for manipulating Pods
type PodHelper interface {
	// WaitPodRunning waits for the Pod to be running for up to given timeout (in seconds) and returns
	// a boolean indicating if the status was reached. If the pod is Failed returns error.
	WaitPodRunning(name string, timeout uint) (bool, error)
}

// podConditionChecker defines a function that checks if a pod satisfies a condition
type podConditionChecker func(*corev1.Pod) (bool, error)

// waitForCondition watches a Pod in a namespace until a podConditionChecker is satisfied or a timeout expires
func (h *helpers) waitForCondition(
	name string,
	timeout time.Duration,
	checker podConditionChecker,
) (bool, error) {
	options := metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	}
	watcher, err := h.client.Structured().Watch(
		"Pod",
		h.namespace,
		options,
	)
	if err != nil {
		return false, err
	}
	defer watcher.Stop()

	expired := time.After(timeout)
	for {
		select {
		case <-expired:
			return false, nil
		case event := <-watcher.ResultChan():
			if event.Type == watch.Error {
				return false, fmt.Errorf("error watching for pod: %v", event.Object)
			}
			if event.Type == watch.Modified {
				pod, isPod := event.Object.(*corev1.Pod)
				if !isPod {
					return false, errors.New("received unknown object while watching for pods")
				}
				condition, err := checker(pod)
				if condition || err != nil {
					return condition, err
				}
			}
		}
	}
}

func (h *helpers) WaitPodRunning(name string, timeout uint) (bool, error) {
	return h.waitForCondition(
		name,
		time.Duration(timeout)*time.Second,
		func(pod *corev1.Pod) (bool, error) {
			if pod.Status.Phase == corev1.PodFailed {
				return false, errors.New("pod has failed")
			}
			if pod.Status.Phase == corev1.PodRunning {
				return true, nil
			}
			return false, nil
		},
	)
}
