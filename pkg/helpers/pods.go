package helpers

import (
	"fmt"
	"time"

	"github.com/grafana/xk6-kubernetes/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

// PodHelper defines helper functions for manipulating Pods
type PodHelper interface {
	// WaitPodRunning waits for the Pod to be running for up to given timeout (in seconds) and returns
	// a boolean indicating if the status was reached. If the pod is Failed returns error.
	WaitPodRunning(name string, timeout uint) (bool, error)
}

func (h *helpers) WaitPodRunning(name string, timeout uint) (bool, error) {
	return utils.Retry(time.Duration(timeout)*time.Second, time.Second, func() (bool, error) {
		pod := &corev1.Pod{}
		err := h.client.Structured().Get("Pod", name, h.namespace, pod)
		if errors.IsNotFound(err) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		if pod.Status.Phase == corev1.PodFailed {
			return false, fmt.Errorf("pod has failed")
		}
		if pod.Status.Phase == corev1.PodRunning {
			return true, nil
		}
		return false, nil
	})
}
