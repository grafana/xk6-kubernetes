// Package api implements helper functions for manipulating resources in a
// Kubernetes cluster.
package api

import (
	"context"

	"github.com/grafana/xk6-kubernetes/pkg/helpers"
	"github.com/grafana/xk6-kubernetes/pkg/resources"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

// Kubernetes defines an interface that extends kubernetes interface[k8s.io/client-go/kubernetes.Interface] adding
// generic functions that operate on any kind of object
type Kubernetes interface {
	resources.UnstructuredOperations
	// Helpers returns helpers for the given namespace. If none is specified, "default" is used
	Helpers(namespace string) helpers.Helpers
}

// KubernetesConfig defines the configuration for creating a Kubernetes instance
type KubernetesConfig struct {
	// Context for executing kubernetes operations
	Context context.Context
	// kubernetes rest config
	Config *rest.Config
	// Client is a pre-configured dynamic client. If provided, the rest config is not used
	Client dynamic.Interface
}

// kubernetes holds references to implementation of the Kubernetes interface
type kubernetes struct {
	ctx context.Context
	*resources.Client
}

// NewFromConfig returns a Kubernetes instance
func NewFromConfig(c KubernetesConfig) (Kubernetes, error) {
	ctx := c.Context
	if ctx == nil {
		ctx = context.TODO()
	}

	var client *resources.Client
	var err error
	if c.Client != nil {
		client = resources.NewFromClient(ctx, c.Client)
	} else {
		client, err = resources.NewFromConfig(ctx, c.Config)
		if err != nil {
			return nil, err
		}
	}

	return &kubernetes{
		ctx:    ctx,
		Client: client,
	}, nil
}

func (k *kubernetes) Helpers(namespace string) helpers.Helpers {
	if namespace == "" {
		namespace = "default"
	}
	return helpers.NewHelper(
		k.ctx,
		k.Client,
		namespace,
	)
}
