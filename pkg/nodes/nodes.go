package nodes

import (
	"context"

	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func New(client kubernetes.Interface, metaOptions metav1.ListOptions, ctx context.Context) *Nodes {
	return &Nodes{
		client,
		metaOptions,
		ctx,
	}
}

type Nodes struct {
	client      kubernetes.Interface
	metaOptions metav1.ListOptions
	ctx         context.Context
}

func (obj *Nodes) List() ([]k8sTypes.Node, error) {
	nodes, err := obj.client.CoreV1().Nodes().List(obj.ctx, obj.metaOptions)
	if err != nil {
		return []k8sTypes.Node{}, err
	}
	return nodes.Items, nil
}
