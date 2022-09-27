package testutils

import (
	"k8s.io/apimachinery/pkg/runtime"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/fake"
)

// NewFakeDynamic creates a new instance of a fake dynamic client with a default scheme
func NewFakeDynamic() *dynamicfake.FakeDynamicClient {
	scheme := runtime.NewScheme()
	fake.AddToScheme(scheme)
	return dynamicfake.NewSimpleDynamicClient(scheme)
}
