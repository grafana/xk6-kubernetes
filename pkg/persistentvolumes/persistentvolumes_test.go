package persistentvolumes

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
	testName = "pv-test"
)

func TestPersistentVolumes_Apply(t *testing.T) {
	t.Parallel()
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
kind: Namespace
metadata:
  name: ` + testName + `
`,
			expectedError: "YAML was not a PersistentVolume",
		},
		{
			testID: "create new persistent volume from yaml",
			yaml: `
apiVersion: v1
kind: PersistentVolume
metadata:
  name: ` + testName + `
  annotations:
  labels:
    app: xk6-kubernetes/unit-test
spec:
  storageClassName: local-storage
  capacity:
    storage: 10G
  accessModes:
  - ReadWriteMany
  hostPath:
    path: "/tmp/xk6-test"
`,
			expectedError: "",
		},
		{
			testID: "update existing persistent volume from yaml",
			yaml: `
apiVersion: v1
kind: PersistentVolume
metadata:
  name: existing
  annotations:
  labels:
    app: xk6-kubernetes/unit-test
spec:
  storageClassName: local-storage
  capacity:
    storage: 10G
  accessModes:
  - ReadWriteMany
  hostPath:
    path: "/tmp/xk6-test"
`,
			expectedError: "persistentvolumes \"existing\" already exists",
		},
	}

	for _, tc := range testCases {
		tc := tc // pin the testcase
		t.Run(tc.testID, func(t *testing.T) {
			t.Parallel()
			fixture := New(context.Background(), fake.NewSimpleClientset(
				testutils.NewPersistentVolume("existing"),
			), metav1.ListOptions{})

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

func TestPersistentVolumes_Create(t *testing.T) {
	t.Parallel()
	fixture := New(context.Background(), fake.NewSimpleClientset(
		testutils.NewPersistentVolume("existing"),
	), metav1.ListOptions{})

	newOne := *testutils.NewPersistentVolume(testName)
	result, err := fixture.Create(newOne, metav1.CreateOptions{})
	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	if result.Name != newOne.Name {
		t.Errorf("incorrect instance was returned")
		return
	}
	pvs, _ := fixture.List()
	if len(pvs) != 2 {
		t.Errorf("expecting 2 persistent volumes in namespace, listing returned %v", len(pvs))
		return
	}
}

func TestPersistentVolumes_List(t *testing.T) {
	t.Parallel()
	fixture := New(context.Background(), fake.NewSimpleClientset(
		testutils.NewPersistentVolume("pv-1"),
		testutils.NewPersistentVolume("pv-2"),
		testutils.NewPersistentVolume("pv-3"),
	), metav1.ListOptions{})

	result, err := fixture.List()
	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	if len(result) != 3 {
		t.Errorf("received %v persistent volume(s), expected %v", len(result), 3)
		return
	}
}

func TestPersistentVolumes_Delete(t *testing.T) {
	t.Parallel()
	fixture := New(context.Background(), fake.NewSimpleClientset(
		testutils.NewPersistentVolume(testName),
	), metav1.ListOptions{})

	err := fixture.Delete(testName, metav1.DeleteOptions{})
	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	pvs, _ := fixture.List()
	if len(pvs) != 0 {
		t.Errorf("expecting 0 persistent volumes, listing returned %v", len(pvs))
		return
	}
}

func TestPersistentVolumes_Get(t *testing.T) {
	t.Parallel()
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
			name:         "pv-unknown",
			expectToFind: false,
		},
	}

	for _, tc := range testCases {
		tc := tc // pin the testcase
		t.Run(tc.testID, func(t *testing.T) {
			t.Parallel()
			fixture := New(context.Background(), fake.NewSimpleClientset(
				testutils.NewPersistentVolume(testName),
			), metav1.ListOptions{})

			result, err := fixture.Get(tc.name, metav1.GetOptions{})

			if err != nil && !errors.IsNotFound(err) {
				t.Errorf("encountered an error: %v", err)
				return
			}
			if tc.expectToFind && err != nil {
				t.Errorf("expected to find persistent volume %v, but received error: %v", tc.name, err)
				return
			}
			if !tc.expectToFind && err == nil {
				t.Errorf("expected an error when trying to find unavailable persistent volume %v", tc.name)
				return
			}
			if tc.expectToFind && result.Name != tc.name {
				t.Errorf("received persistent volume with name %v, expected %v", result.Name, tc.name)
				return
			}
		})
	}
}
