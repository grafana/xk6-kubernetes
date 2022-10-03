package helpers

import (
	"fmt"
	"time"

	"github.com/grafana/xk6-kubernetes/pkg/utils"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

// ServiceHelper implements functions for dealing with services
type ServiceHelper interface {
	// WaitServiceReady waits for the given service to have at least one endpoint available
	// or the timeout (in seconds) expires
	WaitServiceReady(service string, timeout uint) error
}

func (h *helpers) WaitServiceReady(service string, timeout uint) error {
	return utils.Retry(time.Duration(timeout)*time.Second, time.Second, func() (bool, error) {
		ep := &corev1.Endpoints{}
		err := h.client.Structured().Get("Endpoint", service, h.namespace, ep)
		if err != nil {
			if errors.IsNotFound(err) {
				return false, nil
			}
			return false, fmt.Errorf("failed to access service: %w", err)
		}

		for _, subset := range ep.Subsets {
			if len(subset.Addresses) > 0 {
				return true, nil
			}
		}

		return false, nil
	})
}
