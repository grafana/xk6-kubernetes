package testutils

import (
	"k8s.io/apimachinery/pkg/runtime"
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
