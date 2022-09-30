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
	resources.Operations
	Helpers() helpers.Helpers
	NamespacedHelpers(namespace string) helpers.Helpers
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

// Helpers returns Helpers for the default namespace
func (k *kubernetes) Helpers() helpers.Helpers {
	return helpers.NewHelper(
		k.ctx,
		k.Client,
		"default",
	)
}

// NamespacedHelpers returns helpers for the given namespace
func (k *kubernetes) NamespacedHelpers(namespace string) helpers.Helpers {
	return helpers.NewHelper(
		k.ctx,
		k.Client,
		namespace,
	)
}
