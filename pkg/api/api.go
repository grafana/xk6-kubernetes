// Package kubernetes implements helper functions for manipulating resources in a
// Kubernetes cluster.
package api

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

// maps kinds to api resources
// TODO: complete with most common kinds
var knownKinds = map[string]schema.GroupVersionResource{
	"ConfigMap":  {Group: "", Version: "v1", Resource: "configmaps"},
	"Deployment": {Group: "apps", Version: "v1", Resource: "deployments"},
	"Job":        {Group: "batch", Version: "v1", Resource: "jobs"},
	"Pod":        {Group: "", Version: "v1", Resource: "pods"},
	"Namespace":  {Group: "", Version: "v1", Resource: "namespaces"},
	"Node":       {Group: "", Version: "v1", Resource: "nodes"},
	"Secret":     {Group: "", Version: "v1", Resource: "secrets"},
	"Service":    {Group: "", Version: "v1", Resource: "services"},
}

// Defines an interface that extends kubernetes interface[k8s.io/client-go/kubernetes.Interface] adding
// generic functions that operate on any kind of object
type Kubernetes interface {
	Apply(manifest string) error
	Create(obj map[string]interface{}) (map[string]interface{}, error)
	Get(kind string, name string, namespace string) (map[string]interface{}, error)
	List(kind string, namespace string) ([]map[string]interface{}, error)
	Delete(kind string, name string, namespace string) error
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

// kubernetes Holds the reference to the helpers for interacting with kubernetes
type kubernetes struct {
	ctx        context.Context
	client     dynamic.Interface
	serializer runtime.Serializer
}

// NewFromConfig returns a Kubernetes instance
func NewFromConfig(c KubernetesConfig) (Kubernetes, error) {
	client := c.Client
	var err error
	if client == nil {
		client, err = dynamic.NewForConfig(c.Config)
		if err != nil {
			return nil, err
		}
	}

	ctx := c.Context
	if ctx == nil {
		ctx = context.TODO()
	}

	return &kubernetes{
		ctx:        ctx,
		client:     client,
		serializer: yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme),
	}, nil
}

// Apply creates a resource in a kubernetes cluster from a yaml manifest
func (k *kubernetes) Apply(manifest string) error {
	uObj := &unstructured.Unstructured{}
	_, gvk, err := k.serializer.Decode([]byte(manifest), nil, uObj)
	if err != nil {
		return err
	}
	resource, known := knownKinds[gvk.Kind]
	if !known {
		return fmt.Errorf("unknown kind: '%s'", gvk.Kind)
	}

	namespace := uObj.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}

	_, err = k.client.Resource(resource).
		Namespace(namespace).
		Create(
			k.ctx,
			uObj,
			metav1.CreateOptions{},
		)

	return err
}

// Create creates a resource in a kubernetes cluster from a yaml manifest
func (k *kubernetes) Create(obj map[string]interface{}) (map[string]interface{}, error) {
	uObj := &unstructured.Unstructured{
		Object: obj,
	}

	gvk := uObj.GroupVersionKind()
	namespace := uObj.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}
	resource, known := knownKinds[gvk.Kind]
	if !known {
		return nil, fmt.Errorf("unknown kind: '%s'", gvk.Kind)
	}

	resp, err := k.client.Resource(resource).
		Namespace(namespace).
		Create(
			k.ctx,
			uObj,
			metav1.CreateOptions{},
		)

	if err != nil {
		return nil, err
	}
	return resp.UnstructuredContent(), nil
}

// Get returns an object given its kind, name and namespace
func (k *kubernetes) Get(kind string, name string, namespace string) (map[string]interface{}, error) {
	resource, known := knownKinds[kind]
	if !known {
		return nil, fmt.Errorf("unknown kind: '%s'", kind)
	}

	resp, err := k.client.
		Resource(resource).
		Namespace(namespace).
		Get(
			k.ctx,
			name,
			metav1.GetOptions{},
		)

	if err != nil {
		return nil, err
	}
	return resp.UnstructuredContent(), nil
}

// List returns a list of objects given its kind and namespace
func (k *kubernetes) List(kind string, namespace string) ([]map[string]interface{}, error) {
	resource, known := knownKinds[kind]
	if !known {
		return nil, fmt.Errorf("unknown kind: '%s'", kind)
	}

	resp, err := k.client.
		Resource(resource).
		Namespace(namespace).
		List(k.ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	var list []map[string]interface{}
	for _, uObj := range resp.Items {
		list = append(list, uObj.UnstructuredContent())
	}
	return list, nil
}

// Delete deletes an object given its kind, name and namespace
func (k *kubernetes) Delete(kind string, name string, namespace string) error {
	resource, known := knownKinds[kind]
	if !known {
		return fmt.Errorf("unknown kind: '%s'", kind)
	}

	err := k.client.
		Resource(resource).
		Namespace(namespace).
		Delete(k.ctx, name, metav1.DeleteOptions{})

	return err
}
