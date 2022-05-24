package configmaps

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
	testName      = "cm-test"
	testNamespace = "ns-test"
)

func TestConfigMaps_Apply(t *testing.T) {
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
			expectedError: "YAML was not a ConfigMap",
		},
		{
			testID: "create new configmap from yaml",
			yaml: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: ` + testName + `
  annotations:
  labels:
    app: xk6-kubernetes/unit-test
data:
  key-1: value-1
  key-2: value-2
`,
			namespace:     testNamespace,
			expectedError: "",
		},
		{
			testID: "update existing configmap from yaml",
			yaml: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: existing
  annotations:
  labels:
    app: xk6-kubernetes/unit-test
data:
  key-1: updated-value-1
  key-2: updated-value-2
`,
			namespace:     testNamespace,
			expectedError: "configmaps \"existing\" already exists",
		},
	}

	for _, tc := range testCases {
		tc := tc // pin the testcase
		t.Run(tc.testID, func(t *testing.T) {
			t.Parallel()
			fixture := New(context.Background(), fake.NewSimpleClientset(
				testutils.NewConfigMap("existing", testNamespace),
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

func TestConfigMaps_Create(t *testing.T) {
	t.Parallel()
	fixture := New(context.Background(), fake.NewSimpleClientset(
		testutils.NewConfigMap("existing", testNamespace),
	), metav1.ListOptions{})

	newOne := *testutils.NewConfigMap(testName, testNamespace)
	result, err := fixture.Create(newOne, testNamespace, metav1.CreateOptions{})
	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	if result.Name != newOne.Name || result.Namespace != newOne.Namespace {
		t.Errorf("incorrect instance was returned")
		return
	}
	cms, _ := fixture.List(testNamespace)
	if len(cms) != 2 {
		t.Errorf("expecting 2 configmaps in namespace, listing returned %v", len(cms))
		return
	}
}

func TestConfigMaps_List(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		testID        string
		namespace     string
		expectedCount int
	}{
		{
			testID:        "test namespace returns 3 configmaps",
			namespace:     testNamespace,
			expectedCount: 3,
		},
		{
			testID:        "empty namespace returns 0 configmaps",
			namespace:     "ns-empty",
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		tc := tc // pin the testcase
		t.Run(tc.testID, func(t *testing.T) {
			t.Parallel()
			fixture := New(context.Background(), fake.NewSimpleClientset(
				testutils.NewConfigMap("cm-1", testNamespace),
				testutils.NewConfigMap("cm-2", testNamespace),
				testutils.NewConfigMap("cm-3", testNamespace),
			), metav1.ListOptions{})

			result, err := fixture.List(tc.namespace)
			if err != nil {
				t.Errorf("encountered an error: %v", err)
				return
			}
			if len(result) != tc.expectedCount {
				t.Errorf("received %v configmap(s), expected %v", len(result), tc.expectedCount)
				return
			}
			for _, cm := range result {
				if tc.namespace != cm.Namespace {
					t.Errorf("received configmap from %v namespace, only expected %v", cm.Namespace, tc.namespace)
					return
				}
			}
		})
	}
}

func TestConfigMaps_Delete(t *testing.T) {
	t.Parallel()
	fixture := New(context.Background(), fake.NewSimpleClientset(
		testutils.NewConfigMap(testName, testNamespace),
	), metav1.ListOptions{})

	err := fixture.Delete(testName, testNamespace, metav1.DeleteOptions{})
	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	cms, _ := fixture.List(testNamespace)
	if len(cms) != 0 {
		t.Errorf("expecting 0 configmaps in namespace, listing returned %v", len(cms))
		return
	}
}

func TestConfigMaps_Get(t *testing.T) {
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
			name:         "cm-other",
			namespace:    testNamespace,
			expectToFind: false,
		},
		{
			testID:       "fetching unknown name from any namespace returns nothing",
			name:         "cm-unknown",
			namespace:    "any-namespace",
			expectToFind: false,
		},
	}

	for _, tc := range testCases {
		tc := tc // pin the testcase
		t.Run(tc.testID, func(t *testing.T) {
			t.Parallel()
			fixture := New(context.Background(), fake.NewSimpleClientset(
				testutils.NewConfigMap(testName, testNamespace),
				testutils.NewConfigMap("cm-other", "ns-2"),
			), metav1.ListOptions{})

			result, err := fixture.Get(tc.name, tc.namespace, metav1.GetOptions{})

			if err != nil && !errors.IsNotFound(err) {
				t.Errorf("encountered an error: %v", err)
				return
			}
			if tc.expectToFind && err != nil {
				t.Errorf("expected to find configmap %v in %v namespace, but received error: %v", tc.name, tc.namespace, err)
				return
			}
			if !tc.expectToFind && err == nil {
				t.Errorf("expected an error when trying to find unavailable configmap %v in %v namespace", tc.name, tc.namespace)
				return
			}
			if tc.expectToFind && result.Name != tc.name {
				t.Errorf("received configmap with name %v, expected %v", result.Name, tc.name)
				return
			}
			if tc.expectToFind && result.Namespace != tc.namespace {
				t.Errorf("received configmap with namespace %v, expected %v", result.Name, tc.name)
				return
			}
		})
	}
}
