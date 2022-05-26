package ingresses

import (
	"context"
	"strings"
	"testing"

	"github.com/grafana/xk6-kubernetes/internal/testutils"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	testName      = "ingress-test"
	testNamespace = "ns-test"
)

func TestIngresses_Apply(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		testID        string
		yaml          string
		namespace     string
		expectedError string
	}{
		{
			testID: "apply with invalid yaml",
			yaml: `
this is not yaml
`,
			namespace:     testNamespace,
			expectedError: "json parse error",
		},
		{
			testID: "apply with incorrect resource yaml",
			yaml: `
apiVersion: v1
kind: Namespace
metadata:
  name: ` + testName + `
`,
			namespace:     testNamespace,
			expectedError: "YAML was not an Ingress",
		},
		{
			testID: "create new ingress from yaml",
			yaml: `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ` + testName + `
  annotations:
  labels:
    app: xk6-kubernetes/unit-test
spec:
  ingressClassName: nginx-example
  rules:
  - http:
      paths:
      - path: /testpath
        pathType: Prefix
        backend:
          service:
            name: ` + testName + `
            port:
              number: 80
`,
			namespace:     testNamespace,
			expectedError: "",
		},
		{
			testID: "update existing ingress from yaml",
			yaml: `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: existing
  annotations:
  labels:
    app: xk6-kubernetes/unit-test
spec:
  ingressClassName: nginx-example
  rules:
  - http:
      paths:
      - path: /testpath
        pathType: Prefix
        backend:
          service:
            name: existing
            port:
              number: 80
`,
			namespace:     testNamespace,
			expectedError: "ingresses.networking.k8s.io \"existing\" already exists",
		},
	}

	for _, tc := range testCases {
		tc := tc // pin the testcase
		t.Run(tc.testID, func(t *testing.T) {
			t.Parallel()
			fixture := New(context.Background(), fake.NewSimpleClientset(
				testutils.NewIngress("existing", testNamespace),
			), metav1.ListOptions{})

			result, err := fixture.Apply(tc.yaml, tc.namespace)

			if err != nil && (tc.expectedError == "" || !strings.Contains(err.Error(), tc.expectedError)) {
				t.Errorf("encountered error: %v, expected: %v", err, tc.expectedError)
				return
			}
			if err == nil && tc.expectedError != "" {
				t.Errorf("expected error \"%v\", but got none", tc.expectedError)
				return
			}
			if err == nil && (result.Name != testName || result.Namespace != testNamespace) {
				t.Errorf("incorrect instance was returned")
				return
			}
		})
	}
}

func TestIngresses_Create(t *testing.T) {
	t.Parallel()
	fixture := New(context.Background(), fake.NewSimpleClientset(
		testutils.NewIngress("existing", testNamespace),
	), metav1.ListOptions{})

	newOne := *testutils.NewIngress(testName, testNamespace)
	result, err := fixture.Create(newOne, testNamespace, metav1.CreateOptions{})
	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	if result.Name != newOne.Name || result.Namespace != newOne.Namespace {
		t.Errorf("incorrect instance was returned")
		return
	}
	ingresses, _ := fixture.List(testNamespace)
	if len(ingresses) != 2 {
		t.Errorf("expecting 2 ingresses in namespace, listing returned %v", len(ingresses))
		return
	}
}

func TestIngresses_List(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		testID        string
		namespace     string
		expectedCount int
	}{
		{
			testID:        "test namespace returns 3 ingresses",
			namespace:     testNamespace,
			expectedCount: 3,
		},
		{
			testID:        "empty namespace returns 0 ingresses",
			namespace:     "ns-empty",
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		tc := tc // pin the testcase
		t.Run(tc.testID, func(t *testing.T) {
			t.Parallel()
			fixture := New(context.Background(), fake.NewSimpleClientset(
				testutils.NewIngress("ingress-1", testNamespace),
				testutils.NewIngress("ingress-2", testNamespace),
				testutils.NewIngress("ingress-3", testNamespace),
			), metav1.ListOptions{})

			result, err := fixture.List(tc.namespace)
			if err != nil {
				t.Errorf("encountered an error: %v", err)
				return
			}
			if len(result) != tc.expectedCount {
				t.Errorf("received %v ingress(es), expected %v", len(result), tc.expectedCount)
				return
			}
			for _, ingress := range result {
				if tc.namespace != ingress.Namespace {
					t.Errorf("received ingress from %v namespace, only expected %v", ingress.Namespace, tc.namespace)
					return
				}
			}
		})
	}
}

func TestIngresses_Delete(t *testing.T) {
	t.Parallel()
	fixture := New(context.Background(), fake.NewSimpleClientset(
		testutils.NewIngress(testName, testNamespace),
	), metav1.ListOptions{})

	err := fixture.Delete(testName, testNamespace, metav1.DeleteOptions{})
	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	ingresses, _ := fixture.List(testNamespace)
	if len(ingresses) != 0 {
		t.Errorf("expecting 0 ingresses in namespace, listing returned %v", len(ingresses))
		return
	}
}

func TestIngresses_Get(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		testID       string
		name         string
		namespace    string
		expectToFind bool
	}{
		{
			testID:       "fetching valid name within namespace returns correctly",
			name:         testName,
			namespace:    testNamespace,
			expectToFind: true,
		},
		{
			testID:       "fetching valid name from incorrect namespace returns nothing",
			name:         "ingress-other",
			namespace:    testNamespace,
			expectToFind: false,
		},
		{
			testID:       "fetching unknown name from any namespace returns nothing",
			name:         "ingress-unknown",
			namespace:    "any-namespace",
			expectToFind: false,
		},
	}

	for _, tc := range testCases {
		tc := tc // pin the testcase
		t.Run(tc.testID, func(t *testing.T) {
			t.Parallel()
			fixture := New(context.Background(), fake.NewSimpleClientset(
				testutils.NewIngress(testName, testNamespace),
				testutils.NewIngress("ingress-other", "ns-2"),
			), metav1.ListOptions{})

			result, err := fixture.Get(tc.name, tc.namespace, metav1.GetOptions{})

			if err != nil && !errors.IsNotFound(err) {
				t.Errorf("encountered an error: %v", err)
				return
			}
			if tc.expectToFind && err != nil {
				t.Errorf("expected to find ingress %v in %v namespace, but received error: %v", tc.name, tc.namespace, err)
				return
			}
			if !tc.expectToFind && err == nil {
				t.Errorf("expected an error when trying to find unavailable ingress %v in %v namespace", tc.name, tc.namespace)
				return
			}
			if tc.expectToFind && result.Name != tc.name {
				t.Errorf("received ingress with name %v, expected %v", result.Name, tc.name)
				return
			}
			if tc.expectToFind && result.Namespace != tc.namespace {
				t.Errorf("received ingress with namespace %v, expected %v", result.Name, tc.name)
				return
			}
		})
	}
}
