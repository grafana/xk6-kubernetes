// Package nodes provides implementation of Node resources for Kubernetes
//
// Deprecated: Use the resources package instead.
package nodes

import (
	"context"

	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// New creates a new instance backed by the provided client
//
// Deprecated: No longer used.
func New(ctx context.Context, client kubernetes.Interface, metaOptions metav1.ListOptions) *Nodes {
	return &Nodes{
		client,
		metaOptions,
		ctx,
	}
}

// Nodes provides API for manipulating Node resources within a Kubernetes cluster
//
// Deprecated: No longer used in favor of generic resources.
type Nodes struct {
	client      kubernetes.Interface
	metaOptions metav1.ListOptions
	ctx         context.Context
}

// List returns a collection of Nodes comprising the cluster
//
// Deprecated: Use resources.List instead.
func (obj *Nodes) List() ([]k8sTypes.Node, error) {
	nodes, err := obj.client.CoreV1().Nodes().List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.Node{}, err
	}
	return nodes.Items, nil
}
