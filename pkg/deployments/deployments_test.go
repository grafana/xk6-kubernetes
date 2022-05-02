package deployments

import (
	"github.com/grafana/xk6-kubernetes/internal/testutils"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"strings"
	"testing"
)

var (
	testName      = "deployment-test"
	testNamespace = "ns-test"
)

func TestDeployments_Apply(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewDeployment("existing", testNamespace),
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
			expectedError: "Yaml was not a Deployment",
		},
		{
			testID: "create new deployment from yaml",
			yaml: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ` + testName + `
  annotations:
  labels:
    app: xk6-kubernetes/unit-test
spec:
  replicas: 3
  selector:
    matchLabels:
      app: xk6-kubernetes/unit-test
  template:
    metadata:
      labels:
        app: xk6-kubernetes/unit-test
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
`,
			namespace:     testNamespace,
			expectedError: "",
		},
		{
			testID: "update existing deployment from yaml",
			yaml: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ` + testName + `
  annotations:
  labels:
    app: xk6-kubernetes/unit-test
spec:
  replicas: 3
  selector:
    matchLabels:
      app: xk6-kubernetes/unit-test
  template:
    metadata:
      labels:
        app: xk6-kubernetes/unit-test
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
`,
			namespace:     testNamespace,
			expectedError: "deployments.apps \"" + testName + "\" already exists",
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

func TestDeployments_Create(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewDeployment("existing", testNamespace),
	), metav1.ListOptions{}, nil)

	newOne := *testutils.NewDeployment(testName, testNamespace)
	result, err := fixture.Create(newOne, testNamespace, metav1.CreateOptions{})

	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	if result.Name != newOne.Name || result.Namespace != newOne.Namespace {
		t.Errorf("incorrect instance was returned")
		return
	}
	deployments, _ := fixture.List(testNamespace)
	if len(deployments) != 2 {
		t.Errorf("expecting 2 deployments in namespace, listing returned %v", len(deployments))
		return
	}
}

func TestDeployments_List(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewDeployment("deployment-1", testNamespace),
		testutils.NewDeployment("deployment-2", testNamespace),
		testutils.NewDeployment("deployment-3", testNamespace),
	), metav1.ListOptions{}, nil)

	testCases := []struct {
		testID        string
		namespace     string
		expectedCount int
	}{
		{
			testID:        "test namespace returns 3 deployments",
			namespace:     testNamespace,
			expectedCount: 3,
		},
		{
			testID:        "empty namespace returns 0 deployments",
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
				t.Errorf("received %v deployment(s), expected %v", len(result), tc.expectedCount)
				return
			}
			for _, deployment := range result {
				if tc.namespace != deployment.Namespace {
					t.Errorf("received deployment from %v namespace, only expected %v", deployment.Namespace, tc.namespace)
					return
				}
			}
		})
	}
}

func TestDeployments_Delete(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewDeployment(testName, testNamespace),
	), metav1.ListOptions{}, nil)

	err := fixture.Delete(testName, testNamespace, metav1.DeleteOptions{})

	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	deployments, _ := fixture.List(testNamespace)
	if len(deployments) != 0 {
		t.Errorf("expecting 0 deployments in namespace, listing returned %v", len(deployments))
		return
	}
}

func TestDeployments_Get(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewDeployment(testName, testNamespace),
		testutils.NewDeployment("deployment-other", "ns-2"),
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
			name:         "deployment-other",
			namespace:    testNamespace,
			expectToFind: false,
		},
		{
			testID:       "fetching unknown name from any namespace returns nothing",
			name:         "deployment-unknown",
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
				t.Errorf("expected to find deployment %v in %v namespace, but received error: %v", tc.name, tc.namespace, err)
				return
			}
			if !tc.expectToFind && err == nil {
				t.Errorf("expected an error when trying to find unavailable deployment %v in %v namespace", tc.name, tc.namespace)
				return
			}
			if tc.expectToFind && result.Name != tc.name {
				t.Errorf("received deployment with name %v, expected %v", result.Name, tc.name)
				return
			}
			if tc.expectToFind && result.Namespace != tc.namespace {
				t.Errorf("received deployment with namespace %v, expected %v", result.Name, tc.name)
				return
			}
		})
	}
}
