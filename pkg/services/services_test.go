package services

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
	testName      = "svc-test"
	testNamespace = "ns-test"
)

func TestServices_Apply(t *testing.T) {
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
			expectedError: "YAML was not a Service",
		},
		{
			testID: "create new service from yaml",
			yaml: `
apiVersion: v1
kind: Service
metadata:
  name: ` + testName + `
  annotations:
  labels:
    app: xk6-kubernetes/unit-test
spec:
  selector:
    app: MyApp
  ports:
    - protocol: TCP
      port: 80
      targetPort: 9376
`,
			namespace:     testNamespace,
			expectedError: "",
		},
		{
			testID: "update existing service from yaml",
			yaml: `
apiVersion: v1
kind: Service
metadata:
  name: existing
  annotations:
  labels:
    app: xk6-kubernetes/unit-test
spec:
  selector:
    app: MyApp
  ports:
    - protocol: TCP
      port: 80
      targetPort: 9376
`,
			namespace:     testNamespace,
			expectedError: "services \"existing\" already exists",
		},
	}

	for _, tc := range testCases {
		tc := tc // pin the testcase
		t.Run(tc.testID, func(t *testing.T) {
			t.Parallel()
			fixture := New(context.Background(), fake.NewSimpleClientset(
				testutils.NewService("existing", testNamespace),
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

func TestServices_Create(t *testing.T) {
	t.Parallel()
	fixture := New(context.Background(), fake.NewSimpleClientset(
		testutils.NewService("existing", testNamespace),
	), metav1.ListOptions{})

	newOne := *testutils.NewService(testName, testNamespace)
	result, err := fixture.Create(newOne, testNamespace, metav1.CreateOptions{})
	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	if result.Name != newOne.Name || result.Namespace != newOne.Namespace {
		t.Errorf("incorrect instance was returned")
		return
	}
	svcs, _ := fixture.List(testNamespace)
	if len(svcs) != 2 {
		t.Errorf("expecting 2 services in namespace, listing returned %v", len(svcs))
		return
	}
}

func TestServices_List(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		testID        string
		namespace     string
		expectedCount int
	}{
		{
			testID:        "test namespace returns 3 services",
			namespace:     testNamespace,
			expectedCount: 3,
		},
		{
			testID:        "empty namespace returns 0 services",
			namespace:     "ns-empty",
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		tc := tc // pin the testcase
		t.Run(tc.testID, func(t *testing.T) {
			t.Parallel()
			fixture := New(context.Background(), fake.NewSimpleClientset(
				testutils.NewService("svc-1", testNamespace),
				testutils.NewService("svc-2", testNamespace),
				testutils.NewService("svc-3", testNamespace),
			), metav1.ListOptions{})

			result, err := fixture.List(tc.namespace)
			if err != nil {
				t.Errorf("encountered an error: %v", err)
				return
			}
			if len(result) != tc.expectedCount {
				t.Errorf("received %v service(s), expected %v", len(result), tc.expectedCount)
				return
			}
			for _, svc := range result {
				if tc.namespace != svc.Namespace {
					t.Errorf("received services from %v namespace, only expected %v", svc.Namespace, tc.namespace)
					return
				}
			}
		})
	}
}

func TestServices_Delete(t *testing.T) {
	t.Parallel()
	fixture := New(context.Background(), fake.NewSimpleClientset(
		testutils.NewService(testName, testNamespace),
	), metav1.ListOptions{})

	err := fixture.Delete(testName, testNamespace, metav1.DeleteOptions{})
	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	svcs, _ := fixture.List(testNamespace)
	if len(svcs) != 0 {
		t.Errorf("expecting 0 services in namespace, listing returned %v", len(svcs))
		return
	}
}

func TestServices_Get(t *testing.T) {
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
			name:         "svc-other",
			namespace:    testNamespace,
			expectToFind: false,
		},
		{
			testID:       "fetching unknown name from any namespace returns nothing",
			name:         "svc-unknown",
			namespace:    "any-namespace",
			expectToFind: false,
		},
	}

	for _, tc := range testCases {
		tc := tc // pin the testcase
		t.Run(tc.testID, func(t *testing.T) {
			t.Parallel()
			fixture := New(context.Background(), fake.NewSimpleClientset(
				testutils.NewService(testName, testNamespace),
				testutils.NewService("svc-other", "ns-2"),
			), metav1.ListOptions{})

			result, err := fixture.Get(tc.name, tc.namespace, metav1.GetOptions{})

			if err != nil && !errors.IsNotFound(err) {
				t.Errorf("encountered an error: %v", err)
				return
			}
			if tc.expectToFind && err != nil {
				t.Errorf("expected to find service %v in %v namespace, but received error: %v", tc.name, tc.namespace, err)
				return
			}
			if !tc.expectToFind && err == nil {
				t.Errorf("expected an error when trying to find unavailable service %v in %v namespace", tc.name, tc.namespace)
				return
			}
			if tc.expectToFind && result.Name != tc.name {
				t.Errorf("received service with name %v, expected %v", result.Name, tc.name)
				return
			}
			if tc.expectToFind && result.Namespace != tc.namespace {
				t.Errorf("received service with namespace %v, expected %v", result.Name, tc.name)
				return
			}
		})
	}
}
