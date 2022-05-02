package persistentvolumeclaims

import (
	"github.com/grafana/xk6-kubernetes/internal/testutils"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"strings"
	"testing"
)

var (
	testName      = "pvc-test"
	testNamespace = "ns-test"
)

func TestPersistentVolumeClaims_Apply(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewPersistentVolumeClaim("existing", testNamespace),
	), metav1.ListOptions{}, nil)

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
			expectedError: "Yaml was not a PersistentVolumeClaim",
		},
		{
			testID: "create new persistent volume claim from yaml",
			yaml: `
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: ` + testName + `
  annotations:
  labels:
    app: xk6-kubernetes/unit-test
spec:
  storageClassName: local-storage
  accessModes:
  - ReadWriteMany
  resources:
    requests:
      storage: 10G
`,
			namespace:     testNamespace,
			expectedError: "",
		},
		{
			testID: "update existing persistent volume claim from yaml",
			yaml: `
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: ` + testName + `
  annotations:
  labels:
    app: xk6-kubernetes/unit-test
spec:
  storageClassName: local-storage
  accessModes:
  - ReadWriteMany
  resources:
    requests:
      storage: 10G
`,
			namespace:     testNamespace,
			expectedError: "persistentvolumeclaims \"" + testName + "\" already exists",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testID, func(t *testing.T) {
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

func TestPersistentVolumeClaims_Create(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewPersistentVolumeClaim("existing", testNamespace),
	), metav1.ListOptions{}, nil)

	newOne := *testutils.NewPersistentVolumeClaim(testName, testNamespace)
	result, err := fixture.Create(newOne, testNamespace, metav1.CreateOptions{})

	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	if result.Name != newOne.Name || result.Namespace != newOne.Namespace {
		t.Errorf("incorrect instance was returned")
		return
	}
	pvcs, _ := fixture.List(testNamespace)
	if len(pvcs) != 2 {
		t.Errorf("expecting 2 persistent volume claims in namespace, listing returned %v", len(pvcs))
		return
	}
}

func TestPersistentVolumeClaims_List(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewPersistentVolumeClaim("pvc-1", testNamespace),
		testutils.NewPersistentVolumeClaim("pvc-2", testNamespace),
		testutils.NewPersistentVolumeClaim("pvc-3", testNamespace),
	), metav1.ListOptions{}, nil)

	testCases := []struct {
		testID        string
		namespace     string
		expectedCount int
	}{
		{
			testID:        "test namespace returns 3 persistent volume claims",
			namespace:     testNamespace,
			expectedCount: 3,
		},
		{
			testID:        "empty namespace returns 0 persistent volume claims",
			namespace:     "ns-empty",
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testID, func(t *testing.T) {
			result, err := fixture.List(tc.namespace)
			if err != nil {
				t.Errorf("encountered an error: %v", err)
				return
			}
			if len(result) != tc.expectedCount {
				t.Errorf("received %v persistent volume claim(s), expected %v", len(result), tc.expectedCount)
				return
			}
			for _, pvc := range result {
				if tc.namespace != pvc.Namespace {
					t.Errorf("received persistent volume claim from %v namespace, only expected %v", pvc.Namespace, tc.namespace)
					return
				}
			}
		})
	}
}

func TestPersistentVolumeClaims_Delete(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewPersistentVolumeClaim(testName, testNamespace),
	), metav1.ListOptions{}, nil)

	err := fixture.Delete(testName, testNamespace, metav1.DeleteOptions{})

	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	pvcs, _ := fixture.List(testNamespace)
	if len(pvcs) != 0 {
		t.Errorf("expecting 0 persistent volume claims in namespace, listing returned %v", len(pvcs))
		return
	}
}

func TestPersistentVolumeClaims_Get(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewPersistentVolumeClaim(testName, testNamespace),
		testutils.NewPersistentVolumeClaim("pvc-other", "ns-2"),
	), metav1.ListOptions{}, nil)

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
			name:         "pvc-other",
			namespace:    testNamespace,
			expectToFind: false,
		},
		{
			testID:       "fetching unknown name from any namespace returns nothing",
			name:         "pvc-unknown",
			namespace:    "any-namespace",
			expectToFind: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testID, func(t *testing.T) {
			result, err := fixture.Get(tc.name, tc.namespace, metav1.GetOptions{})
			if err != nil && !errors.IsNotFound(err) {
				t.Errorf("encountered an error: %v", err)
				return
			}
			if tc.expectToFind && err != nil {
				t.Errorf("expected to find persistent volume claim %v in %v namespace, but received error: %v", tc.name, tc.namespace, err)
				return
			}
			if !tc.expectToFind && err == nil {
				t.Errorf("expected an error when trying to find unavailable persistent volume claim %v in %v namespace", tc.name, tc.namespace)
				return
			}
			if tc.expectToFind && result.Name != tc.name {
				t.Errorf("received persistent volume claim with name %v, expected %v", result.Name, tc.name)
				return
			}
			if tc.expectToFind && result.Namespace != tc.namespace {
				t.Errorf("received persistent volume claim with namespace %v, expected %v", result.Name, tc.name)
				return
			}
		})
	}
}
