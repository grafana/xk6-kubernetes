package namespaces

import (
	"github.com/grafana/xk6-kubernetes/internal/testutils"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"strings"
	"testing"
)

var (
	testName = "ns-test"
)

func TestNamespaces_Apply(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(testutils.NewNamespace("existing")), metav1.ListOptions{}, nil)

	testCases := []struct {
		testID        string
		yaml          string
		expectedError string
	}{
		{
			testID: "apply with invalid yaml",
			yaml: `
this is not yaml
`,
			expectedError: "json parse error",
		},
		{
			testID: "apply with incorrect resource yaml",
			yaml: `
apiVersion: v1
kind: Secret
metadata:
  name: ` + testName + `
data:
  secret-key: MWYyZDFlMmU2N2Rm
`,
			expectedError: "Yaml was not a Namespace",
		},
		{
			testID: "create new namespace from yaml",
			yaml: `
apiVersion: v1
kind: Namespace
metadata:
  name: ` + testName + `
  annotations:
  labels:
    app: xk6-kubernetes/unit-test
`,
			expectedError: "",
		},
		{
			testID: "update existing namespace from yaml",
			yaml: `
apiVersion: v1
kind: Namespace
metadata:
  name: ` + testName + `
  annotations:
  labels:
    app: xk6-kubernetes/unit-test
`,
			expectedError: "namespaces \"" + testName + "\" already exists",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testID, func(t *testing.T) {
			result, err := fixture.Apply(tc.yaml)

			if err != nil && (tc.expectedError == "" || !strings.Contains(err.Error(), tc.expectedError)) {
				t.Errorf("encountered error: %v, expected: %v", err, tc.expectedError)
				return
			}
			if err == nil && tc.expectedError != "" {
				t.Errorf("expected error \"%v\", but got none", tc.expectedError)
				return
			}
			if err == nil && result.Name != testName {
				t.Errorf("incorrect instance was returned")
				return
			}
		})
	}
}

func TestNamespaces_Create(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(testutils.NewNamespace("existing")), metav1.ListOptions{}, nil)

	newOne := *testutils.NewNamespace(testName)
	result, err := fixture.Create(newOne, metav1.CreateOptions{})

	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	if result.Name != newOne.Name {
		t.Errorf("incorrect instance was returned")
		return
	}
	namespaces, _ := fixture.List()
	if len(namespaces) != 2 {
		t.Errorf("expecting 2 namespaces, listing returned %v", len(namespaces))
		return
	}
}

func TestNamespaces_List(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewNamespace("ns-1"),
		testutils.NewNamespace("ns-2"),
		testutils.NewNamespace("ns-3"),
	), metav1.ListOptions{}, nil)

	result, err := fixture.List()
	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	if len(result) != 3 {
		t.Errorf("received %v namespace(s), expected %v", len(result), 3)
		return
	}
}

func TestNamespaces_Delete(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewNamespace(testName),
	), metav1.ListOptions{}, nil)

	err := fixture.Delete(testName, metav1.DeleteOptions{})

	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	namespaces, _ := fixture.List()
	if len(namespaces) != 0 {
		t.Errorf("expecting 0 namespaces, listing returned %v", len(namespaces))
		return
	}
}

func TestNamespaces_Get(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(testutils.NewNamespace(testName)), metav1.ListOptions{}, nil)

	testCases := []struct {
		testID       string
		name         string
		expectToFind bool
	}{
		{
			testID:       "fetching valid name returns correctly",
			name:         testName,
			expectToFind: true,
		},
		{
			testID:       "fetching unknown name from any namespace returns nothing",
			name:         "ns-unknown",
			expectToFind: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testID, func(t *testing.T) {
			result, err := fixture.Get(tc.name, metav1.GetOptions{})
			if err != nil && !errors.IsNotFound(err) {
				t.Errorf("encountered an error: %v", err)
				return
			}
			if tc.expectToFind && err != nil {
				t.Errorf("expected to find namespace %v, but received error: %v", tc.name, err)
				return
			}
			if !tc.expectToFind && err == nil {
				t.Errorf("expected an error when trying to find unavailable namespace %v", tc.name)
				return
			}
			if tc.expectToFind && result.Name != tc.name {
				t.Errorf("received namespace with name %v, expected %v", result.Name, tc.name)
				return
			}
		})
	}
}
