// Package helpers offers functions to simplify dealing with kubernetes resources.
package helpers

import (
	"context"

	"github.com/grafana/xk6-kubernetes/pkg/resources"
)

// Helpers offers Helper functions grouped by the objects they handle
type Helpers interface {
}

// helpers struct holds the data required by the helpers
type helpers struct {
	client    *resources.Client
	ctx       context.Context
	namespace string
}

// NewHelper creates a set of helper functions on the default namespace
func NewHelper(ctx context.Context, client *resources.Client, namespace string) Helpers {
	return &helpers{
		client:    client,
		ctx:       ctx,
		namespace: namespace,
	}

}
