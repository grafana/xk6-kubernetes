// Package resources implement the interface for accessing kubernetes resources
package resources

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

// Operations defines generic functions that operate on any kind of Kubernetes object
type Operations interface {
	Apply(manifest string) error
	Create(obj map[string]interface{}) (map[string]interface{}, error)
	Get(kind string, name string, namespace string) (map[string]interface{}, error)
	List(kind string, namespace string) ([]map[string]interface{}, error)
	Delete(kind string, name string, namespace string) error
}

// kubernetes holds the reference to
type Client struct {
	ctx        context.Context
	dynamic    dynamic.Interface
	serializer runtime.Serializer
}

// NewFromConfig creates a new Client using the provided kubernetes client configuration
func NewFromConfig(ctx context.Context, config *rest.Config) (*Client, error) {
	dynamic, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return NewFromClient(ctx, dynamic), nil
}

// NewFromClient creates a new client from a dynamic Kubernetes client
func NewFromClient(ctx context.Context, dynamic dynamic.Interface) *Client {
	return &Client{
		ctx:        ctx,
		dynamic:    dynamic,
		serializer: yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme),
	}
}

// maps kinds to api resources
func knownKinds(kind string) (schema.GroupVersionResource, error) {
	kindMapping := map[string]schema.GroupVersionResource{
		"ConfigMap":             {Group: "", Version: "v1", Resource: "configmaps"},
		"Deployment":            {Group: "apps", Version: "v1", Resource: "deployments"},
		"Job":                   {Group: "batch", Version: "v1", Resource: "jobs"},
		"PersistentVolume":      {Group: "", Version: "v1", Resource: "persistentvolumes"},
		"PersistentVolumeClaim": {Group: "", Version: "v1", Resource: "persistentvolumeclaims"},
		"Pod":                   {Group: "", Version: "v1", Resource: "pods"},
		"Namespace":             {Group: "", Version: "v1", Resource: "namespaces"},
		"Node":                  {Group: "", Version: "v1", Resource: "nodes"},
		"Secret":                {Group: "", Version: "v1", Resource: "secrets"},
		"Service":               {Group: "", Version: "v1", Resource: "services"},
		"StatefulSet":           {Group: "apps", Version: "v1", Resource: "statefulsets"},
	}

	gvk, found := kindMapping[kind]
	if !found {
		return schema.GroupVersionResource{}, fmt.Errorf("unknown kind: '%s'", kind)
	}
	return gvk, nil
}

// Apply creates a resource in a kubernetes cluster from a YAML manifest
func (c *Client) Apply(manifest string) error {
	uObj := &unstructured.Unstructured{}
	_, gvk, err := c.serializer.Decode([]byte(manifest), nil, uObj)
	if err != nil {
		return err
	}
	resource, err := knownKinds(gvk.Kind)
	if err != nil {
		return err
	}

	namespace := uObj.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}

	_, err = c.dynamic.Resource(resource).
		Namespace(namespace).
		Create(
			c.ctx,
			uObj,
			metav1.CreateOptions{},
		)
	return err
}

// Create creates a resource in a kubernetes cluster from an object with its specification
func (c *Client) Create(obj map[string]interface{}) (map[string]interface{}, error) {
	uObj := &unstructured.Unstructured{
		Object: obj,
	}

	gvk := uObj.GroupVersionKind()
	namespace := uObj.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}
	resource, err := knownKinds(gvk.Kind)
	if err != nil {
		return nil, err
	}

	resp, err := c.dynamic.Resource(resource).
		Namespace(namespace).
		Create(
			c.ctx,
			uObj,
			metav1.CreateOptions{},
		)
	if err != nil {
		return nil, err
	}
	return resp.UnstructuredContent(), nil
}

// Get returns an object given its kind, name and namespace
func (c *Client) Get(kind string, name string, namespace string) (map[string]interface{}, error) {
	resource, err := knownKinds(kind)
	if err != nil {
		return nil, err
	}

	resp, err := c.dynamic.
		Resource(resource).
		Namespace(namespace).
		Get(
			c.ctx,
			name,
			metav1.GetOptions{},
		)
	if err != nil {
		return nil, err
	}
	return resp.UnstructuredContent(), nil
}

// List returns a list of objects given its kind and namespace
func (c *Client) List(kind string, namespace string) ([]map[string]interface{}, error) {
	resource, err := knownKinds(kind)
	if err != nil {
		return nil, err
	}
	resp, err := c.dynamic.
		Resource(resource).
		Namespace(namespace).
		List(c.ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	list := []map[string]interface{}{}
	for _, uObj := range resp.Items {
		list = append(list, uObj.UnstructuredContent())
	}
	return list, nil
}

// Delete deletes an object given its kind, name and namespace
func (c *Client) Delete(kind string, name string, namespace string) error {
	resource, err := knownKinds(kind)
	if err != nil {
		return err
	}

	err = c.dynamic.
		Resource(resource).
		Namespace(namespace).
		Delete(c.ctx, name, metav1.DeleteOptions{})

	return err
}
