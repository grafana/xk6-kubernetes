package helpers

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceHelper offers utility methods for manipulating namespaces
type NamespaceHelper interface {
	// Creates a new namespace with a random name that starts with the given prefix
	// e.g given the prefix 'test-' in will create a namespace of the form 'test-af8hx5'
	// where 'af8hx' is a random sequence of characters generated for each namespace
	CreateRandomNamespace(prefix string) (string, error)
}

func (h *helpers) CreateRandomNamespace(prefix string) (string, error) {
	ns := corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{GenerateName: prefix},
	}

	created, err := h.client.Structured().Create(ns)
	if err != nil {
		return "", err
	}

	return created.(corev1.Namespace).Name, nil
}
