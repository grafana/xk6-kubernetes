package testutils

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/fake"
)

// NewFakeDynamic creates a new instance of a fake dynamic client with a default scheme
func NewFakeDynamic(objs ...runtime.Object) (*dynamicfake.FakeDynamicClient, error) {
	scheme := runtime.NewScheme()
	err := fake.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}

	return dynamicfake.NewSimpleDynamicClient(scheme, objs...), nil
}

type FakeRESTMapper struct {
	meta.RESTMapper
}

func (f *FakeRESTMapper) RESTMapping(gk schema.GroupKind, versions ...string) (*meta.RESTMapping, error) {
	kindMapping := map[string]schema.GroupVersionResource{
		"ConfigMap":             {Group: "", Version: "v1", Resource: "configmaps"},
		"Deployment":            {Group: "apps", Version: "v1", Resource: "deployments"},
		"Endpoint":              {Group: "", Version: "v1", Resource: "endpoints"},
		"Ingress":               {Group: "networking.k8s.io", Version: "v1", Resource: "ingresses"},
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

	gvr, found := kindMapping[gk.Kind]
	if !found {
		return nil, fmt.Errorf("unknown kind: '%s'", gk.Kind)
	}
	scope := meta.RESTScopeNamespace
	if gk.Kind == "Namespace" || gk.Kind == "Node" {
		scope = meta.RESTScopeRoot
	}

	return &meta.RESTMapping{
		Resource:         gvr,
		GroupVersionKind: gvr.GroupVersion().WithKind(gk.Kind),
		Scope:            scope,
	}, nil
}
